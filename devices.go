package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func CollectDevices(out chan []Device) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var res []Device

	for _, device := range devices {
		d := Device{
			Name:        device.Name,
			Description: device.Description,
			Addresses:   make([]Address, 0),
		}

		for _, address := range device.Addresses {
			a := Address{
				IP:      address.IP.String(),
				Netmask: getNetmask(address),
			}
			if a.IP != "0.0.0.0" && a.Netmask != "" {
				d.Addresses = append(d.Addresses, a)
			}
		}
		if len(d.Addresses) > 0 {
			res = append(res, d)
		}
	}

	out <- res
}

func getNetmask(addr pcap.InterfaceAddress) string {
	if addr.Netmask != nil {
		mask := addr.Netmask
		return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
	} else {
		return ""
	}
}
