package main

import (
	"encoding/json"
	"fmt"
	"sort"
)

func ListDevices() {
	devices := make(chan []Device)

	go CollectDevices(devices)

	res := <-devices

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})

	for index, r := range res {
		fmt.Printf("%d. %s\t %s %s\n", index+1, r.Name, r.Description, r.Addresses)
	}
}

func CaptureDeviceByIndex(index int) {
	devices := make(chan []Device)
	go CollectDevices(devices)
	res := <-devices
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})

	if index <= 0 || index > len(res) {
		fmt.Printf("Device with index %d not found.\n", index)
		return
	}

	device := res[index-1]
	messages := make(chan Packet)

	go DeviceCapture(device.Name, "", messages)

	for {
		select {
		case msg := <-messages:
			jsonOutput, err := json.MarshalIndent(msg, "", " ")
			if err != nil {
				fmt.Printf("could not json marshall reponse item %#v: %v\n", msg, err)
				continue
			}
			fmt.Println(string(jsonOutput))
		}
	}
}

func CaptureHttpByIndex(index int) {
	devices := make(chan []Device)
	go CollectDevices(devices)
	res := <-devices
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name > res[j].Name
	})

	if index <= 0 || index > len(res) {
		fmt.Printf("Device with index %d not found.\n", index)
		return
	}

	device := res[index-1]
	//	messages := make(chan string)

	HttpCapture(device.Name, "") // , messages)
}
