//-----------------------------------------------------------------------------
/*

6502 CPU Emulator

*/
//-----------------------------------------------------------------------------

package cpu

import "fmt"

//-----------------------------------------------------------------------------

func (m *M6502) setC(val uint) {
	if val>>8 != 0 {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
}

func (m *M6502) setN(val uint8) {
	if val&0x80 != 0 {
		m.P |= flagN
	} else {
		m.P &= ^flagN
	}
}

func (m *M6502) setZ(val uint8) {
	if val == 0 {
		m.P |= flagZ
	} else {
		m.P &= ^flagZ
	}
}

func (m *M6502) setNZ(val uint8) {
	m.setN(val)
	m.setZ(val)
}

//-----------------------------------------------------------------------------

func (m *M6502) read16(adr uint16) uint16 {
	l := uint16(m.Mem.Read8(adr))
	h := uint16(m.Mem.Read8(adr + 1))
	return (h << 8) | l
}

func (m *M6502) push8(val uint8) {
	m.Mem.Write8(stkAddress+uint16(m.S), val)
	m.S--
}

func (m *M6502) pop8() uint8 {
	m.S++
	return m.Mem.Read8(stkAddress + uint16(m.S))
}

func (m *M6502) push16(val uint16) {
	m.push8(uint8(val >> 8))
	m.push8(uint8(val))
}

func (m *M6502) pop16() uint16 {
	l := m.pop8()
	h := m.pop8()
	return (uint16(h) << 8) | uint16(l)
}

//-----------------------------------------------------------------------------
// modal write functions

func (m *M6502) writeZeroPage(val uint8) {
	ea := m.Mem.Read8(m.PC + 1)
	m.Mem.Write8(uint16(ea), val)
}

func (m *M6502) writeZeroPageX(val uint8) {
	ea := m.Mem.Read8(m.PC+1) + m.X
	m.Mem.Write8(uint16(ea), val)
}

func (m *M6502) writeZeroPageY(val uint8) {
	ea := m.Mem.Read8(m.PC+1) + m.Y
	m.Mem.Write8(uint16(ea), val)
}

func (m *M6502) writeAbsolute(val uint8) {
	ea := m.read16(m.PC + 1)
	m.Mem.Write8(ea, val)
}

func (m *M6502) writeAbsoluteX(val uint8) {
	ea := m.read16(m.PC+1) + uint16(m.X)
	m.Mem.Write8(ea, val)
}

func (m *M6502) writeAbsoluteY(val uint8) {
	ea := m.read16(m.PC+1) + uint16(m.Y)
	m.Mem.Write8(ea, val)
}

func (m *M6502) writeIndirectX(val uint8) {
	ea := m.read16(uint16(m.Mem.Read8(m.PC+1) + m.X))
	m.Mem.Write8(ea, val)
}

func (m *M6502) writeIndirectY(val uint8) {
	ea := m.read16(uint16(m.Mem.Read8(m.PC+1))) + uint16(m.Y)
	m.Mem.Write8(ea, val)
}

//-----------------------------------------------------------------------------
// modal read functions

func (m *M6502) readImmediate() uint8 {
	return m.Mem.Read8(m.PC + 1)
}

func (m *M6502) readZeroPage() (uint8, uint16) {
	ea := uint16(m.Mem.Read8(m.PC + 1))
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readZeroPageX() (uint8, uint16) {
	ea := uint16(m.Mem.Read8(m.PC+1) + m.X)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readZeroPageY() (uint8, uint16) {
	ea := uint16(m.Mem.Read8(m.PC+1) + m.Y)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readAbsolute() (uint8, uint16) {
	ea := m.read16(m.PC + 1)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readAbsoluteX() (uint8, uint16) {
	ea := m.read16(m.PC+1) + uint16(m.X)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readAbsoluteY() (uint8, uint16) {
	ea := m.read16(m.PC+1) + uint16(m.Y)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readIndirectX() (uint8, uint16) {
	ea := m.read16(uint16(m.Mem.Read8(m.PC+1) + m.X))
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readIndirectY() (uint8, uint16) {
	ea := m.read16(uint16(m.Mem.Read8(m.PC+1))) + uint16(m.Y)
	return m.Mem.Read8(ea), ea
}

func (m *M6502) readAbsoluteXPenalized() (uint8, uint, uint16) {
	ea := m.read16(m.PC + 1)
	var n uint
	if (ea&0xff)+uint16(m.X) > 0xff {
		n = 1
	}
	ea += uint16(m.X)
	return m.Mem.Read8(ea), n, ea
}

func (m *M6502) readAbsoluteYPenalized() (uint8, uint, uint16) {
	ea := m.read16(m.PC + 1)
	var n uint
	if (ea&0xff)+uint16(m.Y) > 0xff {
		n = 1
	}
	ea += uint16(m.Y)
	return m.Mem.Read8(ea), n, ea
}

func (m *M6502) readIndirectYPenalized() (uint8, uint, uint16) {
	ea := m.read16(uint16(m.Mem.Read8(m.PC + 1)))
	var n uint
	if (ea&0xff)+uint16(m.Y) > 0xff {
		n = 1
	}
	ea += uint16(m.Y)
	return m.Mem.Read8(ea), n, ea
}

//-----------------------------------------------------------------------------

// opBranch does a relative branch if a condition is true.
func (m *M6502) opBranch(cond bool) uint {
	cycles := 2
	if cond {
		pc := uint16(m.PC + 2)
		ofs := int8(m.Mem.Read8(m.PC + 1))
		tgt := uint16(int(pc) + int(ofs))
		if (tgt >> 8) == (pc >> 8) {
			// same page: +1 cycle
			cycles++
		} else {
			// different page: +2 cycles
			cycles += 2
		}
		m.PC = tgt
	} else {
		m.PC += 2
	}
	return uint(cycles)
}

// opCompare sets the NZC values for the register/value compare operation.
func (m *M6502) opCompare(reg, val uint8) {
	m.setNZ(reg - val)
	if reg >= val {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
}

// opADC add with carry
func (m *M6502) opADC(v uint8) {

	a := uint(m.A)
	old := a
	rhs := uint(v)
	c := uint(m.P & flagC)

	if m.P&flagD != 0 {
		lo := (old & 0x0F) + (rhs & 0x0F) + c
		if lo >= 0x0A {
			lo = ((lo + 0x06) & 0x0F) + 0x10
		}
		a = (old & 0xF0) + (rhs & 0xF0) + lo
		// overflow
		res := int(int8(old&0xF0)) + int(int8(rhs&0xF0)) + int(int8(lo))
		if (res < -128) || (res > 127) {
			m.P |= flagV
		} else {
			m.P &= ^flagV
		}
		// zero
		if (old+rhs+c)&0xff == 0 {
			m.P |= flagZ
		} else {
			m.P &= ^flagZ
		}
		// negative
		m.setN(uint8(a))
		if a >= 0xA0 {
			a += 0x60
		}
		// carry
		m.setC(a)
		m.A = uint8(a)
	} else {
		a += rhs + c
		m.A = uint8(a)
		// carry
		m.setC(a)
		// overflow
		if (((old ^ rhs) & 0x80) == 0) && (((old ^ a) & 0x80) != 0) {
			m.P |= flagV
		} else {
			m.P &= ^flagV
		}
		// negative, zero
		m.setNZ(m.A)
	}
}

// opSBC subtract with cary
func (m *M6502) opSBC(v uint8) {

	a := uint(m.A)
	old := a
	rhs := uint(v)
	c := uint(m.P & flagC)

	if m.P&flagD != 0 {

		lo := (old & 0x0F) - (rhs & 0x0F) + c - 1
		if lo&0x80 != 0 {
			lo = ((lo - 0x06) & 0x0F) - 0x10
		}

		a = (old & 0xF0) - (rhs & 0xF0) + lo

		if a&0x100 != 0 {
			a -= 0x60
		}

		res := a - rhs + (^c & 1)

		// zero
		m.setZ(uint8(res))

		// negative
		m.setN(uint8(res))

		// carry
		if res <= 0xff {
			m.P |= flagC
		} else {
			m.P &= ^flagC
		}

		// overflow
		if (old^rhs)&(old^res)&0x80 != 0 {
			m.P |= flagV
		} else {
			m.P &= ^flagV
		}
		m.A = uint8(a)

	} else {
		m.opADC(^v)
	}
}

// opBit
func (m *M6502) opBit(v uint8) {
	m.P &= ^flagNVZ
	if v&(1<<7) != 0 {
		m.P |= flagN
	}
	if v&(1<<6) != 0 {
		m.P |= flagV
	}
	if v&m.A == 0 {
		m.P |= flagZ
	}
}

// opASL arithmetic shift left
func (m *M6502) opASL(v uint8) uint8 {
	if v&0x80 != 0 {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
	v <<= 1
	m.setNZ(v)
	return v
}

// opROL rotate left
func (m *M6502) opROL(v uint8) uint8 {
	ci := uint8(m.P & flagC)
	if v&0x80 != 0 {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
	v <<= 1
	v |= ci
	m.setNZ(v)
	return v
}

// opROR rotate right
func (m *M6502) opROR(v uint8) uint8 {
	ci := uint8(m.P & flagC)
	if v&1 != 0 {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
	v >>= 1
	v |= ci << 7
	m.setNZ(v)
	return v
}

// opLSR logical shift right
func (m *M6502) opLSR(v uint8) uint8 {
	if v&1 != 0 {
		m.P |= flagC
	} else {
		m.P &= ^flagC
	}
	v >>= 1
	m.setNZ(v)
	return v
}

//-----------------------------------------------------------------------------
// instructions

// op61, ADC add with carry, X-indexed indirect
func op61(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.opADC(v)
	m.PC += 2
	return 6
}

// op65, ADC add with carry, zeropage
func op65(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opADC(v)
	m.PC += 2
	return 3
}

// op69, ADC add with carry, immediate
func op69(m *M6502) uint {
	v := m.readImmediate()
	m.opADC(v)
	m.PC += 2
	return 2
}

// op6D, ADC add with carry, absolute
func op6D(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opADC(v)
	m.PC += 3
	return 4
}

// op71, ADC add with carry, indirect Y-indexed
func op71(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.opADC(v)
	m.PC += 2
	return 5 + n
}

// op75, ADC add with carry, zeropage X-indexed
func op75(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.opADC(v)
	m.PC += 2
	return 4
}

// op79, ADC add with carry, absolute Y-indexed
func op79(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.opADC(v)
	m.PC += 3
	return 4 + n
}

// op7D, ADC add with carry, absolute X-indexed
func op7D(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.opADC(v)
	m.PC += 3
	return 4 + n
}

// op21, AND and (with accumulator), X-indexed indirect
func op21(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 2
	return 6
}

// op25, AND and (with accumulator), zeropage
func op25(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 2
	return 3
}

// op29, AND and (with accumulator), immediate
func op29(m *M6502) uint {
	v := m.readImmediate()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 2
	return 2
}

// op2D, AND and (with accumulator), absolute
func op2D(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 3
	return 4
}

// op31, AND and (with accumulator), indirect Y-indexed
func op31(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 2
	return 5 + n
}

// op35, AND and (with accumulator), zeropage X-indexed
func op35(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 2
	return 4
}

// op39, AND and (with accumulator), absolute Y-indexed
func op39(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// op3D, AND and (with accumulator), absolute X-indexed
func op3D(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.A &= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// op06, ASL arithmetic shift left, zeropage
func op06(m *M6502) uint {
	v, ea := m.readZeroPage()
	v = m.opASL(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// op0A, ASL arithmetic shift left, accumulator
func op0A(m *M6502) uint {
	m.A = m.opASL(m.A)
	m.PC++
	return 2
}

// op0E, ASL arithmetic shift left, absolute
func op0E(m *M6502) uint {
	v, ea := m.readAbsolute()
	v = m.opASL(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// op16, ASL arithmetic shift left, zeropage X-indexed
func op16(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v = m.opASL(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// op1E, ASL arithmetic shift left, absolute X-indexed
func op1E(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v = m.opASL(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// op90, BCC branch on carry clear, relative
func op90(m *M6502) uint {
	return m.opBranch(m.P&flagC == 0)
}

// opB0, BCS branch on carry set, relative
func opB0(m *M6502) uint {
	return m.opBranch(m.P&flagC != 0)
}

// opF0, BEQ branch on equal (zero set), relative
func opF0(m *M6502) uint {
	return m.opBranch(m.P&flagZ != 0)
}

// op24, BIT bit test, zeropage
func op24(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opBit(v)
	m.PC += 2
	return 3
}

// op2C, BIT bit test, absolute
func op2C(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opBit(v)
	m.PC += 3
	return 4
}

// op30, BMI branch on minus (negative set), relative
func op30(m *M6502) uint {
	return m.opBranch(m.P&flagN != 0)
}

// opD0, BNE branch on not equal (zero clear), relative
func opD0(m *M6502) uint {
	return m.opBranch(m.P&flagZ == 0)
}

// op10, BPL branch on plus (negative clear), relative
func op10(m *M6502) uint {
	return m.opBranch(m.P&flagN == 0)
}

// op00, BRK break/interrupt
func op00(m *M6502) uint {
	m.push16(m.PC + 2)
	m.push8(m.P)
	m.P |= flagI
	m.PC = m.read16(BrkAddress)
	return 7
}

// op50, BVC branch on overflow clear, relative
func op50(m *M6502) uint {
	return m.opBranch(m.P&flagV == 0)
}

// op70, BVS branch on overflow set, relative
func op70(m *M6502) uint {
	return m.opBranch(m.P&flagV != 0)
}

// op18, CLC clear carry
func op18(m *M6502) uint {
	m.P &= ^flagC
	m.PC++
	return 2
}

// opD8, CLD clear decimal
func opD8(m *M6502) uint {
	m.P &= ^flagD
	m.PC++
	return 2
}

// op58, CLI clear interrupt disable
func op58(m *M6502) uint {
	m.P &= ^flagI
	m.PC++
	return 2
}

// opB8, CLV clear overflow
func opB8(m *M6502) uint {
	m.P &= ^flagV
	m.PC++
	return 2
}

// opC1, CMP compare (with accumulator), X-indexed indirect
func opC1(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.opCompare(m.A, v)
	m.PC += 2
	return 6
}

// opC5, CMP compare (with accumulator), zeropage
func opC5(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opCompare(m.A, v)
	m.PC += 2
	return 3
}

// opC9, CMP compare (with accumulator), immediate
func opC9(m *M6502) uint {
	v := m.readImmediate()
	m.opCompare(m.A, v)
	m.PC += 2
	return 2
}

// opCD, CMP compare (with accumulator), absolute
func opCD(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opCompare(m.A, v)
	m.PC += 3
	return 4
}

// opD1, CMP compare (with accumulator), indirect Y-indexed
func opD1(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.opCompare(m.A, v)
	m.PC += 2
	return 5 + n
}

// opD5, CMP compare (with accumulator), zeropage X-indexed
func opD5(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.opCompare(m.A, v)
	m.PC += 2
	return 4
}

// opD9, CMP compare (with accumulator), absolute Y-indexed
func opD9(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.opCompare(m.A, v)
	m.PC += 3
	return 4 + n
}

// opDD, CMP compare (with accumulator), absolute X-indexed
func opDD(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.opCompare(m.A, v)
	m.PC += 3
	return 4 + n
}

// opE0, CPX compare with X, immediate
func opE0(m *M6502) uint {
	v := m.readImmediate()
	m.opCompare(m.X, v)
	m.PC += 2
	return 2
}

// opE4, CPX compare with X, zeropage
func opE4(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opCompare(m.X, v)
	m.PC += 2
	return 3
}

// opEC, CPX compare with X, absolute
func opEC(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opCompare(m.X, v)
	m.PC += 3
	return 4
}

// opC0, CPY compare with Y, immediate
func opC0(m *M6502) uint {
	v := m.readImmediate()
	m.opCompare(m.Y, v)
	m.PC += 2
	return 2
}

// opC4, CPY compare with Y, zeropage
func opC4(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opCompare(m.Y, v)
	m.PC += 2
	return 3
}

// opCC, CPY compare with Y, absolute
func opCC(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opCompare(m.Y, v)
	m.PC += 3
	return 4
}

// opC6, DEC decrement, zeropage
func opC6(m *M6502) uint {
	v, ea := m.readZeroPage()
	v--
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// opCE, DEC decrement, absolute
func opCE(m *M6502) uint {
	v, ea := m.readAbsolute()
	v--
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// opD6, DEC decrement, zeropage X-indexed
func opD6(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v--
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// opDE, DEC decrement, absolute X-indexed
func opDE(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v--
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// opCA, DEX decrement X
func opCA(m *M6502) uint {
	m.X--
	m.setNZ(m.X)
	m.PC++
	return 2
}

// op88, DEY decrement Y
func op88(m *M6502) uint {
	m.Y--
	m.setNZ(m.Y)
	m.PC++
	return 2
}

// op41, EOR exclusive or (with accumulator), X-indexed indirect
func op41(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 2
	return 6
}

// op45, EOR exclusive or (with accumulator), zeropage
func op45(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 2
	return 3
}

// op49, EOR exclusive or (with accumulator), immediate
func op49(m *M6502) uint {
	v := m.readImmediate()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 2
	return 2
}

// op4D, EOR exclusive or (with accumulator), absolute
func op4D(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 3
	return 4
}

// op51, EOR exclusive or (with accumulator), indirect Y-indexed
func op51(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 2
	return 5 + n
}

// op55, EOR exclusive or (with accumulator), zeropage X-indexed
func op55(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 2
	return 4
}

// op59, EOR exclusive or (with accumulator), absolute Y-indexed
func op59(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// op5D, EOR exclusive or (with accumulator), absolute X-indexed
func op5D(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.A ^= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// opXX, ILL illegal
func opXX(m *M6502) uint {
	m.illegal = true
	return 0
}

// opE6, INC increment, zeropage
func opE6(m *M6502) uint {
	v, ea := m.readZeroPage()
	v++
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// opEE, INC increment, absolute
func opEE(m *M6502) uint {
	v, ea := m.readAbsolute()
	v++
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// opF6, INC increment, zeropage X-indexed
func opF6(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v++
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// opFE, INC increment, absolute X-indexed
func opFE(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v++
	m.setNZ(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// opE8, INX increment X
func opE8(m *M6502) uint {
	m.X++
	m.setNZ(m.X)
	m.PC++
	return 2
}

// opC8, INY increment Y
func opC8(m *M6502) uint {
	m.Y++
	m.setNZ(m.Y)
	m.PC++
	return 2
}

// op4C, JMP jump, absolute
func op4C(m *M6502) uint {
	m.PC = m.read16(m.PC + 1)
	m.jmpVSR()
	return 3
}

// op6C, JMP jump, indirect
func op6C(m *M6502) uint {
	m.PC = m.read16(m.read16(m.PC + 1))
	m.jmpVSR()
	return 5
}

// op20, JSR jump subroutine, absolute
func op20(m *M6502) uint {
	m.push16(m.PC + 2)
	m.PC = m.read16(m.PC + 1)
	m.jsrVSR()
	return 6
}

// opA1, LDA load accumulator, X-indexed indirect
func opA1(m *M6502) uint {
	m.A, _ = m.readIndirectX()
	m.setNZ(m.A)
	m.PC += 2
	return 6
}

// opA5, LDA load accumulator, zeropage
func opA5(m *M6502) uint {
	m.A, _ = m.readZeroPage()
	m.setNZ(m.A)
	m.PC += 2
	return 3
}

// opA9, LDA load accumulator, immediate
func opA9(m *M6502) uint {
	m.A = m.readImmediate()
	m.setNZ(m.A)
	m.PC += 2
	return 2
}

// opAD, LDA load accumulator, absolute
func opAD(m *M6502) uint {
	m.A, _ = m.readAbsolute()
	m.setNZ(m.A)
	m.PC += 3
	return 4
}

// opB1, LDA load accumulator, indirect Y-indexed
func opB1(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.A = v
	m.setNZ(m.A)
	m.PC += 2
	return 5 + n
}

// opB5, LDA load accumulator, zeropage X-indexed
func opB5(m *M6502) uint {
	m.A, _ = m.readZeroPageX()
	m.setNZ(m.A)
	m.PC += 2
	return 4
}

// opB9, LDA load accumulator, absolute Y-indexed
func opB9(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.A = v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// opBD, LDA load accumulator, absolute X-indexed
func opBD(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.A = v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// opA2, LDX load X, immediate
func opA2(m *M6502) uint {
	m.X = m.readImmediate()
	m.setNZ(m.X)
	m.PC += 2
	return 2
}

// opA6, LDX load X, zeropage
func opA6(m *M6502) uint {
	m.X, _ = m.readZeroPage()
	m.setNZ(m.X)
	m.PC += 2
	return 3
}

// opAE, LDX load X, absolute
func opAE(m *M6502) uint {
	m.X, _ = m.readAbsolute()
	m.setNZ(m.X)
	m.PC += 3
	return 4
}

// opB6, LDX load X, zeropage Y-indexed
func opB6(m *M6502) uint {
	m.X, _ = m.readZeroPageY()
	m.setNZ(m.X)
	m.PC += 2
	return 4
}

// opBE, LDX load X, absolute Y-indexed
func opBE(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.X = v
	m.setNZ(m.X)
	m.PC += 3
	return 4 + n
}

// opA0, LDY load Y, immediate
func opA0(m *M6502) uint {
	v := m.readImmediate()
	m.Y = v
	m.setNZ(m.Y)
	m.PC += 2
	return 2
}

// opA4, LDY load Y, zeropage
func opA4(m *M6502) uint {
	m.Y, _ = m.readZeroPage()
	m.setNZ(m.Y)
	m.PC += 2
	return 3
}

// opAC, LDY load Y, absolute
func opAC(m *M6502) uint {
	m.Y, _ = m.readAbsolute()
	m.setNZ(m.Y)
	m.PC += 3
	return 4
}

// opB4, LDY load Y, zeropage X-indexed
func opB4(m *M6502) uint {
	m.Y, _ = m.readZeroPageX()
	m.setNZ(m.Y)
	m.PC += 2
	return 4
}

// opBC, LDY load Y, absolute X-indexed
func opBC(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.Y = v
	m.setNZ(m.Y)
	m.PC += 3
	return 4 + n
}

// op46, LSR logical shift right, zeropage
func op46(m *M6502) uint {
	v, ea := m.readZeroPage()
	v = m.opLSR(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// op4A, LSR logical shift right, accumulator
func op4A(m *M6502) uint {
	m.A = m.opLSR(m.A)
	m.PC++
	return 2
}

// op4E, LSR logical shift right, absolute
func op4E(m *M6502) uint {
	v, ea := m.readAbsolute()
	v = m.opLSR(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// op56, LSR logical shift right, zeropage X-indexed
func op56(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v = m.opLSR(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// op5E, LSR logical shift right, absolute X-indexed
func op5E(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v = m.opLSR(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// opEA, NOP no operation
func opEA(m *M6502) uint {
	m.PC++
	return 2
}

// op01, ORA or with accumulator, X-indexed indirect
func op01(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 2
	return 6
}

// op05, ORA or with accumulator, zeropage
func op05(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 2
	return 3
}

// op09, ORA or with accumulator, immediate
func op09(m *M6502) uint {
	v := m.readImmediate()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 2
	return 2
}

// op0D, ORA or with accumulator, absolute
func op0D(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 3
	return 4
}

// op11, ORA or with accumulator, indirect Y-indexed
func op11(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 2
	return 5 + n
}

// op15, ORA or with accumulator, zeropage X-indexed
func op15(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 2
	return 4
}

// op19, ORA or with accumulator, absolute Y-indexed
func op19(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// op1D, ORA or with accumulator, absolute X-indexed
func op1D(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.A |= v
	m.setNZ(m.A)
	m.PC += 3
	return 4 + n
}

// op48, PHA push accumulator
func op48(m *M6502) uint {
	m.push8(m.A)
	m.PC++
	return 3
}

// op08, PHP push processor status (SR)
func op08(m *M6502) uint {
	m.push8(m.P)
	m.PC++
	return 3
}

// op68, PLA pull accumulator
func op68(m *M6502) uint {
	m.A = m.pop8()
	m.setNZ(m.A)
	m.PC++
	return 4
}

// op28, PLP pull processor status (SR)
func op28(m *M6502) uint {
	m.P = m.pop8() | flagB | (1 << 5)
	m.PC++
	return 4
}

// op26, ROL rotate left, zeropage
func op26(m *M6502) uint {
	v, ea := m.readZeroPage()
	v = m.opROL(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// op2A, ROL rotate left, accumulator
func op2A(m *M6502) uint {
	m.A = m.opROL(m.A)
	m.PC++
	return 2
}

// op2E, ROL rotate left, absolute
func op2E(m *M6502) uint {
	v, ea := m.readAbsolute()
	v = m.opROL(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// op36, ROL rotate left, zeropage X-indexed
func op36(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v = m.opROL(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// op3E, ROL rotate left, absolute X-indexed
func op3E(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v = m.opROL(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// op66, ROR rotate right, zeropage
func op66(m *M6502) uint {
	v, ea := m.readZeroPage()
	v = m.opROR(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 5
}

// op6A, ROR rotate right, accumulator
func op6A(m *M6502) uint {
	m.A = m.opROR(m.A)
	m.PC++
	return 2
}

// op6E, ROR rotate right, absolute
func op6E(m *M6502) uint {
	v, ea := m.readAbsolute()
	v = m.opROR(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 6
}

// op76, ROR rotate right, zeropage X-indexed
func op76(m *M6502) uint {
	v, ea := m.readZeroPageX()
	v = m.opROR(v)
	m.Mem.Write8(ea, v)
	m.PC += 2
	return 6
}

// op7E, ROR rotate right, absolute X-indexed
func op7E(m *M6502) uint {
	v, ea := m.readAbsoluteX()
	v = m.opROR(v)
	m.Mem.Write8(ea, v)
	m.PC += 3
	return 7
}

// op40, RTI return from interrupt
func op40(m *M6502) uint {
	m.P = m.pop8()
	m.PC = m.pop16()
	return 6
}

// op60, RTS return from subroutine
func op60(m *M6502) uint {
	m.PC = m.pop16() + 1
	return 6
}

// opE1, SBC subtract with carry, X-indexed indirect
func opE1(m *M6502) uint {
	v, _ := m.readIndirectX()
	m.opSBC(v)
	m.PC += 2
	return 6
}

// opE5, SBC subtract with carry, zeropage
func opE5(m *M6502) uint {
	v, _ := m.readZeroPage()
	m.opSBC(v)
	m.PC += 2
	return 3
}

// opE9, SBC subtract with carry, immediate
func opE9(m *M6502) uint {
	v := m.readImmediate()
	m.opSBC(v)
	m.PC += 2
	return 2
}

// opED, SBC subtract with carry, absolute
func opED(m *M6502) uint {
	v, _ := m.readAbsolute()
	m.opSBC(v)
	m.PC += 3
	return 4
}

// opF1, SBC subtract with carry, indirect Y-indexed
func opF1(m *M6502) uint {
	v, n, _ := m.readIndirectYPenalized()
	m.opSBC(v)
	m.PC += 2
	return 5 + n
}

// opF5, SBC subtract with carry, zeropage X-indexed
func opF5(m *M6502) uint {
	v, _ := m.readZeroPageX()
	m.opSBC(v)
	m.PC += 2
	return 4
}

// opF9, SBC subtract with carry, absolute Y-indexed
func opF9(m *M6502) uint {
	v, n, _ := m.readAbsoluteYPenalized()
	m.opSBC(v)
	m.PC += 3
	return 4 + n
}

// opFD, SBC subtract with carry, absolute X-indexed
func opFD(m *M6502) uint {
	v, n, _ := m.readAbsoluteXPenalized()
	m.opSBC(v)
	m.PC += 3
	return 4 + n
}

// op38, SEC set carry
func op38(m *M6502) uint {
	m.PC++
	m.P |= flagC
	return 2
}

// opF8, SED set decimal
func opF8(m *M6502) uint {
	m.PC++
	m.P |= flagD
	return 2
}

// op78, SEI set interrupt disable
func op78(m *M6502) uint {
	m.PC++
	m.P |= flagI
	return 2
}

// op81, STA store accumulator, X-indexed indirect
func op81(m *M6502) uint {
	m.writeIndirectX(m.A)
	m.PC += 2
	return 6
}

// op85, STA store accumulator, zeropage
func op85(m *M6502) uint {
	m.writeZeroPage(m.A)
	m.PC += 2
	return 3
}

// op8D, STA store accumulator, absolute
func op8D(m *M6502) uint {
	m.writeAbsolute(m.A)
	m.PC += 3
	return 4
}

// op91, STA store accumulator, indirect Y-indexed
func op91(m *M6502) uint {
	m.writeIndirectY(m.A)
	m.PC += 2
	return 6
}

// op95, STA store accumulator, zeropage X-indexed
func op95(m *M6502) uint {
	m.writeZeroPageX(m.A)
	m.PC += 2
	return 4
}

// op99, STA store accumulator, absolute Y-indexed
func op99(m *M6502) uint {
	m.writeAbsoluteY(m.A)
	m.PC += 3
	return 5
}

// op9D, STA store accumulator, absolute X-indexed
func op9D(m *M6502) uint {
	m.writeAbsoluteX(m.A)
	m.PC += 3
	return 5
}

// op86, STX store X, zeropage
func op86(m *M6502) uint {
	m.writeZeroPage(m.X)
	m.PC += 2
	return 3
}

// op8E, STX store X, absolute
func op8E(m *M6502) uint {
	m.writeAbsolute(m.X)
	m.PC += 3
	return 4
}

// op96, STX store X, zeropage Y-indexed
func op96(m *M6502) uint {
	m.writeZeroPageY(m.X)
	m.PC += 2
	return 4
}

// op84, STY store Y, zeropage
func op84(m *M6502) uint {
	m.writeZeroPage(m.Y)
	m.PC += 2
	return 3
}

// op8C, STY store Y, absolute
func op8C(m *M6502) uint {
	m.writeAbsolute(m.Y)
	m.PC += 3
	return 4
}

// op94, STY store Y, zeropage X-indexed
func op94(m *M6502) uint {
	m.writeZeroPageX(m.Y)
	m.PC += 2
	return 4
}

// opAA, TAX transfer accumulator to X
func opAA(m *M6502) uint {
	m.PC++
	m.X = m.A
	m.setNZ(m.X)
	return 2
}

// opA8, TAY transfer accumulator to Y
func opA8(m *M6502) uint {
	m.PC++
	m.Y = m.A
	m.setNZ(m.Y)
	return 2
}

// opBA, TSX transfer stack pointer to X
func opBA(m *M6502) uint {
	m.PC++
	m.X = m.S
	m.setNZ(m.X)
	return 2
}

// op8A, TXA transfer X to accumulator
func op8A(m *M6502) uint {
	m.PC++
	m.A = m.X
	m.setNZ(m.A)
	return 2
}

// op9A, TXS transfer X to stack pointer
func op9A(m *M6502) uint {
	m.PC++
	m.S = m.X
	return 2
}

// op98, TYA transfer Y to accumulator
func op98(m *M6502) uint {
	m.PC++
	m.A = m.Y
	m.setNZ(m.A)
	return 2
}

type opFunc func(m *M6502) uint

var opcodeTable = [256]opFunc{
	op00, op01, opXX, opXX, opXX, op05, op06, opXX, op08, op09, op0A, opXX, opXX, op0D, op0E, opXX,
	op10, op11, opXX, opXX, opXX, op15, op16, opXX, op18, op19, opXX, opXX, opXX, op1D, op1E, opXX,
	op20, op21, opXX, opXX, op24, op25, op26, opXX, op28, op29, op2A, opXX, op2C, op2D, op2E, opXX,
	op30, op31, opXX, opXX, opXX, op35, op36, opXX, op38, op39, opXX, opXX, opXX, op3D, op3E, opXX,
	op40, op41, opXX, opXX, opXX, op45, op46, opXX, op48, op49, op4A, opXX, op4C, op4D, op4E, opXX,
	op50, op51, opXX, opXX, opXX, op55, op56, opXX, op58, op59, opXX, opXX, opXX, op5D, op5E, opXX,
	op60, op61, opXX, opXX, opXX, op65, op66, opXX, op68, op69, op6A, opXX, op6C, op6D, op6E, opXX,
	op70, op71, opXX, opXX, opXX, op75, op76, opXX, op78, op79, opXX, opXX, opXX, op7D, op7E, opXX,
	opXX, op81, opXX, opXX, op84, op85, op86, opXX, op88, opXX, op8A, opXX, op8C, op8D, op8E, opXX,
	op90, op91, opXX, opXX, op94, op95, op96, opXX, op98, op99, op9A, opXX, opXX, op9D, opXX, opXX,
	opA0, opA1, opA2, opXX, opA4, opA5, opA6, opXX, opA8, opA9, opAA, opXX, opAC, opAD, opAE, opXX,
	opB0, opB1, opXX, opXX, opB4, opB5, opB6, opXX, opB8, opB9, opBA, opXX, opBC, opBD, opBE, opXX,
	opC0, opC1, opXX, opXX, opC4, opC5, opC6, opXX, opC8, opC9, opCA, opXX, opCC, opCD, opCE, opXX,
	opD0, opD1, opXX, opXX, opXX, opD5, opD6, opXX, opD8, opD9, opXX, opXX, opXX, opDD, opDE, opXX,
	opE0, opE1, opXX, opXX, opE4, opE5, opE6, opXX, opE8, opE9, opEA, opXX, opEC, opED, opEE, opXX,
	opF0, opF1, opXX, opXX, opXX, opF5, opF6, opXX, opF8, opF9, opXX, opXX, opXX, opFD, opFE, opXX,
}

//-----------------------------------------------------------------------------

// New6502 returns a 6502 CPU in the powered-on and reset state.
func New6502(mem Memory) *M6502 {
	var m M6502
	m.Mem = mem

	// initialise the opcode usage map
	m.usage = make(map[uint8]uint)
	for i := 0; i < 256; i++ {
		x := opcodeLookup(uint8(i))
		if x.ins != "ill" {
			m.usage[uint8(i)] = 0
		}
	}

	return &m
}

// Power on/off the 6502 CPU.
func (m *M6502) Power(state bool) {
	if state {
		m.PC = initialPC
		m.S = initialS
		m.P = initialP
		m.A = initialA
		m.X = initialX
		m.Y = initialY
	} else {
		m.PC = 0
		m.S = 0
		m.P = 0
		m.A = 0
		m.X = 0
		m.Y = 0
	}
	m.irq = false
	m.nmi = false
	m.illegal = false
	m.exit = false
	m.cycles = 0
	m.lastPC = 0
	m.stuckPC = 0
}

// Reset the 6502 CPU.
func (m *M6502) Reset() {
	m.PC = m.read16(RstAddress)
	m.S = initialS
	m.P = initialP
	m.irq = false
	m.nmi = false
}

// NMI generates a non-maskable-interrupt.
func (m *M6502) NMI() {
	m.nmi = true
}

// IRQ generates an interrupt request.
func (m *M6502) IRQ(state bool) {
	m.irq = state
}

// Run the 6502 CPU for a single instruction.
func (m *M6502) Run() error {
	// nmi handling
	if m.nmi {
		m.nmi = false
		m.push16(m.PC)
		m.push8(m.P & ^flagB)
		m.P |= flagI
		m.PC = m.read16(NmiAddress)
		m.cycles += 7
		return nil
	}
	// irq handling
	if m.irq && (m.P&flagI == 0) {
		m.irq = false
		m.push16(m.PC)
		m.push8(m.P & ^flagB)
		m.P |= flagI
		m.PC = m.read16(IrqAddress)
		m.cycles += 7
		return nil
	}
	// normal instructions
	op := m.Mem.Read8(m.PC)
	m.cycles += opcodeTable[op](m)

	if m.illegal {
		return fmt.Errorf("illegal instruction at %04x", m.PC)
	}

	if m.exit {
		return fmt.Errorf("exit at %04x, status %02x, %d cpu cycles, %.2f coverage", m.PC, m.A, m.cycles, m.Coverage())
	}

	// accumulate opcode usage
	m.usage[op]++

	// stuck PC detection
	if m.PC == m.lastPC {
		m.stuckPC++
		if m.stuckPC >= 4 {
			return fmt.Errorf("PC is stuck at %04x, %d cpu cycles", m.PC, m.cycles)
		}
	} else {
		m.stuckPC = 0
		m.lastPC = m.PC
	}

	return nil
}

// Coverage returns the fraction of valid opcodes that have run.
func (m *M6502) Coverage() float32 {
	run := 0
	for _, v := range m.usage {
		if v != 0 {
			run++
		}
	}
	return float32(run) / float32(len(m.usage))
}

// ReadPC returns the 6502 program counter.
func (m *M6502) ReadPC() uint16 {
	return m.PC
}

// Exit sets a status code and exits the emulation
func (m *M6502) Exit(status uint8) {
	m.A = status
	m.exit = true
}

//-----------------------------------------------------------------------------
// virtual JSR/JMP subroutines

func (m *M6502) jsrVSR() {
	if m.vsr != nil {
		if fn, ok := m.vsr[m.PC]; ok {
			// call the hook
			fn(m)
			// simulate RTS
			m.PC = m.pop16() + 1
		}
	}
}

func (m *M6502) jmpVSR() {
	if m.vsr != nil {
		if fn, ok := m.vsr[m.PC]; ok {
			// call the hook
			fn(m)
			// called by a jump: we don't have anywhere to go...
		}
	}
}

// AddVSR adds a virtual subroutine handler at the call address.
func (m *M6502) AddVSR(adr uint16, fn VSRFunc) {
	if m.vsr == nil {
		m.vsr = make(map[uint16]VSRFunc)
	}
	m.vsr[adr] = fn
}

//-----------------------------------------------------------------------------
