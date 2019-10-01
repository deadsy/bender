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
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------
// address modes

type adrMode int

const (
	amNone adrMode = iota
	amAcc          // accumulator
	amAbs          // absolute
	amAbsX         // absolute, X-indexed
	amAbsY         // absolute, Y-indexed
	amImm          // immediate
	amImpl         // implied
	amInd          // indirect
	amXInd         // X-indexed, indirect
	amIndY         // indirect, Y-indexed
	amRel          // relative
	amZpg          // zeropage
	amZpgX         // zeropage, X-indexed
	amZpgY         // zeropage, Y-indexed
)

type adrModeInfo struct {
	suffix string // function suffix
	descr  string // mode description
}

var modeDescr = map[adrMode]adrModeInfo{
	amNone: {"", ""},
	amAcc:  {"acc", "accumulator"},
	amAbs:  {"abs", "absolute"},
	amAbsX: {"absx", "absolute, X-indexed"},
	amAbsY: {"absy", "absolute, Y-indexed"},
	amImm:  {"imm", "immediate"},
	amImpl: {"impl", "implied"},
	amInd:  {"ind", "indirect"},
	amXInd: {"xind", "X-indexed, indirect"},
	amIndY: {"indy", "indirect, Y-indexed"},
	amRel:  {"rel", "relative"},
	amZpg:  {"z", "zeropage"},
	amZpgX: {"zx", "zeropage, X-indexed"},
	amZpgY: {"zy", "zeropage, Y-indexed"},
}

//-----------------------------------------------------------------------------
// opcodes

// instruction information
type insInfo struct {
	ins  string  // mneumonic
	mode adrMode // address mode
}

