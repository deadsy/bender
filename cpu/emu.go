//-----------------------------------------------------------------------------
/*

6502 CPU Emulator

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

func emuTODO() {
	panic("TODO")
}

//-----------------------------------------------------------------------------

func (m *M6502) setNZ(val uint8) {
	var flags uint8
	if val != 0 {
		flags = val & flagN
	} else {
		flags = flagZ
	}
	m.P = (m.P &^ flagNZ) | flags
}

//-----------------------------------------------------------------------------

func (m *M6502) read8(adr uint16) uint8 {
	return m.mem.Read8(adr)
}

func (m *M6502) read16(adr uint16) uint16 {
	l := uint16(m.mem.Read8(adr))
	h := uint16(m.mem.Read8(adr + 1))
	return (h << 8) | l
}

func (m *M6502) readPointer(adr uint16) uint16 {
	return m.read16(adr)
}

//-----------------------------------------------------------------------------

func (m *M6502) push8(val uint8) {
	m.mem.Write8(stkAddress+uint16(m.S), val)
	m.S--
}

func (m *M6502) pop8() uint8 {
	m.S++
	return m.mem.Read8(stkAddress + uint16(m.S))
}

func (m *M6502) push16(val uint16) {
	m.mem.Write8(stkAddress+uint16(m.S), uint8(val>>8))
	m.mem.Write8(stkAddress+uint16(m.S-1), uint8(val))
	m.S -= 2
}

func (m *M6502) pop16() uint16 {
	l := uint16(m.mem.Read8(stkAddress + uint16(m.S+1)))
	h := uint16(m.mem.Read8(stkAddress + uint16(m.S+2)))
	m.S += 2
	return (h << 8) | l
}

//-----------------------------------------------------------------------------

func (m *M6502) readIndirectX() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readZeroPage() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readImmediate() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readAbsolute() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readIndirectY() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readZeroPageX() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readAbsoluteY() (uint8, uint) {
	return 0, 0
}

func (m *M6502) readAbsoluteX() (uint8, uint) {
	return 0, 0
}

//-----------------------------------------------------------------------------
// ADC add with carry

func (m *M6502) opADC(v uint8) {
	c := m.P & flagC
	if m.P&flagD != 0 {
		l := uint(m.A&0x0F) + uint(v&0x0F) + uint(c)
		h := uint(m.A&0xF0) + uint(v&0xF0)
		m.P &= ^(flagV | flagC | flagN | flagZ)
		if (l+h)&0xFF == 0 {
			m.P |= flagZ
		}
		if l > 0x09 {
			h += 0x10
			l += 0x06
		}
		if h&0x80 != 0 {
			m.P |= flagN
		}
		if ^(m.A^v)&(m.A^uint8(h))&0x80 != 0 {
			m.P |= flagV
		}
		if h > 0x90 {
			h += 0x60
		}
		if h>>8 != 0 {
			m.P |= flagC
		}
		m.A = uint8(l&0x0F) | uint8(h&0xF0)
	} else {
		t := uint(m.A) + uint(v) + uint(c)
		m.P &= ^(flagV | flagC)
		if ^(m.A^v)&(m.A^uint8(t))&0x80 != 0 {
			m.P |= flagV
		}
		if t>>8 != 0 {
			m.P |= flagC
		}
		m.A = uint8(t)
		m.setNZ(m.A)
	}
}

// opADCabs, add with carry, absolute mode
func opADCabs(m *M6502) uint {
	v, n := m.readAbsolute()
	m.opADC(v)
	return n
}

// opADCabsx, add with carry, absolute, X-indexed mode
func opADCabsx(m *M6502) uint {
	v, n := m.readAbsoluteX()
	m.opADC(v)
	return n
}

// opADCabsy, add with carry, absolute, Y-indexed mode
func opADCabsy(m *M6502) uint {
	v, n := m.readAbsoluteY()
	m.opADC(v)
	return n
}

// opADCimm, add with carry, immediate mode
func opADCimm(m *M6502) uint {
	v, n := m.readImmediate()
	m.opADC(v)
	return n
}

// opADCindy, add with carry, indirect, Y-indexed mode
func opADCindy(m *M6502) uint {
	v, n := m.readIndirectY()
	m.opADC(v)
	return n
}

// opADCxind, add with carry, X-indexed, indirect mode
func opADCxind(m *M6502) uint {
	v, n := m.readIndirectX()
	m.opADC(v)
	return n
}

// opADCz, add with carry, zeropage mode
func opADCz(m *M6502) uint {
	v, n := m.readZeroPage()
	m.opADC(v)
	return n
}

// opADCzx, add with carry, zeropage, X-indexed mode
func opADCzx(m *M6502) uint {
	v, n := m.readZeroPageX()
	m.opADC(v)
	return n
}

//-----------------------------------------------------------------------------

// opANDabs, and (with accumulator), absolute mode
func opANDabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDabsx, and (with accumulator), absolute, X-indexed mode
func opANDabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDabsy, and (with accumulator), absolute, Y-indexed mode
func opANDabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDimm, and (with accumulator), immediate mode
func opANDimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDindy, and (with accumulator), indirect, Y-indexed mode
func opANDindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDxind, and (with accumulator), X-indexed, indirect mode
func opANDxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDz, and (with accumulator), zeropage mode
func opANDz(m *M6502) uint {
	emuTODO()
	return 0
}

// opANDzx, and (with accumulator), zeropage, X-indexed mode
func opANDzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opASLabs, arithmetic shift left, absolute mode
func opASLabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opASLabsx, arithmetic shift left, absolute, X-indexed mode
func opASLabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opASLacc, arithmetic shift left, accumulator mode
func opASLacc(m *M6502) uint {
	emuTODO()
	return 0
}

// opASLz, arithmetic shift left, zeropage mode
func opASLz(m *M6502) uint {
	emuTODO()
	return 0
}

// opASLzx, arithmetic shift left, zeropage, X-indexed mode
func opASLzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opBCCrel, branch on carry clear, relative mode
func opBCCrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBCSrel, branch on carry set, relative mode
func opBCSrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBEQimm, branch on equal (zero set), immediate mode
func opBEQimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opBITabs, bit test, absolute mode
func opBITabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opBITz, bit test, zeropage mode
func opBITz(m *M6502) uint {
	emuTODO()
	return 0
}

// opBMIrel, branch on minus (negative set), relative mode
func opBMIrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBNErel, branch on not equal (zero clear), relative mode
func opBNErel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBPLrel, branch on plus (negative clear), relative mode
func opBPLrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBRKimpl, break / interrupt, implied mode
func opBRKimpl(m *M6502) uint {
	m.read8(m.PC + 1)
	m.push16(m.PC + 2)
	m.push8(m.P | flagB)
	m.P |= flagB | flagI
	m.PC = m.readPointer(brkAddress)
	return 7
}

// opBVCrel, branch on overflow clear, relative mode
func opBVCrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opBVSrel, branch on overflow set, relative mode
func opBVSrel(m *M6502) uint {
	emuTODO()
	return 0
}

// opCLCimpl, clear carry, implied mode
func opCLCimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opCLDimpl, clear decimal, implied mode
func opCLDimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opCLIimpl, clear interrupt disable, implied mode
func opCLIimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opCLVimpl, clear overflow, implied mode
func opCLVimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPabs, compare (with accumulator), absolute mode
func opCMPabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPabsx, compare (with accumulator), absolute, X-indexed mode
func opCMPabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPabsy, compare (with accumulator), absolute, Y-indexed mode
func opCMPabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPimm, compare (with accumulator), immediate mode
func opCMPimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPindy, compare (with accumulator), indirect, Y-indexed mode
func opCMPindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPxind, compare (with accumulator), X-indexed, indirect mode
func opCMPxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPz, compare (with accumulator), zeropage mode
func opCMPz(m *M6502) uint {
	emuTODO()
	return 0
}

// opCMPzx, compare (with accumulator), zeropage, X-indexed mode
func opCMPzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPXabs, compare with X, absolute mode
func opCPXabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPXimm, compare with X, immediate mode
func opCPXimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPXz, compare with X, zeropage mode
func opCPXz(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPYabs, compare with Y, absolute mode
func opCPYabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPYimm, compare with Y, immediate mode
func opCPYimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opCPYz, compare with Y, zeropage mode
func opCPYz(m *M6502) uint {
	emuTODO()
	return 0
}

// opDECabs, decrement, absolute mode
func opDECabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opDECabsx, decrement, absolute, X-indexed mode
func opDECabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opDECz, decrement, zeropage mode
func opDECz(m *M6502) uint {
	emuTODO()
	return 0
}

// opDECzx, decrement, zeropage, X-indexed mode
func opDECzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opDEXimpl, decrement X, implied mode
func opDEXimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opDEYimpl, decrement Y, implied mode
func opDEYimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORabs, exclusive or (with accumulator), absolute mode
func opEORabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORabsx, exclusive or (with accumulator), absolute, X-indexed mode
func opEORabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORabsy, exclusive or (with accumulator), absolute, Y-indexed mode
func opEORabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORimm, exclusive or (with accumulator), immediate mode
func opEORimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORindy, exclusive or (with accumulator), indirect, Y-indexed mode
func opEORindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORxind, exclusive or (with accumulator), X-indexed, indirect mode
func opEORxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORz, exclusive or (with accumulator), zeropage mode
func opEORz(m *M6502) uint {
	emuTODO()
	return 0
}

// opEORzx, exclusive or (with accumulator), zeropage, X-indexed mode
func opEORzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opILL, illegal instruction
func opILL(m *M6502) uint {
	return 2
}

// opINCabs, increment, absolute mode
func opINCabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opINCabsx, increment, absolute, X-indexed mode
func opINCabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opINCz, increment, zeropage mode
func opINCz(m *M6502) uint {
	emuTODO()
	return 0
}

// opINCzx, increment, zeropage, X-indexed mode
func opINCzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opINXimpl, increment X, implied mode
func opINXimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opINYimpl, increment Y, implied mode
func opINYimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opJMPabs, jump, absolute mode
func opJMPabs(m *M6502) uint {
	m.PC = m.read16(m.PC + 1)
	return 3
}

// opJMPind, jump, indirect mode
func opJMPind(m *M6502) uint {
	m.PC = m.read16(m.read16(m.PC + 1))
	return 5
}

// opJSRabs, jump subroutine, absolute mode
func opJSRabs(m *M6502) uint {
	m.push16(m.PC + 2)
	m.PC = m.read16(m.PC + 1)
	return 6
}

// opLDAabs, load accumulator, absolute mode
func opLDAabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAabsx, load accumulator, absolute, X-indexed mode
func opLDAabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAabsy, load accumulator, absolute, Y-indexed mode
func opLDAabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAimm, load accumulator, immediate mode
func opLDAimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAindy, load accumulator, indirect, Y-indexed mode
func opLDAindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAxind, load accumulator, X-indexed, indirect mode
func opLDAxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAz, load accumulator, zeropage mode
func opLDAz(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDAzx, load accumulator, zeropage, X-indexed mode
func opLDAzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDXabs, load X, absolute mode
func opLDXabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDXabsy, load X, absolute, Y-indexed mode
func opLDXabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDXimm, load X, immediate mode
func opLDXimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDXz, load X, zeropage mode
func opLDXz(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDXzy, load X, zeropage, Y-indexed mode
func opLDXzy(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDYabs, load Y, absolute mode
func opLDYabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDYabsx, load Y, absolute, X-indexed mode
func opLDYabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDYimm, load Y, immediate mode
func opLDYimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDYz, load Y, zeropage mode
func opLDYz(m *M6502) uint {
	emuTODO()
	return 0
}

// opLDYzx, load Y, zeropage, X-indexed mode
func opLDYzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opLSRabs, logical shift right, absolute mode
func opLSRabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opLSRabsx, logical shift right, absolute, X-indexed mode
func opLSRabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opLSRacc, logical shift right, accumulator mode
func opLSRacc(m *M6502) uint {
	emuTODO()
	return 0
}

// opLSRz, logical shift right, zeropage mode
func opLSRz(m *M6502) uint {
	emuTODO()
	return 0
}

// opLSRzx, logical shift right, zeropage, X-indexed mode
func opLSRzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opNOPimpl, no operation, implied mode
func opNOPimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAabs, or with accumulator, absolute mode
func opORAabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAabsx, or with accumulator, absolute, X-indexed mode
func opORAabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAabsy, or with accumulator, absolute, Y-indexed mode
func opORAabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAimm, or with accumulator, immediate mode
func opORAimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAindy, or with accumulator, indirect, Y-indexed mode
func opORAindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAxind, or with accumulator, X-indexed, indirect mode
func opORAxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAz, or with accumulator, zeropage mode
func opORAz(m *M6502) uint {
	emuTODO()
	return 0
}

// opORAzx, or with accumulator, zeropage, X-indexed mode
func opORAzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opPHAimpl, push accumulator, implied mode
func opPHAimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opPHPimpl, push processor status (SR), implied mode
func opPHPimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opPLAimpl, pull accumulator, implied mode
func opPLAimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opPLPimpl, pull processor status (SR), implied mode
func opPLPimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opROLabs, rotate left, absolute mode
func opROLabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opROLabsx, rotate left, absolute, X-indexed mode
func opROLabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opROLacc, rotate left, accumulator mode
func opROLacc(m *M6502) uint {
	emuTODO()
	return 0
}

// opROLz, rotate left, zeropage mode
func opROLz(m *M6502) uint {
	emuTODO()
	return 0
}

// opROLzx, rotate left, zeropage, X-indexed mode
func opROLzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opRORabs, rotate right, absolute mode
func opRORabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opRORabsx, rotate right, absolute, X-indexed mode
func opRORabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opRORacc, rotate right, accumulator mode
func opRORacc(m *M6502) uint {
	emuTODO()
	return 0
}

// opRORz, rotate right, zeropage mode
func opRORz(m *M6502) uint {
	emuTODO()
	return 0
}

// opRORzx, rotate right, zeropage, X-indexed mode
func opRORzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opRTIimpl, return from interrupt, implied mode
func opRTIimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opRTSimpl, return from subroutine, implied mode
func opRTSimpl(m *M6502) uint {
	m.PC = m.pop16() + 1
	return 6
}

// opSBCabs, subtract with carry, absolute mode
func opSBCabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCabsx, subtract with carry, absolute, X-indexed mode
func opSBCabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCabsy, subtract with carry, absolute, Y-indexed mode
func opSBCabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCimm, subtract with carry, immediate mode
func opSBCimm(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCindy, subtract with carry, indirect, Y-indexed mode
func opSBCindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCxind, subtract with carry, X-indexed, indirect mode
func opSBCxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCz, subtract with carry, zeropage mode
func opSBCz(m *M6502) uint {
	emuTODO()
	return 0
}

// opSBCzx, subtract with carry, zeropage, X-indexed mode
func opSBCzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opSECimpl, set carry, implied mode
func opSECimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opSEDimpl, set decimal, implied mode
func opSEDimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opSEIimpl, set interrupt disable, implied mode
func opSEIimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAabs, store accumulator, absolute mode
func opSTAabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAabsx, store accumulator, absolute, X-indexed mode
func opSTAabsx(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAabsy, store accumulator, absolute, Y-indexed mode
func opSTAabsy(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAindy, store accumulator, indirect, Y-indexed mode
func opSTAindy(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAxind, store accumulator, X-indexed, indirect mode
func opSTAxind(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAz, store accumulator, zeropage mode
func opSTAz(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTAzx, store accumulator, zeropage, X-indexed mode
func opSTAzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTXabs, store X, absolute mode
func opSTXabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTXz, store X, zeropage mode
func opSTXz(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTXzy, store X, zeropage, Y-indexed mode
func opSTXzy(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTYabs, store Y, absolute mode
func opSTYabs(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTYz, store Y, zeropage mode
func opSTYz(m *M6502) uint {
	emuTODO()
	return 0
}

// opSTYzx, store Y, zeropage, X-indexed mode
func opSTYzx(m *M6502) uint {
	emuTODO()
	return 0
}

// opTAXimpl, transfer accumulator to X, implied mode
func opTAXimpl(m *M6502) uint {
	m.PC++
	m.X = m.A
	m.setNZ(m.X)
	return 2
}

// opTAYimpl, transfer accumulator to Y, implied mode
func opTAYimpl(m *M6502) uint {
	m.PC++
	m.Y = m.A
	m.setNZ(m.Y)
	return 2
}

// opTSXimpl, transfer stack pointer to X, implied mode
func opTSXimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opTXAimpl, transfer X to accumulator, implied mode
func opTXAimpl(m *M6502) uint {
	m.PC++
	m.A = m.X
	m.setNZ(m.A)
	return 2
}

// opTXSimpl, transfer X to stack pointer, implied mode
func opTXSimpl(m *M6502) uint {
	emuTODO()
	return 0
}

// opTYAimpl, transfer Y to accumulator, implied mode
func opTYAimpl(m *M6502) uint {
	m.PC++
	m.A = m.Y
	m.setNZ(m.A)
	return 2
}

type opFunc func(m *M6502) uint

var opcodeTable = [256]opFunc{
	opBRKimpl, opORAxind, opILL, opILL, opILL, opORAz, opASLz, opILL, opPHPimpl, opORAimm, opASLacc, opILL, opILL, opORAabs, opASLabs, opILL,
	opBPLrel, opORAindy, opILL, opILL, opILL, opORAzx, opASLzx, opILL, opCLCimpl, opORAabsy, opILL, opILL, opILL, opORAabsx, opASLabsx, opILL,
	opJSRabs, opANDxind, opILL, opILL, opBITz, opANDz, opROLz, opILL, opPLPimpl, opANDimm, opROLacc, opILL, opBITabs, opANDabs, opROLabs, opILL,
	opBMIrel, opANDindy, opILL, opILL, opILL, opANDzx, opROLzx, opILL, opSECimpl, opANDabsy, opILL, opILL, opILL, opANDabsx, opROLabsx, opILL,
	opRTIimpl, opEORxind, opILL, opILL, opILL, opEORz, opLSRz, opILL, opPHAimpl, opEORimm, opLSRacc, opILL, opJMPabs, opEORabs, opLSRabs, opILL,
	opBVCrel, opEORindy, opILL, opILL, opILL, opEORzx, opLSRzx, opILL, opCLIimpl, opEORabsy, opILL, opILL, opILL, opEORabsx, opLSRabsx, opILL,
	opRTSimpl, opADCxind, opILL, opILL, opILL, opADCz, opRORz, opILL, opPLAimpl, opADCimm, opRORacc, opILL, opJMPind, opADCabs, opRORabs, opILL,
	opBVSrel, opADCindy, opILL, opILL, opILL, opADCzx, opRORzx, opILL, opSEIimpl, opADCabsy, opILL, opILL, opILL, opADCabsx, opRORabsx, opILL,
	opILL, opSTAxind, opILL, opILL, opSTYz, opSTAz, opSTXz, opILL, opDEYimpl, opILL, opTXAimpl, opILL, opSTYabs, opSTAabs, opSTXabs, opILL,
	opBCCrel, opSTAindy, opILL, opILL, opSTYzx, opSTAzx, opSTXzy, opILL, opTYAimpl, opSTAabsy, opTXSimpl, opILL, opILL, opSTAabsx, opILL, opILL,
	opLDYimm, opLDAxind, opLDXimm, opILL, opLDYz, opLDAz, opLDXz, opILL, opTAYimpl, opLDAimm, opTAXimpl, opILL, opLDYabs, opLDAabs, opLDXabs, opILL,
	opBCSrel, opLDAindy, opILL, opILL, opLDYzx, opLDAzx, opLDXzy, opILL, opCLVimpl, opLDAabsy, opTSXimpl, opILL, opLDYabsx, opLDAabsx, opLDXabsy, opILL,
	opCPYimm, opCMPxind, opILL, opILL, opCPYz, opCMPz, opDECz, opILL, opINYimpl, opCMPimm, opDEXimpl, opILL, opCPYabs, opCMPabs, opDECabs, opILL,
	opBNErel, opCMPindy, opILL, opILL, opILL, opCMPzx, opDECzx, opILL, opCLDimpl, opCMPabsy, opILL, opILL, opILL, opCMPabsx, opDECabsx, opILL,
	opCPXimm, opSBCxind, opILL, opILL, opCPXz, opSBCz, opINCz, opILL, opINXimpl, opSBCimm, opNOPimpl, opILL, opCPXabs, opSBCabs, opINCabs, opILL,
	opBEQimm, opSBCindy, opILL, opILL, opILL, opSBCzx, opINCzx, opILL, opSEDimpl, opSBCabsy, opILL, opILL, opILL, opSBCabsx, opINCabsx, opILL,
}

//-----------------------------------------------------------------------------

// New6502 returns a 6502 CPU in the powered-on and reset state.
func New6502(mem Memory) *M6502 {
	var m M6502
	m.mem = mem
	m.Power(true)
	m.Reset()
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
		m.irq = false
		m.nmi = false
	} else {
		m.PC = 0
		m.S = 0
		m.P = 0
		m.A = 0
		m.X = 0
		m.Y = 0
		m.irq = false
		m.nmi = false
	}
}

// Reset the 6502 CPU.
func (m *M6502) Reset() {
	m.PC = m.readPointer(rstAddress)
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

// Run the 6502 CPU for a number of clock cycles.
func (m *M6502) Run(cycles uint) uint {

	var clks uint

	for clks < cycles {

		// nmi handling
		if m.nmi {
			m.nmi = false                    // clear the nmi
			m.P &= ^flagB                    // clear the break flag
			m.push16(m.PC)                   // save return addres in the stack.
			m.push8(m.P)                     // save current status in the stack.
			m.PC = m.readPointer(nmiAddress) // make PC point to the NMI routine.
			m.P |= flagI                     // disable interrupts
			clks += 7                        // accepting an NMI consumes 7 ticks.
			continue
		}

		// irq handling
		if m.irq && (m.P&flagI == 0) {
			m.P &= ^flagB
			m.push16(m.PC)
			m.push8(m.P)
			m.PC = m.readPointer(irqAddress)
			m.P |= flagI
			clks += 7
			continue
		}

		opcode := m.read8(m.PC)
		clks += opcodeTable[opcode](m)
	}

	return clks
}

//-----------------------------------------------------------------------------
