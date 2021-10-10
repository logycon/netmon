package main

import (
	"flag"
	"fmt"
	interactive2 "github.com/logycon/net/netmon/interactive"
	"github.com/logycon/net/netmon/server"
	"os"
)

var (
	interactive   bool
	listDevices   bool
	captureDevice int
	captureHttp   int

	serve      bool
	portNumber int

	testPrintPackets bool
	testFilter       string
	testPrintPayload bool
	helpFlag         bool
	outputFile       string
	inputFile        string
)

func main() {
	flag.BoolVar(&helpFlag, "h", false, "Help (this screen)")
	flag.BoolVar(&serve, "s", false, "Http server mode, e.g. -s")
	flag.IntVar(&portNumber, "p", 8080, "Server's port number, e.g. -p=12345")

	flag.BoolVar(&interactive, "i", false, "Interactive mode, e.g. -i")
	flag.BoolVar(&listDevices, "ld", false, "List Devices, e.g. -i -ld")
	flag.IntVar(&captureDevice, "cd", 0, "Capture Device by index, e.g. -i -cd=1")
	flag.IntVar(&captureHttp, "ch", 0, "Capture Http from device by index, e.g. -i ch=1")

	flag.BoolVar(&testPrintPackets, "tpp", false, "Test Mode: Also print gopacket, e.g. -tpp")
	flag.BoolVar(&testPrintPayload, "tp", true, "Test Mode: Print payload, e.g. -tp=false")
	flag.StringVar(&testFilter, "tf", "", "Test Mode: BPF, e.g. -tf=\"host 192.168.1.1\"")
	flag.StringVar(&outputFile, "of", "", "Output pcap file name, dirtest mode only, e.g. -t -o=out.pcap")
	flag.StringVar(&inputFile, "if", "", "Input pcap file to play, e.g. -t -i=in.pcap")

	flag.Parse()

	if helpFlag {
		fmt.Println("\nNetwork Monitor 0.1\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if interactive {
		switch {
		case listDevices:
			interactive2.ListDevices()
		case captureDevice > 0:
			interactive2.CaptureDeviceByIndex(captureDevice)
		case captureHttp > 0:
			interactive2.CaptureHttpByIndex(captureHttp)
		}

		// dirtest()

	}

	if serve {
		server.Serve(portNumber)
	}
}

func test() {
	if len(inputFile) > 0 {
		HandleInput(inputFile, testPrintPackets, testPrintPayload)
	} else {
		HandleTest(testPrintPackets, testFilter, testPrintPayload, outputFile)
	}
}
