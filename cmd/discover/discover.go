package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	. "github.com/thirdeye-labs/go-onvif"
)

func main() {

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var addrs []net.Addr

	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "en") || strings.Contains(iface.Name, "wlp") || strings.Contains(iface.Name, "wlan") || strings.Contains(iface.Name, "br0") || strings.Contains(iface.Name, "enx") {
			ifaceAddrs, err := iface.Addrs()

			if err != nil {
				panic(err)
			}

			addrs = append(addrs, ifaceAddrs...)
		}
	}

	if len(addrs) == 0 {
		fmt.Println("No addresses")
		return
	}

	var found = 0

	for found == 0 {
		time.Sleep(1 * time.Second)

		fmt.Println("Discovering...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		d, err := StartDiscoveryWithContext(ctx, addrs, time.Second*15)
		if err != nil {
			panic(err)
		}

		found = len(d)
		fmt.Println(found)

		if len(d) == 0 {
			fmt.Println("No cameras were found")
			continue
		}

		for i := 0; i < found; i++ {
			Info("XAddr", d[i].XAddr)

			d[i].User = "admin"
			d[i].Password = "admin"
			nps, err := d[i].GetNetworkProtocols()
			if err != nil {
				Error(err)
			}
			parsed, _ := url.Parse(d[i].XAddr)
			host, _, _ := net.SplitHostPort(parsed.Host)

			for _, np := range nps {
				fmt.Println("Joined", net.JoinHostPort(host, fmt.Sprintf("%d", np.Port)))
			}

			profiles, _ := d[i].GetProfiles()
			for _, p := range profiles {
				fmt.Println(d[i].GetStreamURI(p.Token, "UDP"))
				fmt.Println(d[i].GetStreamURI(p.Token, "RTSP"))
				fmt.Println(d[i].GetStreamURI(p.Token, "HTTP"))
				fmt.Println(d[i].GetSnapshotURI(p.Token))
			}
		}
	}
}
