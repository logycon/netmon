package capture

import (
	"github.com/logycon/net/netmon/common"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func DeviceCapture(deviceName string, filter string, out chan common.Packet) {
	var snapshotLen int32 = 65535
	var promiscuous bool = true
	var timeout time.Duration = -1 * time.Second
	var handle *pcap.Handle

	handle, err := pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Printf("Error setting filter: %s\n", filter)
		filter = "*none*"
		//os.Exit(1)
	}

	// log.Printf("Capturing %s with filter %s\n", deviceName, filter)

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		out <- packetToMessage(packet)
	}
}

func packetToMessage(packet gopacket.Packet) common.Packet {
	v := common.Packet{}
	v.Timestamp = packet.Metadata().Timestamp

	for _, layer := range packet.Layers() {
		v.Layers = append(v.Layers, strings.ToLower(layer.LayerType().String()))
	}

	if ethLayer := packet.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		eth, _ := ethLayer.(*layers.Ethernet)
		v.Ethernet.Src = eth.SrcMAC.String()
		v.Ethernet.Dest = eth.DstMAC.String()
		v.Ethernet.PayloadSize = len(eth.Payload)
	}

	if common.Contains(v.Layers, "tcp") {
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)

			v.TCP = &common.TCP{
				FromPort: int(tcp.SrcPort),
				ToPort:   int(tcp.DstPort),
				Flags: common.TCPFlags{
					ACK: strconv.FormatBool(tcp.ACK),
					FIN: strconv.FormatBool(tcp.FIN),
					URG: strconv.FormatBool(tcp.URG),
					PSH: strconv.FormatBool(tcp.PSH),
					RST: strconv.FormatBool(tcp.RST),
					ECE: strconv.FormatBool(tcp.ECE),
					CWR: strconv.FormatBool(tcp.CWR),
					NS:  strconv.FormatBool(tcp.NS),
				},
			}
		}
	}

	if common.Contains(v.Layers, "arp") {
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp, _ := arpLayer.(*layers.ARP)

			v.ARP = &common.ARP{
				Operation: common.IIF(arp.Operation == layers.ARPRequest, "ArpRequest", "ArpReply"),
				Sender: common.ARPParticipant{
					MAC: net.HardwareAddr(arp.SourceHwAddress).String(),
					IP:  net.IP(arp.SourceProtAddress).String(),
				},
				Target: common.ARPParticipant{
					MAC: net.HardwareAddr(arp.DstHwAddress).String(),
					IP:  net.IP(arp.DstProtAddress).String(),
				},
			}
		}
	}

	if common.Contains(v.Layers, "udp") {
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			v.UDP = &common.UDP{
				SrcPort: int(udp.SrcPort),
				DstPort: int(udp.DstPort),
			}
		}
	}

	if common.Contains(v.Layers, "ipv4") {
		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			ipv4, _ := ipv4Layer.(*layers.IPv4)

			//ipv4.Flags

			v.IPV4 = &common.IPV4{
				Src:      ipv4.SrcIP.String(),
				Dst:      ipv4.DstIP.String(),
				Protocol: ipv4.Protocol.String(),
				ID:       int(ipv4.Id),
				Offset:   int(ipv4.FragOffset),
				Flags:    ipv4.Flags.String(),
				TTL:      int(ipv4.TTL),
			}
		}
	}

	if common.Contains(v.Layers, "payload") {
		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			payload := string(appLayer.Payload())
			v.Payload = &payload
		}
	}

	return v
}
