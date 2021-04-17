package main

import "time"

type Address struct {
	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
}

type Device struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Addresses   []Address `json:"addresses"`
}

type Ethernet struct {
	Src         string `json:"from"`
	Dest        string `json:"to"`
	PayloadSize int    `json:"payloadSize"`
}

type TCPFlags struct {
	ACK string `json:"ack,omitempty"`
	FIN string `json:"fin,omitempty"`
	URG string `json:"urg,omitempty"`
	PSH string `json:"psh,omitempty"`
	RST string `json:"rst,omitempty"`
	ECE string `json:"ece,omitempty"`
	CWR string `json:"cwr,omitempty"`
	NS  string `json:"ns,omitempty"`
}

type TCP struct {
	FromPort int      `json:"fromPort,omitempty"`
	ToPort   int      `json:"toPort,omitempty"`
	Flags    TCPFlags `json:"flags,omitempty"`
}

type ARP struct {
	Operation string         `json:"operation,omitempty"`
	Sender    ARPParticipant `json:"sender,omitempty"`
	Target    ARPParticipant `json:"target,omitempty"`
}

type ARPParticipant struct {
	MAC string `json:"mac,omitempty"`
	IP  string `json:"ip,omitempty"`
}

type UDP struct {
	SrcPort int `json:"srcPort,omitempty"`
	DstPort int `json:"dstPort,omitempty"`
}

type IPV4 struct {
	Src      string `json:"fromIP,omitempty"`
	Dst      string `json:"toIP,omitempty"`
	ID       int    `json:"id,omitempty"`
	Offset   int    `json:"offset,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Flags    string `json:"flags,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
}

type Packet struct {
	Timestamp time.Time `json:"timestamp"`
	Layers    []string  `json:"layers,omitempty"`
	Ethernet  Ethernet  `json:"ethernet,omitempty"`
	TCP       *TCP      `json:"tcp,omitempty"`
	ARP       *ARP      `json:"arp,omitempty"`
	UDP       *UDP      `json:"udp,omitempty"`
	IPV4      *IPV4     `json:"ipv4,omitempty"`
	Payload   *string   `json:"payload,omitempty"`
}