// opcode table as a map, unspecified opcodes are illegal instructions
var opcodeTable = map[uint8]insInfo{

	0x00: insInfo{"BRK", amImpl},
	0x10: insInfo{"BPL", amRel},
	0x20: insInfo{"JSR", amAbs},
	0x30: insInfo{"BMI", amRel},
	0x40: insInfo{"RTI", amImpl},
	0x50: insInfo{"BVC", amRel},
	0x60: insInfo{"RTS", amImpl},
	0x70: insInfo{"BVS", amRel},
	0x80: insInfo{"ILL", amNone},
	0x90: insInfo{"BCC", amRel},
	0xa0: insInfo{"LDY", amImm},
	0xb0: insInfo{"BCS", amRel},
	0xc0: insInfo{"CPY", amImm},
	0xd0: insInfo{"BNE", amRel},
	0xe0: insInfo{"CPX", amImm},
	0xf0: insInfo{"BEQ", amImm},

	0x01: insInfo{"ORA", amXInd},
	0x11: insInfo{"ORA", amIndY},
	0x21: insInfo{"AND", amXInd},
	0x31: insInfo{"AND", amIndY},
	0x41: insInfo{"EOR", amXInd},
	0x51: insInfo{"EOR", amIndY},
	0x61: insInfo{"ADC", amXInd},
	0x71: insInfo{"ADC", amIndY},
	0x81: insInfo{"STA", amXInd},
	0x91: insInfo{"STA", amIndY},
	0xa1: insInfo{"LDA", amXInd},
	0xb1: insInfo{"LDA", amIndY},
	0xc1: insInfo{"CMP", amXInd},
	0xd1: insInfo{"CMP", amIndY},
	0xe1: insInfo{"SBC", amXInd},
	0xf1: insInfo{"SBC", amIndY},

	0xa2: insInfo{"LDX", amImm},

	0x04: insInfo{"ILL", amNone},
	0x14: insInfo{"ILL", amNone},
	0x24: insInfo{"BIT", amZpg},
	0x34: insInfo{"ILL", amNone},
	0x44: insInfo{"ILL", amNone},
	0x54: insInfo{"ILL", amNone},
	0x64: insInfo{"ILL", amNone},
	0x74: insInfo{"ILL", amNone},
	0x84: insInfo{"STY", amZpg},
	0x94: insInfo{"STY", amZpgX},
	0xa4: insInfo{"LDY", amZpg},
	0xb4: insInfo{"LDY", amZpgX},
	0xc4: insInfo{"CPY", amZpg},
	0xd4: insInfo{"ILL", amNone},
	0xe4: insInfo{"CPX", amZpg},
	0xf4: insInfo{"ILL", amNone},

	0x05: insInfo{"ORA", amZpg},
	0x15: insInfo{"ORA", amZpgX},
	0x25: insInfo{"AND", amZpg},
	0x35: insInfo{"AND", amZpgX},
	0x45: insInfo{"EOR", amZpg},
	0x55: insInfo{"EOR", amZpgX},
	0x65: insInfo{"ADC", amZpg},
	0x75: insInfo{"ADC", amZpgX},
	0x85: insInfo{"STA", amZpg},
	0x95: insInfo{"STA", amZpgX},
	0xa5: insInfo{"LDA", amZpg},
	0xb5: insInfo{"LDA", amZpgX},
	0xc5: insInfo{"CMP", amZpg},
	0xd5: insInfo{"CMP", amZpgX},
	0xe5: insInfo{"SBC", amZpg},
	0xf5: insInfo{"SBC", amZpgX},

	0x06: insInfo{"ASL", amZpg},
	0x16: insInfo{"ASL", amZpgX},
	0x26: insInfo{"ROL", amZpg},
	0x36: insInfo{"ROL", amZpgX},
	0x46: insInfo{"LSR", amZpg},
	0x56: insInfo{"LSR", amZpgX},
	0x66: insInfo{"ROR", amZpg},
	0x76: insInfo{"ROR", amZpgX},
	0x86: insInfo{"STX", amZpg},
	0x96: insInfo{"STX", amZpgY},
	0xa6: insInfo{"LDX", amZpg},
	0xb6: insInfo{"LDX", amZpgY},
	0xc6: insInfo{"DEC", amZpg},
	0xd6: insInfo{"DEC", amZpgX},
	0xe6: insInfo{"INC", amZpg},
	0xf6: insInfo{"INC", amZpgX},

	0x08: insInfo{"PHP", amImpl},
	0x18: insInfo{"CLC", amImpl},
	0x28: insInfo{"PLP", amImpl},
	0x38: insInfo{"SEC", amImpl},
	0x48: insInfo{"PHA", amImpl},
	0x58: insInfo{"CLI", amImpl},
	0x68: insInfo{"PLA", amImpl},
	0x78: insInfo{"SEI", amImpl},
	0x88: insInfo{"DEY", amImpl},
	0x98: insInfo{"TYA", amImpl},
	0xa8: insInfo{"TAY", amImpl},
	0xb8: insInfo{"CLV", amImpl},
	0xc8: insInfo{"INY", amImpl},
	0xd8: insInfo{"CLD", amImpl},
	0xe8: insInfo{"INX", amImpl},
	0xf8: insInfo{"SED", amImpl},

	0x09: insInfo{"ORA", amImm},
	0x19: insInfo{"ORA", amAbsY},
	0x29: insInfo{"AND", amImm},
	0x39: insInfo{"AND", amAbsY},
	0x49: insInfo{"EOR", amImm},
	0x59: insInfo{"EOR", amAbsY},
	0x69: insInfo{"ADC", amImm},
	0x79: insInfo{"ADC", amAbsY},
	0x89: insInfo{"ILL", amNone},
	0x99: insInfo{"STA", amAbsY},
	0xa9: insInfo{"LDA", amImm},
	0xb9: insInfo{"LDA", amAbsY},
	0xc9: insInfo{"CMP", amImm},
	0xd9: insInfo{"CMP", amAbsY},
	0xe9: insInfo{"SBC", amImm},
	0xf9: insInfo{"SBC", amAbsY},

	0x0a: insInfo{"ASL", amAcc},
	0x1a: insInfo{"ILL", amNone},
	0x2a: insInfo{"ROL", amAcc},
	0x3a: insInfo{"ILL", amNone},
	0x4a: insInfo{"LSR", amAcc},
	0x5a: insInfo{"ILL", amNone},
	0x6a: insInfo{"ROR", amAcc},
	0x7a: insInfo{"ILL", amNone},
	0x8a: insInfo{"TXA", amImpl},
	0x9a: insInfo{"TXS", amImpl},
	0xaa: insInfo{"TAX", amImpl},
	0xba: insInfo{"TSX", amImpl},
	0xca: insInfo{"DEX", amImpl},
	0xda: insInfo{"ILL", amNone},
	0xea: insInfo{"NOP", amImpl},
	0xfa: insInfo{"ILL", amNone},

	0x0c: insInfo{"ILL", amNone},
	0x1c: insInfo{"ILL", amNone},
	0x2c: insInfo{"BIT", amAbs},
	0x3c: insInfo{"ILL", amNone},
	0x4c: insInfo{"JMP", amAbs},
	0x5c: insInfo{"ILL", amNone},
	0x6c: insInfo{"JMP", amInd},
	0x7c: insInfo{"ILL", amNone},
	0x8c: insInfo{"STY", amAbs},
	0x9c: insInfo{"ILL", amNone},
	0xac: insInfo{"LDY", amAbs},
	0xbc: insInfo{"LDY", amAbsX},
	0xcc: insInfo{"CPY", amAbs},
	0xdc: insInfo{"ILL", amNone},
	0xec: insInfo{"CPX", amAbs},
	0xfc: insInfo{"ILL", amNone},

	0x0d: insInfo{"ORA", amAbs},
	0x1d: insInfo{"ORA", amAbsX},
	0x2d: insInfo{"AND", amAbs},
	0x3d: insInfo{"AND", amAbsX},
	0x4d: insInfo{"EOR", amAbs},
	0x5d: insInfo{"EOR", amAbsX},
	0x6d: insInfo{"ADC", amAbs},
	0x7d: insInfo{"ADC", amAbsX},
	0x8d: insInfo{"STA", amAbs},
	0x9d: insInfo{"STA", amAbsX},
	0xad: insInfo{"LDA", amAbs},
	0xbd: insInfo{"LDA", amAbsX},
	0xcd: insInfo{"CMP", amAbs},
	0xdd: insInfo{"CMP", amAbsX},
	0xed: insInfo{"SBC", amAbs},
	0xfd: insInfo{"SBC", amAbsX},

	0x0e: insInfo{"ASL", amAbs},
	0x1e: insInfo{"ASL", amAbsX},
	0x2e: insInfo{"ROL", amAbs},
	0x3e: insInfo{"ROL", amAbsX},
	0x4e: insInfo{"LSR", amAbs},
	0x5e: insInfo{"LSR", amAbsX},
	0x6e: insInfo{"ROR", amAbs},
	0x7e: insInfo{"ROR", amAbsX},
	0x8e: insInfo{"STX", amAbs},
	0x9e: insInfo{"ILL", amNone},
	0xae: insInfo{"LDX", amAbs},
	0xbe: insInfo{"LDX", amAbsY},
	0xce: insInfo{"DEC", amAbs},
	0xde: insInfo{"DEC", amAbsX},
	0xee: insInfo{"INC", amAbs},
	0xfe: insInfo{"INC", amAbsX},
}

