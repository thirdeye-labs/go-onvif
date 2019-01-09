package onvif

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type NetworkProtocolsResponse struct {
	NetworkProtocols []networkProtocol `json:"NetworkProtocols"`
}
type networkProtocol struct {
	Enabled string `json:"Enabled"`
	Name    string `json:"Name"`
	Port    string `json:"Port"`
}

type NetworkProtocol struct {
	Enabled bool
	Name    string
	Port    int64
}

// GetNetworkProtocols fetches the network protocols that you can access the
// device on.
func (device Device) GetNetworkProtocols() ([]NetworkProtocol, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: mediaXMLNs,
		Body: `<trt:GetNetworkProtocols>
		</trt:GetNetworkProtocols>`,
		User:     device.User,
		Password: device.Password,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return nil, errors.Wrap(err, "GetNetworkProtocols: Could not send SOAP request")
	}

	// Parse response to interface
	npResp, err := response.ValueForPath("Envelope.Body.GetNetworkProtocolsResponse")
	if err != nil {
		return nil, errors.Wrap(err, "GetNetworkProtocols: Parse network protocols response")
	}

	m, err := json.MarshalIndent(npResp, "  ", "  ")

	var npr NetworkProtocolsResponse
	if err := json.Unmarshal(m, &npr); err != nil {
		return nil, errors.Wrap(err, "GetNetworkProtocols: Could not unmarshal protocols response")
	}

	var nps = make([]NetworkProtocol, len(npr.NetworkProtocols))
	for idx, np := range npr.NetworkProtocols {
		port, err := strconv.ParseInt(np.Port, 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		nps[idx] = NetworkProtocol{
			Enabled: np.Enabled == "true",
			Name:    np.Name,
			Port:    port,
		}
	}

	return nps, nil
}
