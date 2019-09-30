//-----------------------------------------------------------------------------
/*

6502 CPU Emulator

See also:

https://github.com/redcode/6502
https://www.masswerk.at/6502/6502_instruction_set.html

*/
//-----------------------------------------------------------------------------

package cpu

type addressMode int

const (
	amNone addressMode = iota
	amA                // accumulator
	amAbs              // absolute
	amAbsX             // absolute, X-indexed
	amAbsY             // absolute, Y-indexed
	amImm              // immediate
	amImpl             // implied
	amInd              // indirect
	amXInd             // X-indexed, indirect
	amIndY             // indirect, Y-indexed
	amRel              // relative
	amZpg              // zeropage
	amZpgX             // zeropage, X-indexed
	amZpgY             // zeropage, Y-indexed
)

type opInfo struct {
	mneumonic string
	mode      addressMode
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


}

/*


HI	LO-NIBBLE
06	07	08	09	0A	0B	0C	0D	0E	0F
ASL zpg	  ---	PHP impl	ORA #	ASL A	---	---	ORA abs	ASL abs	---
ASL zpg,X	---	CLC impl	ORA abs,Y	---	---	---	ORA abs,X	ASL abs,X	---
ROL zpg	  ---	PLP impl	AND #	ROL A	---	BIT abs	AND abs	ROL abs	---
ROL zpg,X	---	SEC impl	AND abs,Y	---	---	---	AND abs,X	ROL abs,X	---
LSR zpg	  ---	PHA impl	EOR #	LSR A	---	JMP abs	EOR abs	LSR abs	---
LSR zpg,X	---	CLI impl	EOR abs,Y	---	---	---	EOR abs,X	LSR abs,X	---
ROR zpg	  ---	PLA impl	ADC #	ROR A	---	JMP ind	ADC abs	ROR abs	---
ROR zpg,X	---	SEI impl	ADC abs,Y	---	---	---	ADC abs,X	ROR abs,X	---
STX zpg	  ---	DEY impl	---	TXA impl	---	STY abs	STA abs	STX abs	---
STX zpg,Y	---	TYA impl	STA abs,Y	TXS impl	---	---	STA abs,X	---	---
LDX zpg	  ---	TAY impl	LDA #	TAX impl	---	LDY abs	LDA abs	LDX abs	---
LDX zpg,Y	---	CLV impl	LDA abs,Y	TSX impl	---	LDY abs,X	LDA abs,X	LDX abs,Y	---
DEC zpg	  ---	INY impl	CMP #	DEX impl	---	CPY abs	CMP abs	DEC abs	---
DEC zpg,X	---	CLD impl	CMP abs,Y	---	---	---	CMP abs,X	DEC abs,X	---
INC zpg	  ---	INX impl	SBC #	NOP impl	---	CPX abs	SBC abs	INC abs	---
INC zpg,X	---	SED impl	SBC abs,Y	---	---	---	SBC abs,X	INC abs,X	---

*/

//-----------------------------------------------------------------------------
