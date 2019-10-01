//-----------------------------------------------------------------------------
/*

6502 Disassembler

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"

	"github.com/deadsy/bender/cpu"
)

//-----------------------------------------------------------------------------

func memString(adr uint16, mem []byte) string {
	s := make([]string, len(mem))
	for i, v := range mem {
		s[i] = fmt.Sprintf("%02x", v)
	}
	return fmt.Sprintf("%04x: %-15s", uint16(adr), strings.Join(s, " "))
}

//-----------------------------------------------------------------------------

func main() {

	code := []byte{
		0xa9, 0x00, 0xa8, 0x91, 0x02, 0xc8, 0xca, 0xd0,
		0xfa, 0x60, 0x18, 0x36, 0x00, 0x88, 0xd0, 0x01,
		0x60, 0xe8, 0x4c, 0x0b, 0x02, 0x18, 0x76, 0x00,
		0x88, 0xd0, 0x01, 0x60, 0xca, 0x4c, 0x16, 0x02,
		0x38, 0xa9, 0xff, 0x55, 0x00, 0x69, 0x00, 0x95,
		0x00, 0xe8, 0x88, 0xd0, 0xf4, 0x60, 0x18, 0xa0,
		0x00, 0xb1, 0x02, 0x71, 0x00, 0x91, 0x02, 0xc8,
		0xca, 0xd0, 0xf6, 0x60, 0x38, 0xa0, 0x00, 0xb1,
		0xff, 0xff, 0xff,
	}

	adr := uint16(0x200)
	var ofs int

	for ofs < len(code) {
		sIns, n := cpu.Disassemble(adr, code[ofs:])
		sMem := memString(adr, code[ofs:ofs+n])
		fmt.Printf("%s %s\n", sMem, sIns)
		ofs += n
		adr += uint16(n)
	}

}

//-----------------------------------------------------------------------------
