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
	m.reg.P = (m.reg.P &^ flagNZ) | flags
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
	m.mem.Write8(stkAddress+uint16(m.reg.S), val)
	m.reg.S--
}

func (m *M6502) pop8() uint8 {
	m.reg.S++
	return m.mem.Read8(stkAddress + uint16(m.reg.S))
}

func (m *M6502) push16(val uint16) {
	m.mem.Write8(stkAddress+uint16(m.reg.S), uint8(val>>8))
	m.mem.Write8(stkAddress+uint16(m.reg.S-1), uint8(val))
	m.reg.S -= 2
}

func (m *M6502) pop16() uint16 {
	l := uint16(m.mem.Read8(stkAddress + uint16(m.reg.S+1)))
	h := uint16(m.mem.Read8(stkAddress + uint16(m.reg.S+2)))
	m.reg.S += 2
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

// opADC, add with carry
func opADC(m *M6502, op uint8) uint {

	v := uint8(0) // TODO
	n := uint(3)  // TODO

	c := m.reg.P & flagC
	if m.reg.P&flagD != 0 {
		l := uint(m.reg.A&0x0F) + uint(v&0x0F) + uint(c)
		h := uint(m.reg.A&0xF0) + uint(v&0xF0)
		m.reg.P &= ^(flagV | flagC | flagN | flagZ)
		if (l+h)&0xFF == 0 {
			m.reg.P |= flagZ
		}
		if l > 0x09 {
			h += 0x10
			l += 0x06
		}
		if h&0x80 != 0 {
			m.reg.P |= flagN
		}
		if ^(m.reg.A^v)&(m.reg.A^uint8(h))&0x80 != 0 {
			m.reg.P |= flagV
		}
		if h > 0x90 {
			h += 0x60
		}
		if h>>8 != 0 {
			m.reg.P |= flagC
		}
		m.reg.A = uint8(l&0x0F) | uint8(h&0xF0)
	} else {
		t := uint(m.reg.A) + uint(v) + uint(c)
		m.reg.P &= ^(flagV | flagC)
		if ^(m.reg.A^v)&(m.reg.A^uint8(t))&0x80 != 0 {
			m.reg.P |= flagV
		}
		if t>>8 != 0 {
			m.reg.P |= flagC
		}
		m.reg.A = uint8(t)
		m.setNZ(m.reg.A)
	}

	return n
}

// opAND, and (with accumulator)
func opAND(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opASL, arithmetic shift left
func opASL(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBCC, branch on carry clear
func opBCC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBCS, branch on carry set
func opBCS(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBEQ, branch on equal (zero set)
func opBEQ(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBIT, bit test
func opBIT(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBMI, branch on minus (negative set)
func opBMI(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBNE, branch on not equal (zero clear)
func opBNE(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBPL, branch on plus (negative clear)
func opBPL(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBRK, break / interrupt
func opBRK(m *M6502, op uint8) uint {
	m.read8(m.reg.PC + 1)
	m.push16(m.reg.PC + 2)
	m.push8(m.reg.P | flagB)
	m.reg.P |= flagB | flagI
	m.reg.PC = m.readPointer(brkAddress)
	return 7
}

// opBVC, branch on overflow clear
func opBVC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opBVS, branch on overflow set
func opBVS(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCLC, clear carry
func opCLC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCLD, clear decimal
func opCLD(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCLI, clear interrupt disable
func opCLI(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCLV, clear overflow
func opCLV(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCMP, compare (with accumulator)
func opCMP(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCPX, compare with X
func opCPX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opCPY, compare with Y
func opCPY(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opDEC, decrement
func opDEC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opDEX, decrement X
func opDEX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opDEY, decrement Y
func opDEY(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opEOR, exclusive or (with accumulator)
func opEOR(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opILL,
func opILL(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opINC, increment
func opINC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opINX, increment X
func opINX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opINY, increment Y
func opINY(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opJMPabs, jump, absolute mode
func opJMPabs(m *M6502, op uint8) uint {
	m.reg.PC = m.read16(m.reg.PC + 1)
	return 3
}

// opJMPind, jump, indirect mode
func opJMPind(m *M6502, op uint8) uint {
	m.reg.PC = m.read16(m.read16(m.reg.PC + 1))
	return 5
}

// opJSR, jump subroutine
func opJSR(m *M6502, op uint8) uint {
	m.push16(m.reg.PC + 2)
	m.reg.PC = m.read16(m.reg.PC + 1)
	return 6
}

// opLDA, load accumulator
func opLDA(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opLDX, load X
func opLDX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opLDY, load Y
func opLDY(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opLSR, logical shift right
func opLSR(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opNOP, no operation
func opNOP(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opORA, or with accumulator
func opORA(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opPHA, push accumulator
func opPHA(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opPHP, push processor status (SR)
func opPHP(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opPLA, pull accumulator
func opPLA(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opPLP, pull processor status (SR)
func opPLP(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opROL, rotate left
func opROL(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opROR, rotate right
func opROR(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opRTI, return from interrupt
func opRTI(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opRTS, return from subroutine
func opRTS(m *M6502, op uint8) uint {
	m.reg.PC = m.pop16() + 1
	return 6
}

// opSBC, subtract with carry
func opSBC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSEC, set carry
func opSEC(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSED, set decimal
func opSED(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSEI, set interrupt disable
func opSEI(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSTA, store accumulator
func opSTA(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSTX, store X
func opSTX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opSTY, store Y
func opSTY(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opTAX, transfer accumulator to X
func opTAX(m *M6502, op uint8) uint {
	m.reg.PC++
	m.reg.X = m.reg.A
	m.setNZ(m.reg.X)
	return 2
}

// opTAY, transfer accumulator to Y
func opTAY(m *M6502, op uint8) uint {
	m.reg.PC++
	m.reg.Y = m.reg.A
	m.setNZ(m.reg.Y)
	return 2
}

// opTSX, transfer stack pointer to X
func opTSX(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opTXA, transfer X to accumulator
func opTXA(m *M6502, op uint8) uint {
	m.reg.PC++
	m.reg.A = m.reg.X
	m.setNZ(m.reg.A)
	return 2
}

// opTXS, transfer X to stack pointer
func opTXS(m *M6502, op uint8) uint {
	emuTODO()
	return 0
}

// opTYA, transfer Y to accumulator
func opTYA(m *M6502, op uint8) uint {
	m.reg.PC++
	m.reg.A = m.reg.Y
	m.setNZ(m.reg.A)
	return 2
}

type opFunc func(m *M6502, op uint8) uint

var opcodeTable = [256]opFunc{
	opBRK, opORA, opILL, opILL, opILL, opORA, opASL, opILL, opPHP, opORA, opASL, opILL, opILL, opORA, opASL, opILL,
	opBPL, opORA, opILL, opILL, opILL, opORA, opASL, opILL, opCLC, opORA, opILL, opILL, opILL, opORA, opASL, opILL,
	opJSR, opAND, opILL, opILL, opBIT, opAND, opROL, opILL, opPLP, opAND, opROL, opILL, opBIT, opAND, opROL, opILL,
	opBMI, opAND, opILL, opILL, opILL, opAND, opROL, opILL, opSEC, opAND, opILL, opILL, opILL, opAND, opROL, opILL,
	opRTI, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL, opPHA, opEOR, opLSR, opILL, opJMPabs, opEOR, opLSR, opILL,
	opBVC, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL, opCLI, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL,
	opRTS, opADC, opILL, opILL, opILL, opADC, opROR, opILL, opPLA, opADC, opROR, opILL, opJMPind, opADC, opROR, opILL,
	opBVS, opADC, opILL, opILL, opILL, opADC, opROR, opILL, opSEI, opADC, opILL, opILL, opILL, opADC, opROR, opILL,
	opILL, opSTA, opILL, opILL, opSTY, opSTA, opSTX, opILL, opDEY, opILL, opTXA, opILL, opSTY, opSTA, opSTX, opILL,
	opBCC, opSTA, opILL, opILL, opSTY, opSTA, opSTX, opILL, opTYA, opSTA, opTXS, opILL, opILL, opSTA, opILL, opILL,
	opLDY, opLDA, opLDX, opILL, opLDY, opLDA, opLDX, opILL, opTAY, opLDA, opTAX, opILL, opLDY, opLDA, opLDX, opILL,
	opBCS, opLDA, opILL, opILL, opLDY, opLDA, opLDX, opILL, opCLV, opLDA, opTSX, opILL, opLDY, opLDA, opLDX, opILL,
	opCPY, opCMP, opILL, opILL, opCPY, opCMP, opDEC, opILL, opINY, opCMP, opDEX, opILL, opCPY, opCMP, opDEC, opILL,
	opBNE, opCMP, opILL, opILL, opILL, opCMP, opDEC, opILL, opCLD, opCMP, opILL, opILL, opILL, opCMP, opDEC, opILL,
	opCPX, opSBC, opILL, opILL, opCPX, opSBC, opINC, opILL, opINX, opSBC, opNOP, opILL, opCPX, opSBC, opINC, opILL,
	opBEQ, opSBC, opILL, opILL, opILL, opSBC, opINC, opILL, opSED, opSBC, opILL, opILL, opILL, opSBC, opINC, opILL,
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
		m.reg.PC = initialPC
		m.reg.S = initialS
		m.reg.P = initialP
		m.reg.A = initialA
		m.reg.X = initialX
		m.reg.Y = initialY
		m.irq = false
		m.nmi = false
	} else {
		m.reg.PC = 0
		m.reg.S = 0
		m.reg.P = 0
		m.reg.A = 0
		m.reg.X = 0
		m.reg.Y = 0
		m.irq = false
		m.nmi = false
	}
}

// Reset the 6502 CPU.
func (m *M6502) Reset() {
	m.reg.PC = m.readPointer(rstAddress)
	m.reg.S = initialS
	m.reg.P = initialP
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
			m.nmi = false                        // clear the nmi
			m.reg.P &= ^flagB                    // clear the break flag
			m.push16(m.reg.PC)                   // save return addres in the stack.
			m.push8(m.reg.P)                     // save current status in the stack.
			m.reg.PC = m.readPointer(nmiAddress) // make PC point to the NMI routine.
			m.reg.P |= flagI                     // disable interrupts
			clks += 7                            // accepting an NMI consumes 7 ticks.
			continue
		}

		// irq handling
		if m.irq && (m.reg.P&flagI == 0) {
			m.reg.P &= ^flagB
			m.push16(m.reg.PC)
			m.push8(m.reg.P)
			m.reg.PC = m.readPointer(irqAddress)
			m.reg.P |= flagI
			clks += 7
			continue
		}

		op := m.read8(m.reg.PC)
		clks += opcodeTable[op](m, op)
	}

	return clks
}

// ReadRegisters returns a copy of the 6502 CPU registers.
func (m *M6502) ReadRegisters() *Registers {
	x := m.reg
	return &x
}

// ReadPC returns the 6502 program counter.
func (m *M6502) ReadPC() uint16 {
	return m.reg.PC
}

//-----------------------------------------------------------------------------
