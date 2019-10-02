//-----------------------------------------------------------------------------
/*

6502 CPU Emulator

See also:

https://github.com/redcode/6502
https://www.masswerk.at/6502/6502_instruction_set.html

*/
//-----------------------------------------------------------------------------

package cpu

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

// Code Segment defines a contiguous area of memory.
type CodeSegment struct {
	Base   uint16            // base address
	Memory []byte            // memory content
	Symbol map[uint16]string // symbol table
}

// Disassembly returns the result of the disassembler.
type Disassembly struct {
	Memory      string
	Symbol      string
	Instruction string
	Comment     string
	n           uint // length of decode
}

//-----------------------------------------------------------------------------

func (cs *CodeSegment) getMem(adr uint16, n uint) ([]byte, error) {
	ofs := int(adr) - int(cs.Base)
	if ofs < 0 || ofs >= len(cs.Memory) {
		return nil, fmt.Errorf("address is out of segment")
	}
	if ofs+int(n) >= len(cs.Memory) {
		return nil, fmt.Errorf("length is out of segment")
	}
	return cs.Memory[ofs : ofs+int(n)], nil
}

func (cs *CodeSegment) daMemory(adr uint16, mem []byte) string {
	s := make([]string, len(mem))
	for i, v := range mem {
		s[i] = fmt.Sprintf("%02x", v)
	}
	return fmt.Sprintf("%04x: %s", uint16(adr), strings.Join(s, " "))
}

func (cs *CodeSegment) daSymbol(adr uint16) string {
	if cs.Symbol != nil {
		return cs.Symbol[adr]
	}
	return ""
}

func (cs *CodeSegment) daInstruction(adr uint16, mem []byte) string {

	var s []string

	info := opcodeLookup(mem[0])
	ofs := 1

	// instruction mneumonic
	s = append(s, info.ins)

	switch info.mode {
	case amNone:
		// illegal - no operands
	case amAcc:
		// accumulator - no operands
		s = append(s, "a")
	case amAbs:
		// absolute - 2 byte operand
		operand := int(mem[ofs]) + (int(mem[ofs+1]) << 8)
		s = append(s, fmt.Sprintf("$%04x", operand))
	case amAbsX:
		s = append(s, "TODO absolute, X-indexed")
	case amAbsY:
		s = append(s, "TODO absolute, Y-indexed")
	case amImm:
		// immediate - 1 byte operand
		operand := mem[ofs]
		s = append(s, fmt.Sprintf("#$%02x", operand))
	case amImpl:
		// implied - no operands
	case amInd:
		s = append(s, "TODO indirect")
	case amXInd:
		s = append(s, "TODO X-indexed, indirect")
	case amIndY:
		// indirect, Y-indexed - 1 byte operand
		operand := mem[ofs]
		s = append(s, fmt.Sprintf("($%02x),y", operand))
	case amRel:
		// relative - 1 byte operand
		operand := mem[ofs]
		dst := uint16(int(adr) + int(int8(operand)) + 2)
		s = append(s, fmt.Sprintf("$%02x\t\t; $%04x", operand, dst))
	case amZpg:
		// zeropage - 1 byte operand
		operand := mem[ofs]
		s = append(s, fmt.Sprintf("$%02x", operand))
	case amZpgX:
		// zeropage, X-indexed - 1 byte operand
		operand := mem[ofs]
		s = append(s, fmt.Sprintf("$%02x,x", operand))
	case amZpgY:
		s = append(s, "TODO zeropage, Y-indexed")
	default:
		panic("bad address mode")
	}

	return strings.Join(s, " ")
}

func (cs *CodeSegment) Disassemble(adr uint16) (*Disassembly, error) {

	// get the instruction memory
	mem, err := cs.getMem(adr, 1)
	if err != nil {
		return nil, err
	}
	mem, err = cs.getMem(adr, insLength(mem[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

//-----------------------------------------------------------------------------
