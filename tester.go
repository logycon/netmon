package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

var (
	deviceName = "\\\\Device\\\\NPF_{8C8CFC6A-A95C-4203-9F9F-EDB615E0D2F0}"
)

func HandleInput(inputFile string, printPackets bool, printPayload bool) {

	handle, err := pcap.OpenOffline(inputFile)
	if err != nil {
		fmt.Printf("Error opening %s :%#v\n", inputFile, err)
		os.Exit(1)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacket(packet, printPackets, printPayload)
	}
}

func HandleTest(printPackets bool, filter string, printPayload bool, outputFile string) {
	var timeout = -1 * time.Second
	var handle *pcap.Handle
	var w *pcapgo.Writer

	if len(outputFile) > 0 {
		f, _ := os.Create(outputFile)
		w = pcapgo.NewWriter(f)
		err := w.WriteFileHeader(65535, layers.LinkTypeEthernet)
		if err != nil {
			outputFile = ""
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Printf("Error closing %s :%#v\n", outputFile, err)
			}
		}(f)
	}

	handle, err := pcap.OpenLive(deviceName, 65535, true, timeout)
	if err != nil {
		log.Printf("Error opening device %s: %v", deviceName, err)
		os.Exit(1)
	}

	if len(filter) > 0 {
		if err := handle.SetBPFFilter(filter); err != nil {
			log.Printf("Error setting filter: %s\n", filter)
		} else {
			fmt.Printf("Filtering on %s\n", filter)
		}
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		if len(outputFile) > 0 {
			err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			if err != nil {
				fmt.Printf("Error writing to %s :%#v\n", outputFile, err)
			}
		}
		printPacket(packet, printPackets, printPayload)
	}
}

func printPacket(packet gopacket.Packet, printPackets bool, printPayload bool) {

	var fromPort string
	var toPort string

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		fromPort = fmt.Sprintf(":%d", tcp.SrcPort)
		toPort = fmt.Sprintf(":%d", tcp.DstPort)
	} else {
		fromPort = ""
		toPort = ""
	}

	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		ip4, _ := ip4Layer.(*layers.IPv4)
		fmt.Printf("--- %s: %s From %s%s to %s%s  ---\n",
			packet.Metadata().Timestamp, ip4.Protocol, ip4.SrcIP, fromPort, ip4.DstIP, toPort,
		)
	} else {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer != nil {
			arp, _ := arpLayer.(*layers.ARP)
			fmt.Printf("--- %s: %s From %s%s to %s%s ---\n",
				packet.Metadata().Timestamp, arp.Protocol, arp.SourceProtAddress, fromPort, arp.DstProtAddress, toPort,
			)
		}
	}

	if printPayload {
		fmt.Println("---- Payload Begin ----\n")
		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			fmt.Print(string(appLayer.Payload()))
		}
		fmt.Println("\n---- End Payload ----\n")
	}

	if printPackets {
		fmt.Println("----- GoPacket Begin -----\n")
		fmt.Println(packet)
		fmt.Println("----- End GoPacket -----\n")
	}
}
