package main

import (
	"fmt"
	"log"
	"syscall"
)

const BUFFER_SIZE = 100

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
	protocol           byte
	headerChecksum     uint16
	sourceAddress      [4]byte
	destinationAddress [4]byte
	options            []byte
	segment            Segment
}

func (p Packet) getVersion() string {
	return fmt.Sprintf("IPv%d", p.version)
}

func (p Packet) getSourceAddress() string {
	return fmt.Sprintf("%d.%d.%d.%d", p.sourceAddress[0], p.sourceAddress[1], p.sourceAddress[2], p.sourceAddress[3])
}

func (p Packet) getDestinationAddress() string {
	return fmt.Sprintf("%d.%d.%d.%d", p.destinationAddress[0], p.destinationAddress[1], p.destinationAddress[2], p.destinationAddress[3])
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
			totalLength:        uint16(buffer[2])<<8 + uint16(buffer[3]),
			identification:     uint16(buffer[4])<<8 + uint16(buffer[5]),
			flags:              buffer[6] >> 5,
			fragmentOffset:     (uint16(buffer[6])<<8 + uint16(buffer[7])) & 0b00011111,
			timeToLive:         buffer[8],
			protocol:           buffer[9],
			headerChecksum:     uint16(buffer[10])<<8 + uint16(buffer[11]),
			sourceAddress:      [4]byte{buffer[12], buffer[13], buffer[14], buffer[15]},
			destinationAddress: [4]byte{buffer[16], buffer[17], buffer[18], buffer[19]},
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
		fmt.Printf("   | %d          |       %d       |               %d              |\n", packet.timeToLive, packet.protocol, packet.headerChecksum)
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |                       Source Address                          |")
		fmt.Printf("   |                       %s                                     |\n", packet.getSourceAddress())
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")
		fmt.Println("   |                   Destination Address                        |")
		fmt.Printf("   |                       %s                                     |\n", packet.getDestinationAddress())
		fmt.Println("   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+")

		fmt.Printf("%+v", packet)
	}
}
