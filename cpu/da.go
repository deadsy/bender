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

	0x00: insInfo{"brk", amImpl},
	0x10: insInfo{"bpl", amRel},
	0x20: insInfo{"jsr", amAbs},
	0x30: insInfo{"bmi", amRel},
	0x40: insInfo{"rti", amImpl},
	0x50: insInfo{"bvc", amRel},
	0x60: insInfo{"rts", amImpl},
	0x70: insInfo{"bvs", amRel},
	0x80: insInfo{"ill", amNone},
	0x90: insInfo{"bcc", amRel},
	0xa0: insInfo{"ldy", amImm},
	0xb0: insInfo{"bcs", amRel},
	0xc0: insInfo{"cpy", amImm},
	0xd0: insInfo{"bne", amRel},
	0xe0: insInfo{"cpx", amImm},
	0xf0: insInfo{"beq", amImm},

	0x01: insInfo{"ora", amXInd},
	0x11: insInfo{"ora", amIndY},
	0x21: insInfo{"and", amXInd},
	0x31: insInfo{"and", amIndY},
	0x41: insInfo{"eor", amXInd},
	0x51: insInfo{"eor", amIndY},
	0x61: insInfo{"adc", amXInd},
	0x71: insInfo{"adc", amIndY},
	0x81: insInfo{"sta", amXInd},
	0x91: insInfo{"sta", amIndY},
	0xa1: insInfo{"lda", amXInd},
	0xb1: insInfo{"lda", amIndY},
	0xc1: insInfo{"cmp", amXInd},
	0xd1: insInfo{"cmp", amIndY},
	0xe1: insInfo{"sbc", amXInd},
	0xf1: insInfo{"sbc", amIndY},

	0xa2: insInfo{"ldx", amImm},

	0x04: insInfo{"ill", amNone},
	0x14: insInfo{"ill", amNone},
	0x24: insInfo{"bit", amZpg},
	0x34: insInfo{"ill", amNone},
	0x44: insInfo{"ill", amNone},
	0x54: insInfo{"ill", amNone},
	0x64: insInfo{"ill", amNone},
	0x74: insInfo{"ill", amNone},
	0x84: insInfo{"sty", amZpg},
	0x94: insInfo{"sty", amZpgX},
	0xa4: insInfo{"ldy", amZpg},
	0xb4: insInfo{"ldy", amZpgX},
	0xc4: insInfo{"cpy", amZpg},
	0xd4: insInfo{"ill", amNone},
	0xe4: insInfo{"cpx", amZpg},
	0xf4: insInfo{"ill", amNone},

	0x05: insInfo{"ora", amZpg},
	0x15: insInfo{"ora", amZpgX},
	0x25: insInfo{"and", amZpg},
	0x35: insInfo{"and", amZpgX},
	0x45: insInfo{"eor", amZpg},
	0x55: insInfo{"eor", amZpgX},
	0x65: insInfo{"adc", amZpg},
	0x75: insInfo{"adc", amZpgX},
	0x85: insInfo{"sta", amZpg},
	0x95: insInfo{"sta", amZpgX},
	0xa5: insInfo{"lda", amZpg},
	0xb5: insInfo{"lda", amZpgX},
	0xc5: insInfo{"cmp", amZpg},
	0xd5: insInfo{"cmp", amZpgX},
	0xe5: insInfo{"sbc", amZpg},
	0xf5: insInfo{"sbc", amZpgX},

	0x06: insInfo{"asl", amZpg},
	0x16: insInfo{"asl", amZpgX},
	0x26: insInfo{"rol", amZpg},
	0x36: insInfo{"rol", amZpgX},
	0x46: insInfo{"lsr", amZpg},
	0x56: insInfo{"lsr", amZpgX},
	0x66: insInfo{"ror", amZpg},
	0x76: insInfo{"ror", amZpgX},
	0x86: insInfo{"stx", amZpg},
	0x96: insInfo{"stx", amZpgY},
	0xa6: insInfo{"ldx", amZpg},
	0xb6: insInfo{"ldx", amZpgY},
	0xc6: insInfo{"dec", amZpg},
	0xd6: insInfo{"dec", amZpgX},
	0xe6: insInfo{"inc", amZpg},
	0xf6: insInfo{"inc", amZpgX},

	0x08: insInfo{"php", amImpl},
	0x18: insInfo{"clc", amImpl},
	0x28: insInfo{"plp", amImpl},
	0x38: insInfo{"sec", amImpl},
	0x48: insInfo{"pha", amImpl},
	0x58: insInfo{"cli", amImpl},
	0x68: insInfo{"pla", amImpl},
	0x78: insInfo{"sei", amImpl},
	0x88: insInfo{"dey", amImpl},
	0x98: insInfo{"tya", amImpl},
	0xa8: insInfo{"tay", amImpl},
	0xb8: insInfo{"clv", amImpl},
	0xc8: insInfo{"iny", amImpl},
	0xd8: insInfo{"cld", amImpl},
	0xe8: insInfo{"inx", amImpl},
	0xf8: insInfo{"sed", amImpl},

	0x09: insInfo{"ora", amImm},
	0x19: insInfo{"ora", amAbsY},
	0x29: insInfo{"and", amImm},
	0x39: insInfo{"and", amAbsY},
	0x49: insInfo{"eor", amImm},
	0x59: insInfo{"eor", amAbsY},
	0x69: insInfo{"adc", amImm},
	0x79: insInfo{"adc", amAbsY},
	0x89: insInfo{"ill", amNone},
	0x99: insInfo{"sta", amAbsY},
	0xa9: insInfo{"lda", amImm},
	0xb9: insInfo{"lda", amAbsY},
	0xc9: insInfo{"cmp", amImm},
	0xd9: insInfo{"cmp", amAbsY},
	0xe9: insInfo{"sbc", amImm},
	0xf9: insInfo{"sbc", amAbsY},

	0x0a: insInfo{"asl", amAcc},
	0x1a: insInfo{"ill", amNone},
	0x2a: insInfo{"rol", amAcc},
	0x3a: insInfo{"ill", amNone},
	0x4a: insInfo{"lsr", amAcc},
	0x5a: insInfo{"ill", amNone},
	0x6a: insInfo{"ror", amAcc},
	0x7a: insInfo{"ill", amNone},
	0x8a: insInfo{"txa", amImpl},
	0x9a: insInfo{"txs", amImpl},
	0xaa: insInfo{"tax", amImpl},
	0xba: insInfo{"tsx", amImpl},
	0xca: insInfo{"dex", amImpl},
	0xda: insInfo{"ill", amNone},
	0xea: insInfo{"nop", amImpl},
	0xfa: insInfo{"ill", amNone},

	0x0c: insInfo{"ill", amNone},
	0x1c: insInfo{"ill", amNone},
	0x2c: insInfo{"bit", amAbs},
	0x3c: insInfo{"ill", amNone},
	0x4c: insInfo{"jmp", amAbs},
	0x5c: insInfo{"ill", amNone},
	0x6c: insInfo{"jmp", amInd},
	0x7c: insInfo{"ill", amNone},
	0x8c: insInfo{"sty", amAbs},
	0x9c: insInfo{"ill", amNone},
	0xac: insInfo{"ldy", amAbs},
	0xbc: insInfo{"ldy", amAbsX},
	0xcc: insInfo{"cpy", amAbs},
	0xdc: insInfo{"ill", amNone},
	0xec: insInfo{"cpx", amAbs},
	0xfc: insInfo{"ill", amNone},

	0x0d: insInfo{"ora", amAbs},
	0x1d: insInfo{"ora", amAbsX},
	0x2d: insInfo{"and", amAbs},
	0x3d: insInfo{"and", amAbsX},
	0x4d: insInfo{"eor", amAbs},
	0x5d: insInfo{"eor", amAbsX},
	0x6d: insInfo{"adc", amAbs},
	0x7d: insInfo{"adc", amAbsX},
	0x8d: insInfo{"sta", amAbs},
	0x9d: insInfo{"sta", amAbsX},
	0xad: insInfo{"lda", amAbs},
	0xbd: insInfo{"lda", amAbsX},
	0xcd: insInfo{"cmp", amAbs},
	0xdd: insInfo{"cmp", amAbsX},
	0xed: insInfo{"sbc", amAbs},
	0xfd: insInfo{"sbc", amAbsX},

	0x0e: insInfo{"asl", amAbs},
	0x1e: insInfo{"asl", amAbsX},
	0x2e: insInfo{"rol", amAbs},
	0x3e: insInfo{"rol", amAbsX},
	0x4e: insInfo{"lsr", amAbs},
	0x5e: insInfo{"lsr", amAbsX},
	0x6e: insInfo{"ror", amAbs},
	0x7e: insInfo{"ror", amAbsX},
	0x8e: insInfo{"stx", amAbs},
	0x9e: insInfo{"ill", amNone},
	0xae: insInfo{"ldx", amAbs},
	0xbe: insInfo{"ldx", amAbsY},
	0xce: insInfo{"dec", amAbs},
	0xde: insInfo{"dec", amAbsX},
	0xee: insInfo{"inc", amAbs},
	0xfe: insInfo{"inc", amAbsX},
}

