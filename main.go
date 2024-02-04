package main

import (
	"log"
	"syscall"
)

const BUFFER_SIZE = 100

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
