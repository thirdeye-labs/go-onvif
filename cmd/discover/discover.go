package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/apex/log"
	onvif "github.com/byronwilliams/go-onvif"
)

func main() {
	log.SetLevelFromString("debug")
	log.Debug("Starting")

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var addrs []net.Addr

	for _, iface := range ifaces {
		if iface.Name == "wlp58s0" || iface.Name == "wlan0" || iface.Name == "wlan1" {
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*35)
		defer cancel()
		d, err := onvif.StartDiscoveryWithContext(ctx, addrs, time.Second*30)
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
			fmt.Println(d[i].XAddr)
			nps, _ := d[i].GetNetworkProtocols()

			parsed, _ := url.Parse(d[i].XAddr)
			host, _, _ := net.SplitHostPort(parsed.Host)

			for _, np := range nps {
				fmt.Println(net.JoinHostPort(host, fmt.Sprintf("%d", np.Port)))
			}

			profiles, _ := d[i].GetProfiles()
			for _, p := range profiles {
				fmt.Println(d[i].GetStreamURI(p.Token, "UDP"))
				fmt.Println(d[i].GetStreamURI(p.Token, "RTSP"))
				fmt.Println(d[i].GetStreamURI(p.Token, "HTTP"))
			}
		}
	}
}
