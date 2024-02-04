package main

import (
	"fmt"
	"strings"
)

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

	var valuesToPrint []string

	for index, def := range r.pdu.HeadersDefs() {
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

		spaces = cellSize - len(r.pdu.HeadersValues()[index])

		val := strings.Repeat(" ", spaces/2) + r.pdu.HeadersValues()[index] + strings.Repeat(" ", spaces/2)

		if spaces%2 != 0 {
			val = val + " "
		}

		valuesToPrint = append(valuesToPrint, val)

		if printedHeadersSize/8 == r.width {
			r.builder.WriteString("|\n")
			r.printLineSeparator('-')
			printedHeadersSize = 0

			for _, value := range valuesToPrint {
				r.builder.WriteString("|")
				r.builder.WriteString(value)
			}

			r.builder.WriteString("|\n")
			r.printLineSeparator('=')
			valuesToPrint = []string{}
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
