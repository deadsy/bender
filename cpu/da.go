//-----------------------------------------------------------------------------
/*

6502 CPU Disassembler

*/
//-----------------------------------------------------------------------------

package cpu

import (
	"fmt"
	"strings"
)

//-----------------------------------------------------------------------------

type SymbolTable map[uint16]string

// Disassembly returns the result of the disassembler call.
type Disassembly struct {
	Dump        string  // address and memory bytes
	Symbol      string  // symbol for the address (if any)
	Instruction string  // instruction decode
	Comment     string  // useful comment
	Bytes       []uint8 // decoded bytes
}

func (da *Disassembly) String() string {
	s := make([]string, 2)
	s[0] = fmt.Sprintf("%-16s %8s %-13s", da.Dump, da.Symbol, da.Instruction)
	if da.Comment != "" {
		s[1] = fmt.Sprintf(" ; %s", da.Comment)
	}
	return strings.Join(s, "")
}

//-----------------------------------------------------------------------------

func daDump(adr uint16, mem []byte) string {
	s := make([]string, len(mem))
	for i, v := range mem {
		s[i] = fmt.Sprintf("%02x", v)
	}
	return fmt.Sprintf("%04x: %s", adr, strings.Join(s, " "))
}

func daSymbol(adr uint16, st SymbolTable) string {
	if st != nil {
		return st[adr]
	}
	return ""
}

func daInstruction(adr uint16, mem []uint8) (string, string) {

	var s []string
	var comment string

	info := opcodeLookup(mem[0])

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
		operand := int(mem[1]) + (int(mem[2]) << 8)
		s = append(s, fmt.Sprintf("$%04x", operand))
	case amAbsX:
		s = append(s, "TODO absolute, X-indexed")
	case amAbsY:
		s = append(s, "TODO absolute, Y-indexed")
	case amImm:
		// immediate - 1 byte operand
		operand := mem[1]
		s = append(s, fmt.Sprintf("#$%02x", operand))
	case amImpl:
		// implied - no operands
	case amInd:
		s = append(s, "TODO indirect")
	case amXInd:
		s = append(s, "TODO X-indexed, indirect")
	case amIndY:
		// indirect, Y-indexed - 1 byte operand
		operand := mem[1]
		s = append(s, fmt.Sprintf("($%02x),y", operand))
	case amRel:
		// relative - 1 byte operand
		operand := mem[1]
		s = append(s, fmt.Sprintf("$%02x", operand))
		dst := uint16(int(adr) + int(int8(operand)) + 2)
		comment = fmt.Sprintf("$%04x", dst)
	case amZpg:
		// zeropage - 1 byte operand
		operand := mem[1]
		s = append(s, fmt.Sprintf("$%02x", operand))
	case amZpgX:
		// zeropage, X-indexed - 1 byte operand
		operand := mem[1]
		s = append(s, fmt.Sprintf("$%02x,x", operand))
	case amZpgY:
		s = append(s, "TODO zeropage, Y-indexed")
	default:
		panic("bad address mode")
	}

	return strings.Join(s, " "), comment
}

// Disassemble a 6502 instruction from the memory at the address.
func Disassemble(m Memory, adr uint16, st SymbolTable) *Disassembly {
	// get the instruction bytes
	mem := make([]uint8, insLength(m.Read8(adr)))
	for i := range mem {
		mem[i] = m.Read8(adr + uint16(i))
	}

	instruction, comment := daInstruction(adr, mem)

	return &Disassembly{
		Dump:        daDump(adr, mem),
		Symbol:      daSymbol(adr, st),
		Instruction: instruction,
		Comment:     comment,
		Bytes:       mem,
	}
}

//-----------------------------------------------------------------------------

// Disassemble returns the disassembly for a region of the CPU memory.
func (m *M6502) Disassemble(adr uint16, size int) string {
	s := make([]string, 0, 16)
	for size > 0 {
		da := Disassemble(m.mem, adr, nil)
		s = append(s, da.String())
		n := len(da.Bytes)
		size -= n
		adr += uint16(n)
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
