package main

import (
	"fmt"
	"log"
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

type Segment struct {
}

type Packet struct {
	version            byte // 4 bits
	ihl                byte // 4 bits
	typeOfService      byte
	totalLength        uint16
	identification     uint16
	flags              byte   // 3 bits
	fragmentOffset     uint16 // 13 bits
	timeToLive         byte
	protocol           Protocol
	headerChecksum     uint16
	sourceAddress      IPAddress
	destinationAddress IPAddress
	options            []byte
	segment            Segment
}

func (p Packet) getVersion() string {
	return fmt.Sprintf("IPv%d", p.version)
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
			flags:              buffer[6] >> 5,
			fragmentOffset:     (BytesToUInt16(buffer[6], buffer[7])) & 0b00011111,
			timeToLive:         buffer[8],
			protocol:           Protocol(buffer[9]),
			headerChecksum:     BytesToUInt16(buffer[10], buffer[11]),
			sourceAddress:      IPAddress(buffer[12:16]),
			destinationAddress: IPAddress(buffer[16:21]),
		}

		fmt.Println("    0                   1                   2                   3")
		fmt.Println("    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1")
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |Version|  IHL  |Type of Service|          Total Length         |")
		fmt.Printf("   | %s  |   %d   |       %d       |               %d              |\n", packet.getVersion(), packet.ihl, packet.typeOfService, packet.totalLength)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |         Identification        |Flags|      Fragment Offset    |")
		fmt.Printf("   |            %d              |  %d  |            %d            |\n", packet.identification, packet.flags, packet.fragmentOffset)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |  Time to Live |    Protocol   |         Header Checksum       |")
		fmt.Printf("   | %d          |       %s       |               %d              |\n", packet.timeToLive, packet.protocol, packet.headerChecksum)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |                       Source Address                          |")
		fmt.Printf("   |                       %s                                     |\n", packet.sourceAddress)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |                   Destination Address                        |")
		fmt.Printf("   |                       %s                                     |\n", packet.destinationAddress)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")

		fmt.Printf("%+v", packet)
	}
}

func BytesToUInt16(msByte byte, lsByte byte) uint16 {
	return uint16(msByte)<<8 + uint16(lsByte)
}