// opcodeLookup returns the instruction information for this opcode.
func opcodeLookup(code uint8) *insInfo {
	if info, ok := opcodeTable[code]; ok {
		return &info
	}
	return &insInfo{"ILL", amNone}
}

// insDescr maps the instruction mneumonic onto a full description.
var insDescr = map[string]string{
	"ADC": "add with carry",
	"AND": "and (with accumulator)",
	"ASL": "arithmetic shift left",
	"BCC": "branch on carry clear",
	"BCS": "branch on carry set",
	"BEQ": "branch on equal (zero set)",
	"BIT": "bit test",
	"BMI": "branch on minus (negative set)",
	"BNE": "branch on not equal (zero clear)",
	"BPL": "branch on plus (negative clear)",
	"BRK": "break / interrupt",
	"BVC": "branch on overflow clear",
	"BVS": "branch on overflow set",
	"CLC": "clear carry",
	"CLD": "clear decimal",
	"CLI": "clear interrupt disable",
	"CLV": "clear overflow",
	"CMP": "compare (with accumulator)",
	"CPX": "compare with X",
	"CPY": "compare with Y",
	"DEC": "decrement",
	"DEX": "decrement X",
	"DEY": "decrement Y",
	"EOR": "exclusive or (with accumulator)",
	"INC": "increment",
	"INX": "increment X",
	"INY": "increment Y",
	"JMP": "jump",
	"JSR": "jump subroutine",
	"LDA": "load accumulator",
	"LDX": "load X",
	"LDY": "load Y",
	"LSR": "logical shift right",
	"NOP": "no operation",
	"ORA": "or with accumulator",
	"PHA": "push accumulator",
	"PHP": "push processor status (SR)",
	"PLA": "pull accumulator",
	"PLP": "pull processor status (SR)",
	"ROL": "rotate left",
	"ROR": "rotate right",
	"RTI": "return from interrupt",
	"RTS": "return from subroutine",
	"SBC": "subtract with carry",
	"SEC": "set carry",
	"SED": "set decimal",
	"SEI": "set interrupt disable",
	"STA": "store accumulator",
	"STX": "store X",
	"STY": "store Y",
	"TAX": "transfer accumulator to X",
	"TAY": "transfer accumulator to Y",
	"TSX": "transfer stack pointer to X",
	"TXA": "transfer X to accumulator",
	"TXS": "transfer X to stack pointer",
	"TYA": "transfer Y to accumulator",
}

