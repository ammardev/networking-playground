package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"
)

const BUFFER_SIZE = 100

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
		p.protocol.String(),
		fmt.Sprintf("%d", p.ihl),
		fmt.Sprintf("%d", p.typeOfService),
		fmt.Sprintf("%d", p.totalLength),
		fmt.Sprintf("%d", p.identification),
		p.flags.String(),
		fmt.Sprintf("%d", p.fragmentOffset),
		fmt.Sprintf("%d", p.timeToLive),
		fmt.Sprintf("%d", p.typeOfService),
		fmt.Sprintf("%d", p.headerChecksum),
		p.sourceAddress.String(),
		p.destinationAddress.String(),
	}
}

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_UDP)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start listening...")

	for {
		buffer := make([]byte, BUFFER_SIZE)

		_, _, err = syscall.Recvfrom(fd, buffer, 0)

		packet := Packet{
			version:            buffer[0] >> 4,
			ihl:                buffer[0] & 0x0F,
			typeOfService:      buffer[1],
			totalLength:        BytesToUInt16(buffer[2], buffer[3]),
			identification:     BytesToUInt16(buffer[4], buffer[5]),
			flags:              PacketFlags{buffer[6]&0b010 != 0, buffer[6]&0b001 != 0},
			fragmentOffset:     (BytesToUInt16(buffer[6], buffer[7])) & 0b00011111,
			timeToLive:         buffer[8],
			protocol:           Protocol(buffer[9]),
			headerChecksum:     BytesToUInt16(buffer[10], buffer[11]),
			sourceAddress:      IPAddress(buffer[12:16]),
			destinationAddress: IPAddress(buffer[16:21]),
		}

		pduRenderer := newProtocolDataUnitRenderer(packet, 4)

		pduRenderer.Print()
	}
}

func BytesToUInt16(msByte byte, lsByte byte) uint16 {
	return uint16(msByte)<<8 + uint16(lsByte)
}

func boolToInt(b bool) byte {
	if b {
		return 1
	}

	return 0
}

type ProtocolDataUnitRenderer struct {
	pdu     PDU
	width   byte
	builder strings.Builder
}

func newProtocolDataUnitRenderer(pdu PDU, w byte) ProtocolDataUnitRenderer {
	return ProtocolDataUnitRenderer{
		width: w,
		pdu:   pdu,
	}
}

func (r *ProtocolDataUnitRenderer) printHeaders() {
	r.builder.WriteString("  ")

	var i byte
	for i = 0; i < r.width; i++ {
		r.builder.WriteByte(i + 48)

		r.builder.WriteString(
			strings.Repeat(" ", 39),
		)
	}

	r.builder.WriteString("\n  ")

	for i = 0; i < r.width*8; i++ {
		r.builder.WriteByte((i % 10) + 48)

		r.builder.WriteString("   ")
	}

	r.builder.WriteString("\n")
}

func (r *ProtocolDataUnitRenderer) Print() {
	r.printHeaders()

	r.printLineSeparator('=')

	printedHeadersSize := byte(0)

	for _, def := range r.pdu.HeadersDefs() {
		printedHeadersSize += def.bitSize
		cellSize := int(4*def.bitSize - 1)
		r.builder.WriteString("|")

		spaces := cellSize - len(def.name)
		if spaces/2 < 0 {
			spaces = 0
		}

		r.builder.WriteString(strings.Repeat(" ", spaces/2))
		r.builder.WriteString(def.name)
		r.builder.WriteString(strings.Repeat(" ", spaces/2))

		if spaces%2 != 0 {
			r.builder.WriteString(" ")
		}

		if printedHeadersSize/8 == r.width {
			r.builder.WriteString("|\n")
			printedHeadersSize = 0
			r.printLineSeparator('-')
			// Print items
			r.printLineSeparator('=')
		}
	}

	fmt.Print(r.builder.String())
}

func (r *ProtocolDataUnitRenderer) printLineSeparator(char rune) {
	r.builder.WriteString("|")
	r.builder.WriteString(
		strings.Repeat(string(char), 127),
	)
	r.builder.WriteString("|\n")
}
