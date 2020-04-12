package onvif

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/clbanning/mxj"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var errWrongDiscoveryResponse = errors.New("Response is not related to discovery request")

// StartDiscovery send a WS-Discovery message and wait for all matching device to respond
func StartDiscovery(duration time.Duration) ([]*Device, error) {
	// Get list of interface address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []*Device{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	return StartDiscoveryWithContext(ctx, addrs, duration)
}

func StartDiscoveryWithContext(ctx context.Context, addrs []net.Addr, duration time.Duration) ([]*Device, error) {
	eg, ctx := errgroup.WithContext(ctx)

	// Create initial discovery results
	discoveryResults := []*Device{}

	// Fetch IPv4 address
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if ok && !ipAddr.IP.IsLoopback() && ipAddr.IP.To4() != nil {
			eg.Go(func() error {
				devices, err := discoverDevices(ipAddr, duration)
				if err != nil {
					return err
				}

				discoveryResults = append(discoveryResults, devices...)

				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.Wrap(err, "Error waiting for discovery to complete")
	}

	return discoveryResults, nil
}

func discoverDevices(ipAddr *net.IPNet, duration time.Duration) ([]*Device, error) {
	log.Debugf("discoverDevices. IP: %s. Duration: %s", ipAddr, duration)
	var now = time.Now()
	// Create WS-Discovery request
	messageID := "uuid:" + uuid.Must(uuid.NewV4()).String()

	// var request = `
	// <?xml version="1.0" encoding="utf-8"?>
	// <Envelope xmlns:dn="http://www.onvif.org/ver10/network/wsdl"
	// xmlns="http://www.w3.org/2003/05/soap-envelope">
	// <Header>
	// <wsa:MessageID xmlns:wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing">` + messageID + `</wsa:MessageID>
	// <wsa:To xmlns:wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing">urn:schemas-xmlsoap-org:ws:2005:04:discovery</wsa:To>
	// <wsa:Action xmlns:wsa="http://schemas.xmlsoap.org/ws/2004/08/addressing">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</wsa:Action>
	// </Header>
	// <Body>
	// <Probe xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	// xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	// xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery">
	// <Types>dn:NetworkVideoTransmitter</Types>
	// <Scopes />
	// </Probe>
	// </Body>
	// </Envelope>
	// `
	var request = `
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
  <s:Header>
    <a:Action s:mustUnderstand="1">http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
    <a:MessageID>` + messageID + `</a:MessageID>
    <a:ReplyTo>
      <a:Address>http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <a:To s:mustUnderstand="1">urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
  </s:Header>
  <s:Body>
    <Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery">
      <d:Types xmlns:d="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:dp0="http://www.onvif.org/ver10/network/wsdl">dp0:NetworkVideoTransmitter</d:Types>
    </Probe>
  </s:Body>
</s:Envelope>
`
	// Clean WS-Discovery message
	request = regexp.MustCompile(`\>\s+\<`).ReplaceAllString(request, "><")
	request = regexp.MustCompile(`\s+`).ReplaceAllString(request, " ")
	// Create UDP address for local and multicast address
	localAddress, err := net.ResolveUDPAddr("udp4", ipAddr.IP.String()+":0")
	if err != nil {
		return []*Device{}, err
	}

	// fmt.Println("ip", ipAddr, duration, localAddress)

	multicastAddress, err := net.ResolveUDPAddr("udp4", "239.255.255.250:3702")
	if err != nil {
		return []*Device{}, err
	}

	// Create UDP connection to listen for respond from matching device
	conn, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		return []*Device{}, err
	}
	defer conn.Close()

	// Set connection's timeout
	err = conn.SetDeadline(now.Add(duration))
	if err != nil {
		return []*Device{}, err
	}

	// Send WS-Discovery request to multicast address
	_, err = conn.WriteToUDP([]byte(request), multicastAddress)
	if err != nil {
		return []*Device{}, err
	}

	// Create initial discovery results
	discoveryResults := []*Device{}

	// Keep reading UDP message until timeout
	for {
		// Create buffer and receive UDP response
		buffer := make([]byte, 16*1024)
		_, udpAddr, err := conn.ReadFromUDP(buffer)

		// Check if connection timeout
		if err != nil {
			if udpErr, ok := err.(net.Error); ok && udpErr.Timeout() {
				break
			} else {
				return discoveryResults, err
			}
		}

		log.Debugf("Camera replied. Data: %s", string(buffer))
		// Read and parse WS-Discovery response
		devices, err := readDiscoveryResponse(messageID, buffer)
		if err != nil && err != errWrongDiscoveryResponse {
			return discoveryResults, err
		}

		// fmt.Println(now, device, err)

		// Push device to results
		for _, device := range devices {
			parsed, _ := url.Parse(device.XAddr)

			ipA, _, _ := net.ParseCIDR(fmt.Sprintf("%s/32", parsed.Hostname()))

			// log.Debug(ipAddr.IP.String(), ipAddr.Mask.String(), "contains", parsed.Hostname(), " = ", ipAddr.Contains(ipA))

			if ipAddr.Contains(ipA) {
				device.IPAddress = udpAddr.IP.String()
				discoveryResults = append(discoveryResults, device)
			}
		}
	}

	return discoveryResults, nil
}

// readDiscoveryResponse reads and parses WS-Discovery response
func readDiscoveryResponse(messageID string, buffer []byte) ([]*Device, error) {
	// Parse XML to map
	mapXML, err := mxj.NewMapXml(buffer)
	if err != nil {
		return nil, err
	}

	// Check if this response is for our request
	responseMessageID, _ := mapXML.ValueForPathString("Envelope.Header.RelatesTo")
	if responseMessageID != messageID {
		return nil, errWrongDiscoveryResponse
	}

	// Get device's ID and clean it
	deviceID, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.EndpointReference.Address")
	deviceID = strings.Replace(deviceID, "urn:uuid:", "", 1)

	// Get device's name
	deviceName := ""
	deviceMAC := ""
	scopes, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.Scopes")
	for _, scope := range strings.Split(scopes, " ") {
		if strings.HasPrefix(scope, "onvif://www.onvif.org/name/") {
			deviceName = strings.Replace(scope, "onvif://www.onvif.org/name/", "", 1)
			deviceName = strings.Replace(deviceName, "_", " ", -1)
		}
		if strings.HasPrefix(scope, "onvif://www.onvif.org/MAC/") {
			deviceMAC = strings.Replace(scope, "onvif://www.onvif.org/MAC/", "", 1)
		}
	}

	// Get device's xAddrs
	xAddrs, _ := mapXML.ValueForPathString("Envelope.Body.ProbeMatches.ProbeMatch.XAddrs")
	listXAddr := strings.Split(xAddrs, " ")
	if len(listXAddr) == 0 {
		return nil, errors.New("Device does not have any xAddr")
	}

	var devices = make([]*Device, len(listXAddr))

	for idx, xAddr := range listXAddr {
		devices[idx] = &Device{
			ID:      deviceID,
			Name:    deviceName,
			MACAddr: deviceMAC,
			XAddr:   xAddr,
		}
	}

	return devices, nil
}
