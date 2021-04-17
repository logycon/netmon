package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Serve() {
	http.HandleFunc("/devices", handleDevicesRequest)
	http.HandleFunc("/capture", handleCaptureRequest)

	port := fmt.Sprintf(":%d", portNumber)
	log.Printf("Serving..%s.\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleDevicesRequest(rw http.ResponseWriter, req *http.Request) {
	devices := make(chan []Device)

	go CollectDevices(devices)

	res := <-devices

	json, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal("Failed to generate json")
	}

	fmt.Fprintf(rw, "%s\n", string(json))
}

func handleCaptureRequest(rw http.ResponseWriter, req *http.Request) {
	device := req.URL.Query().Get("d")
	if len(device) <= 0 {
		http.NotFound(rw, req)
		return
	}

	filter := req.URL.Query().Get("f")
	if len(filter) <= 0 {
		filter = ""
	}

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.NotFound(rw, req)
		return
	}

	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Content-Type", "application/stream+json; charset=utf-8")

	messages := make(chan Packet)

	go DeviceCapture(device, filter, messages)

	for {
		select {
		case <-cn.CloseNotify():
			log.Println("Client stopped listening")
			return
		case msg := <-messages:
			jsonOutput, err := json.MarshalIndent(msg, "", " ")
			if err != nil {
				log.Printf("could not json marshall reponse item %#v: %v\n", msg, err)
				continue
			}

			// write to output
			_, err = fmt.Fprintf(rw, "%s\n", string(jsonOutput))
			if err != nil {
				log.Fatal(err)
			}
			f.Flush()
		}
	}
}
