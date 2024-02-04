package main

import "fmt"

type IPAddress [4]byte

func (ip IPAddress) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

type Protocol byte

func (p Protocol) String() string {
	return map[Protocol]string{
		1:  "ICMP",
		6:  "TCP",
		17: "UDP",
	}[p]
}

type HeaderDefinition struct {
	name    string
	bitSize byte
}

type PDU interface {
	HeadersDefs() []HeaderDefinition
	HeadersValues() []string
}

type Segment struct {
}

type PacketFlags struct {
	dontFragment  bool
	moreFragments bool
}

func (pf PacketFlags) String() string {
	return fmt.Sprintf("0   %d   %d", boolToInt(pf.dontFragment), boolToInt(pf.moreFragments))
}

type Packet struct {
	version            byte // 4 bits
	ihl                byte // 4 bits
	typeOfService      byte
	totalLength        uint16
	identification     uint16
	flags              PacketFlags
	fragmentOffset     uint16 // 13 bits
	timeToLive         byte
	protocol           Protocol
	headerChecksum     uint16
	sourceAddress      IPAddress
	destinationAddress IPAddress
	options            []byte
	segment            Segment
}

func (p Packet) HeadersDefs() []HeaderDefinition {
	return []HeaderDefinition{
		{"Version", 4},
		{"IHL", 4},
		{"TOS", 8},
		{"Total Length", 16},
		{"Identification", 16},
		{"Flags", 3},
		{"Fragment Offset", 13},
		{"TTL", 8},
		{"Protocol", 8},
		{"Header Checksum", 16},
		{"Source Address", 32},
		{"Destination Address", 32},
	}
}

func (p Packet) HeadersValues() []string {
	return []string{
		fmt.Sprintf("IPv%d", p.version),
		fmt.Sprintf("%d", p.ihl),
		fmt.Sprintf("%d", p.typeOfService),
		fmt.Sprintf("%d", p.totalLength),
		fmt.Sprintf("%d", p.identification),
		p.flags.String(),
		fmt.Sprintf("%d", p.fragmentOffset),
		fmt.Sprintf("%d", p.timeToLive),
		p.protocol.String(),
		fmt.Sprintf("%d", p.headerChecksum),
		p.sourceAddress.String(),
		p.destinationAddress.String(),
	}
}