//-----------------------------------------------------------------------------

func opcodeFuncName(code uint8) string {
	x := opcodeLookup(code)
	return fmt.Sprintf("op%s%s", x.ins, modeDescr[x.mode].suffix)
}

func genOpcodeFunc() string {

	// get the set of unique opcode function names
	fs := make(map[string]bool)
	for code := 0; code < 256; code++ {
		fs[opcodeFuncName(uint8(code))] = true
	}
	f := make([]string, len(fs))
	i := 0
	for name := range fs {
		f[i] = name
		i++
	}
	sort.Strings(f)
	return strings.Join(f, "\n")
}

func genOpcodeTable() string {

	f := make([]string, 16)
	for i := 0; i < 16; i++ {
		l := make([]string, 16)
		for j := 0; j < 16; j++ {
			code := uint8((i * 16) + j)
			l[j] = opcodeFuncName(code)
		}
		f[i] = strings.Join(l, ",")
	}

	s := make([]string, 3)
	s[0] = "var opcodeTable = [256]opFunc{"
	s[1] = strings.Join(f, "\n")
	s[2] = "}"
	return strings.Join(s, "\n")
}

func GenCode() string {

	s := make([]string, 2)

	s[0] = genOpcodeFunc()
	s[1] = genOpcodeTable()

	return strings.Join(s, "\n")

}

//-----------------------------------------------------------------------------

func Disassemble(mem []byte) (string, uint) {

	var ofs uint
	var s []string

	opcode := mem[ofs]
	ofs++
	info := opcodeLookup(opcode)

	// instruction mneumonic
	s = append(s, info.ins)

	switch info.mode {
	case amNone:
		// illegal - no operands
	case amAcc:
		s = append(s, "accumulator")
	case amAbs:
		s = append(s, "absolute")
	case amAbsX:
		s = append(s, "absolute, X-indexed")
	case amAbsY:
		s = append(s, "absolute, Y-indexed")
	case amImm:
		// immediate - 1 byte operand
		operand := mem[ofs]
		ofs++
		s = append(s, fmt.Sprintf("$%02x", operand))
	case amImpl:
		// implied - no operands
	case amInd:
		s = append(s, "indirect")
	case amXInd:
		s = append(s, "X-indexed, indirect")
	case amIndY:
		// indirect, Y-indexed - 1 byte operand
		operand := mem[ofs]
		ofs++
		s = append(s, fmt.Sprintf("($%02x),Y", operand))
	case amRel:
		// relative - 1 byte operand
		operand := mem[ofs]
		ofs++
		s = append(s, fmt.Sprintf("$%02x", operand))
	case amZpg:
		s = append(s, "zeropage")
	case amZpgX:
		s = append(s, "zeropage, X-indexed")
	case amZpgY:
		s = append(s, "zeropage, Y-indexed")
	default:
		panic("bad address mode")
	}

	return strings.Join(s, " "), ofs
}

//-----------------------------------------------------------------------------