// opcodeLookup returns the instruction information for this opcode.
func opcodeLookup(code uint8) *insInfo {
	if info, ok := opcodeTable[code]; ok {
		return &info
	}
	return &insInfo{"ill", amNone}
}

// insDescr maps the instruction mneumonic onto a full description.
var insDescr = map[string]string{
	"adc": "add with carry",
	"and": "and (with accumulator)",
	"asl": "arithmetic shift left",
	"bcc": "branch on carry clear",
	"bcs": "branch on carry set",
	"beq": "branch on equal (zero set)",
	"bit": "bit test",
	"bmi": "branch on minus (negative set)",
	"bne": "branch on not equal (zero clear)",
	"bpl": "branch on plus (negative clear)",
	"brk": "break / interrupt",
	"bvc": "branch on overflow clear",
	"bvs": "branch on overflow set",
	"clc": "clear carry",
	"cld": "clear decimal",
	"cli": "clear interrupt disable",
	"clv": "clear overflow",
	"cmp": "compare (with accumulator)",
	"cpx": "compare with X",
	"cpy": "compare with Y",
	"dec": "decrement",
	"dex": "decrement X",
	"dey": "decrement Y",
	"eor": "exclusive or (with accumulator)",
	"inc": "increment",
	"inx": "increment X",
	"iny": "increment Y",
	"jmp": "jump",
	"jsr": "jump subroutine",
	"lda": "load accumulator",
	"ldx": "load X",
	"ldy": "load Y",
	"lsr": "logical shift right",
	"nop": "no operation",
	"ora": "or with accumulator",
	"pha": "push accumulator",
	"php": "push processor status (SR)",
	"pla": "pull accumulator",
	"plp": "pull processor status (SR)",
	"rol": "rotate left",
	"ror": "rotate right",
	"rti": "return from interrupt",
	"rts": "return from subroutine",
	"sbc": "subtract with carry",
	"sec": "set carry",
	"sed": "set decimal",
	"sei": "set interrupt disable",
	"sta": "store accumulator",
	"stx": "store X",
	"sty": "store Y",
	"tax": "transfer accumulator to X",
	"tay": "transfer accumulator to Y",
	"tsx": "transfer stack pointer to X",
	"txa": "transfer X to accumulator",
	"txs": "transfer X to stack pointer",
	"tya": "transfer Y to accumulator",
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

func Disassemble(adr uint16, mem []byte) (string, int) {

	var ofs int
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
		// accumulator - no operands
		s = append(s, "a")
	case amAbs:
		// absolute - 2 byte operand
		operand := int(mem[ofs]) + (int(mem[ofs+1]) << 8)
		ofs += 2
		s = append(s, fmt.Sprintf("$%04x", operand))
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
		s = append(s, fmt.Sprintf("($%02x),y", operand))
	case amRel:
		// relative - 1 byte operand
		operand := mem[ofs]
		ofs++
		dst := uint16(int(adr) + int(int8(operand)) + 2)
		s = append(s, fmt.Sprintf("$%02x\t\t; $%04x", operand, dst))
	case amZpg:
		s = append(s, "zeropage")
	case amZpgX:
		// zeropage, X-indexed - 1 byte operand
		operand := mem[ofs]
		ofs++
		s = append(s, fmt.Sprintf("$%02x,x", operand))
	case amZpgY:
		s = append(s, "zeropage, Y-indexed")
	default:
		panic("bad address mode")
	}

	return strings.Join(s, " "), ofs
}

//-----------------------------------------------------------------------------
