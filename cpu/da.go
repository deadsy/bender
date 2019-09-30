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

type opInfo struct {
	ins  string  // instruction mneumonic
	mode adrMode // address mode
}

var opcodes = map[uint8]opInfo{

	0x00: opInfo{"BRK", amImpl},
	0x10: opInfo{"BPL", amRel},
	0x20: opInfo{"JSR", amAbs},
	0x30: opInfo{"BMI", amRel},
	0x40: opInfo{"RTI", amImpl},
	0x50: opInfo{"BVC", amRel},
	0x60: opInfo{"RTS", amImpl},
	0x70: opInfo{"BVS", amRel},
	0x80: opInfo{"ILL", amNone},
	0x90: opInfo{"BCC", amRel},
	0xa0: opInfo{"LDY", amImm},
	0xb0: opInfo{"BCS", amRel},
	0xc0: opInfo{"CPY", amImm},
	0xd0: opInfo{"BNE", amRel},
	0xe0: opInfo{"CPX", amImm},
	0xf0: opInfo{"BEQ", amImm},

	0x01: opInfo{"ORA", amXInd},
	0x11: opInfo{"ORA", amIndY},
	0x21: opInfo{"AND", amXInd},
	0x31: opInfo{"AND", amIndY},
	0x41: opInfo{"EOR", amXInd},
	0x51: opInfo{"EOR", amIndY},
	0x61: opInfo{"ADC", amXInd},
	0x71: opInfo{"ADC", amIndY},
	0x81: opInfo{"STA", amXInd},
	0x91: opInfo{"STA", amIndY},
	0xa1: opInfo{"LDA", amXInd},
	0xb1: opInfo{"LDA", amIndY},
	0xc1: opInfo{"CMP", amXInd},
	0xd1: opInfo{"CMP", amIndY},
	0xe1: opInfo{"SBC", amXInd},
	0xf1: opInfo{"SBC", amIndY},

	0xa2: opInfo{"LDX", amImm},

	0x04: opInfo{"ILL", amNone},
	0x14: opInfo{"ILL", amNone},
	0x24: opInfo{"BIT", amZpg},
	0x34: opInfo{"ILL", amNone},
	0x44: opInfo{"ILL", amNone},
	0x54: opInfo{"ILL", amNone},
	0x64: opInfo{"ILL", amNone},
	0x74: opInfo{"ILL", amNone},
	0x84: opInfo{"STY", amZpg},
	0x94: opInfo{"STY", amZpgX},
	0xa4: opInfo{"LDY", amZpg},
	0xb4: opInfo{"LDY", amZpgX},
	0xc4: opInfo{"CPY", amZpg},
	0xd4: opInfo{"ILL", amNone},
	0xe4: opInfo{"CPX", amZpg},
	0xf4: opInfo{"ILL", amNone},

	0x05: opInfo{"ORA", amZpg},
	0x15: opInfo{"ORA", amZpgX},
	0x25: opInfo{"AND", amZpg},
	0x35: opInfo{"AND", amZpgX},
	0x45: opInfo{"EOR", amZpg},
	0x55: opInfo{"EOR", amZpgX},
	0x65: opInfo{"ADC", amZpg},
	0x75: opInfo{"ADC", amZpgX},
	0x85: opInfo{"STA", amZpg},
	0x95: opInfo{"STA", amZpgX},
	0xa5: opInfo{"LDA", amZpg},
	0xb5: opInfo{"LDA", amZpgX},
	0xc5: opInfo{"CMP", amZpg},
	0xd5: opInfo{"CMP", amZpgX},
	0xe5: opInfo{"SBC", amZpg},
	0xf5: opInfo{"SBC", amZpgX},

	0x06: opInfo{"ASL", amZpg},
	0x16: opInfo{"ASL", amZpgX},
	0x26: opInfo{"ROL", amZpg},
	0x36: opInfo{"ROL", amZpgX},
	0x46: opInfo{"LSR", amZpg},
	0x56: opInfo{"LSR", amZpgX},
	0x66: opInfo{"ROR", amZpg},
	0x76: opInfo{"ROR", amZpgX},
	0x86: opInfo{"STX", amZpg},
	0x96: opInfo{"STX", amZpgY},
	0xa6: opInfo{"LDX", amZpg},
	0xb6: opInfo{"LDX", amZpgY},
	0xc6: opInfo{"DEC", amZpg},
	0xd6: opInfo{"DEC", amZpgX},
	0xe6: opInfo{"INC", amZpg},
	0xf6: opInfo{"INC", amZpgX},

	0x08: opInfo{"PHP", amImpl},
	0x18: opInfo{"CLC", amImpl},
	0x28: opInfo{"PLP", amImpl},
	0x38: opInfo{"SEC", amImpl},
	0x48: opInfo{"PHA", amImpl},
	0x58: opInfo{"CLI", amImpl},
	0x68: opInfo{"PLA", amImpl},
	0x78: opInfo{"SEI", amImpl},
	0x88: opInfo{"DEY", amImpl},
	0x98: opInfo{"TYA", amImpl},
	0xa8: opInfo{"TAY", amImpl},
	0xb8: opInfo{"CLV", amImpl},
	0xc8: opInfo{"INY", amImpl},
	0xd8: opInfo{"CLD", amImpl},
	0xe8: opInfo{"INX", amImpl},
	0xf8: opInfo{"SED", amImpl},

	0x09: opInfo{"ORA", amImm},
	0x19: opInfo{"ORA", amAbsY},
	0x29: opInfo{"AND", amImm},
	0x39: opInfo{"AND", amAbsY},
	0x49: opInfo{"EOR", amImm},
	0x59: opInfo{"EOR", amAbsY},
	0x69: opInfo{"ADC", amImm},
	0x79: opInfo{"ADC", amAbsY},
	0x89: opInfo{"ILL", amNone},
	0x99: opInfo{"STA", amAbsY},
	0xa9: opInfo{"LDA", amImm},
	0xb9: opInfo{"LDA", amAbsY},
	0xc9: opInfo{"CMP", amImm},
	0xd9: opInfo{"CMP", amAbsY},
	0xe9: opInfo{"SBC", amImm},
	0xf9: opInfo{"SBC", amAbsY},

	0x0a: opInfo{"ASL", amAcc},
	0x1a: opInfo{"ILL", amNone},
	0x2a: opInfo{"ROL", amAcc},
	0x3a: opInfo{"ILL", amNone},
	0x4a: opInfo{"LSR", amAcc},
	0x5a: opInfo{"ILL", amNone},
	0x6a: opInfo{"ROR", amAcc},
	0x7a: opInfo{"ILL", amNone},
	0x8a: opInfo{"TXA", amImpl},
	0x9a: opInfo{"TXS", amImpl},
	0xaa: opInfo{"TAX", amImpl},
	0xba: opInfo{"TSX", amImpl},
	0xca: opInfo{"DEX", amImpl},
	0xda: opInfo{"ILL", amNone},
	0xea: opInfo{"NOP", amImpl},
	0xfa: opInfo{"ILL", amNone},

	0x0c: opInfo{"ILL", amNone},
	0x1c: opInfo{"ILL", amNone},
	0x2c: opInfo{"BIT", amAbs},
	0x3c: opInfo{"ILL", amNone},
	0x4c: opInfo{"JMP", amAbs},
	0x5c: opInfo{"ILL", amNone},
	0x6c: opInfo{"JMP", amInd},
	0x7c: opInfo{"ILL", amNone},
	0x8c: opInfo{"STY", amAbs},
	0x9c: opInfo{"ILL", amNone},
	0xac: opInfo{"LDY", amAbs},
	0xbc: opInfo{"LDY", amAbsX},
	0xcc: opInfo{"CPY", amAbs},
	0xdc: opInfo{"ILL", amNone},
	0xec: opInfo{"CPX", amAbs},
	0xfc: opInfo{"ILL", amNone},

	0x0d: opInfo{"ORA", amAbs},
	0x1d: opInfo{"ORA", amAbsX},
	0x2d: opInfo{"AND", amAbs},
	0x3d: opInfo{"AND", amAbsX},
	0x4d: opInfo{"EOR", amAbs},
	0x5d: opInfo{"EOR", amAbsX},
	0x6d: opInfo{"ADC", amAbs},
	0x7d: opInfo{"ADC", amAbsX},
	0x8d: opInfo{"STA", amAbs},
	0x9d: opInfo{"STA", amAbsX},
	0xad: opInfo{"LDA", amAbs},
	0xbd: opInfo{"LDA", amAbsX},
	0xcd: opInfo{"CMP", amAbs},
	0xdd: opInfo{"CMP", amAbsX},
	0xed: opInfo{"SBC", amAbs},
	0xfd: opInfo{"SBC", amAbsX},

	0x0e: opInfo{"ASL", amAbs},
	0x1e: opInfo{"ASL", amAbsX},
	0x2e: opInfo{"ROL", amAbs},
	0x3e: opInfo{"ROL", amAbsX},
	0x4e: opInfo{"LSR", amAbs},
	0x5e: opInfo{"LSR", amAbsX},
	0x6e: opInfo{"ROR", amAbs},
	0x7e: opInfo{"ROR", amAbsX},
	0x8e: opInfo{"STX", amAbs},
	0x9e: opInfo{"ILL", amNone},
	0xae: opInfo{"LDX", amAbs},
	0xbe: opInfo{"LDX", amAbsY},
	0xce: opInfo{"DEC", amAbs},
	0xde: opInfo{"DEC", amAbsX},
	0xee: opInfo{"INC", amAbs},
	0xfe: opInfo{"INC", amAbsX},
}

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
	fname := "opILL"
	x, ok := opcodes[code]
	if ok {
		fname = fmt.Sprintf("op%s%s", x.ins, modeDescr[x.mode].suffix)
	}
	return fname
}

func genOpcodeFuncs() string {
	s := make([]string, 256)
	for i := 0; i < 256; i++ {
		s = append(s, opcodeFuncName(uint8(i)))
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------
