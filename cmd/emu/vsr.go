//-----------------------------------------------------------------------------
/*

6502 Emulator: Virtual Subroutines

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/deadsy/bender/cpu"
)

//-----------------------------------------------------------------------------

func getAX(m *cpu.M6502) uint16 {
	return uint16(m.A) + uint16(m.X<<8)
}

func setAX(m *cpu.M6502, val uint16) {
	m.A = uint8(val)
	m.X = uint8(val >> 8)
}

const spAddress = 0

func popParam(m *cpu.M6502, inc uint16) uint16 {
	mem := m.Mem.(*memory)
	sp := mem.read16(spAddress)
	val := mem.read16(sp)
	mem.write16(spAddress, sp+inc)
	return val
}

//-----------------------------------------------------------------------------

func vsrOpen(m *cpu.M6502) {
	fmt.Printf("*** vsrOpen ***\n")
}
func vsrClose(m *cpu.M6502) {
	fmt.Printf("*** vsrClose ***\n")
}
func vsrRead(m *cpu.M6502) {
	fmt.Printf("*** vsrRead ***\n")
}
func vsrWrite(m *cpu.M6502) {

	n := getAX(m)
	buf := popParam(m, 2)
	fd := popParam(m, 2)

	fmt.Printf("vsrWrite ($%04X, $%04X, $%04X)\n", fd, buf, n)

	// 0: PVWrite ($0001, $03C7, $000C)

	//Data = xmalloc (Count);
	//while (I < Count) {
	//    Data[I++] = MemReadByte (Buf++);
	//}

	//RetVal = write (FD, Data, Count);

	setAX(m, 0)

}
func vsrArgs(m *cpu.M6502) {
	fmt.Printf("*** vsrArgs ***\n")
}
func vsrExit(m *cpu.M6502) {
	fmt.Printf("*** vsrExit ***\n")
}

//-----------------------------------------------------------------------------