//-----------------------------------------------------------------------------
/*

6502 CPU Emulator

*/
//-----------------------------------------------------------------------------

package cpu

//-----------------------------------------------------------------------------

func (m *M6502) setNZ(val uint8) {
	var flags uint8
	if val != 0 {
		flags = val & flagN
	} else {
		flags = flagZ
	}
	m.p = (m.p &^ flagNZ) | flags
}

//-----------------------------------------------------------------------------

func (m *M6502) read8(adr uint16) uint8 {
	return 0
}

func (m *M6502) read16(adr uint16) uint16 {
	return 0
}

func (m *M6502) readPointer(adr uint16) uint16 {
	return 0
}

//-----------------------------------------------------------------------------

func (m *M6502) push8(val uint8) {
}

func (m *M6502) push16(val uint16) {
}

func (m *M6502) pop16() uint16 {
	return 0
}

//-----------------------------------------------------------------------------

// jmp WORD
func opJMPabs(m *M6502) uint {
	m.pc = m.read16(m.pc + 1)
	return 3
}

// jmp (WORD)
func opJMPind(m *M6502) uint {
	m.pc = m.read16(m.read16(m.pc + 1))
	return 5
}

// jsr WORD
func opJSR(m *M6502) uint {
	m.push16(m.pc + 2)
	m.pc = m.read16(m.pc + 1)
	return 6
}

// rts
func opRTS(m *M6502) uint {
	m.pc = m.pop16() + 1
	return 6
}

func opADC(m *M6502) uint {
	return 0
}
func opAND(m *M6502) uint {
	return 0
}
func opASL(m *M6502) uint {
	return 0
}
func opBCC(m *M6502) uint {
	return 0
}
func opBCS(m *M6502) uint {
	return 0
}
func opBEQ(m *M6502) uint {
	return 0
}
func opBIT(m *M6502) uint {
	return 0
}
func opBMI(m *M6502) uint {
	return 0
}
func opBNE(m *M6502) uint {
	return 0
}
func opBPL(m *M6502) uint {
	return 0
}

func opBVC(m *M6502) uint {
	return 0
}
func opBVS(m *M6502) uint {
	return 0
}
func opCLC(m *M6502) uint {
	return 0
}
func opCLD(m *M6502) uint {
	return 0
}
func opCLI(m *M6502) uint {
	return 0
}
func opCLV(m *M6502) uint {
	return 0
}
func opCMP(m *M6502) uint {
	return 0
}
func opCPX(m *M6502) uint {
	return 0
}
func opCPY(m *M6502) uint {
	return 0
}
func opDEC(m *M6502) uint {
	return 0
}
func opDEX(m *M6502) uint {
	return 0
}
func opDEY(m *M6502) uint {
	return 0
}
func opEOR(m *M6502) uint {
	return 0
}
func opINC(m *M6502) uint {
	return 0
}
func opINX(m *M6502) uint {
	return 0
}
func opINY(m *M6502) uint {
	return 0
}

func opLDA(m *M6502) uint {
	return 0
}
func opLDX(m *M6502) uint {
	return 0
}
func opLDY(m *M6502) uint {
	return 0
}
func opLSR(m *M6502) uint {
	return 0
}
func opNOP(m *M6502) uint {
	return 0
}
func opORA(m *M6502) uint {
	return 0
}
func opPHA(m *M6502) uint {
	return 0
}
func opPHP(m *M6502) uint {
	return 0
}
func opPLA(m *M6502) uint {
	return 0
}
func opPLP(m *M6502) uint {
	return 0
}
func opROL(m *M6502) uint {
	return 0
}
func opROR(m *M6502) uint {
	return 0
}
func opRTI(m *M6502) uint {
	return 0
}

func opSBC(m *M6502) uint {
	return 0
}
func opSEC(m *M6502) uint {
	return 0
}
func opSED(m *M6502) uint {
	return 0
}
func opSEI(m *M6502) uint {
	return 0
}
func opSTA(m *M6502) uint {
	return 0
}
func opSTX(m *M6502) uint {
	return 0
}
func opSTY(m *M6502) uint {
	return 0
}

func opTSX(m *M6502) uint {
	return 0
}

func opTXS(m *M6502) uint {
	return 0
}

func opTAX(m *M6502) uint {
	m.pc++
	m.x = m.a
	m.setNZ(m.x)
	return 2
}

func opTAY(m *M6502) uint {
	m.pc++
	m.y = m.a
	m.setNZ(m.y)
	return 2
}

func opTXA(m *M6502) uint {
	m.pc++
	m.a = m.x
	m.setNZ(m.a)
	return 2
}

func opTYA(m *M6502) uint {
	m.pc++
	m.a = m.y
	m.setNZ(m.a)
	return 2
}

func opBRK(m *M6502) uint {
	m.read8(m.pc + 1)
	m.push16(m.pc + 2)
	m.push8(m.p | flagB)
	m.p |= flagB | flagI
	m.pc = m.readPointer(brkAddress)
	return 7
}

func opILL(m *M6502) uint {
	return 2
}

type opFunc func(m *M6502) uint

var opcodeFunc = [256]opFunc{
	// 00  01     02     03     04     05     06     07     08     09     0a     0b     0c     0d     0e     0f
	opBRK, opORA, opILL, opILL, opILL, opORA, opASL, opILL, opPHP, opORA, opASL, opILL, opILL, opORA, opASL, opILL,
	// 10  11     12     13     14     15     16     17     18     19     1a     1b     1c     1d     1e     1f
	opBPL, opORA, opILL, opILL, opILL, opORA, opASL, opILL, opCLC, opORA, opILL, opILL, opILL, opORA, opASL, opILL,
	// 20  21     22     23     24     25     26     27     28     29     2a     2b     2c     2d     2e     2f
	opJSR, opAND, opILL, opILL, opBIT, opAND, opROL, opILL, opPLP, opAND, opROL, opILL, opBIT, opAND, opROL, opILL,
	// 30  31     32     33     34     35     36     37     38     39     3a     3b     3c     3d     3e     3f
	opBMI, opAND, opILL, opILL, opILL, opAND, opROL, opILL, opSEC, opAND, opILL, opILL, opILL, opAND, opROL, opILL,
	// 40  41     42     43     44     45     46     47     48     49     4a     4b     4c     4d     4e     4f
	opRTI, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL, opPHA, opEOR, opLSR, opILL, opJMPabs, opEOR, opLSR, opILL,
	// 50  51     52     53     54     55     56     57     58     59     5a     5b     5c     5d     5e     5f
	opBVC, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL, opCLI, opEOR, opILL, opILL, opILL, opEOR, opLSR, opILL,
	// 60  61     62     63     64     65     66     67     68     69     6a     6b     6c     6d     6e     6f
	opRTS, opADC, opILL, opILL, opILL, opADC, opROR, opILL, opPLA, opADC, opROR, opILL, opJMPind, opADC, opROR, opILL,
	// 70  71     72     73     74     75     76     77     78     79     7a     7b     7c     7d     7e     7f
	opBVS, opADC, opILL, opILL, opILL, opADC, opROR, opILL, opSEI, opADC, opILL, opILL, opILL, opADC, opROR, opILL,
	// 80  81     82     83     84     85     86     87     88     89     8a     8b     8c     8d     8e     8f
	opILL, opSTA, opILL, opILL, opSTY, opSTA, opSTX, opILL, opDEY, opILL, opTXA, opILL, opSTY, opSTA, opSTX, opILL,
	// 90  91     92     93     94     95     96     97     98     99     9a     9b     9c     9d     9e     9f
	opBCC, opSTA, opILL, opILL, opSTY, opSTA, opSTX, opILL, opTYA, opSTA, opTXS, opILL, opILL, opSTA, opILL, opILL,
	// a0  a1     a2     a3     a4     a5     a6     a7     a8     a9     aa     ab     ac     ad     ae     af
	opLDY, opLDA, opLDX, opILL, opLDY, opLDA, opLDX, opILL, opTAY, opLDA, opTAX, opILL, opLDY, opLDA, opLDX, opILL,
	// b0  b1     b2     b3     b4     b5     b6     b7     b8     b9     ba     bb     bc     bd     be     bf
	opBCS, opLDA, opILL, opILL, opLDY, opLDA, opLDX, opILL, opCLV, opLDA, opTSX, opILL, opLDY, opLDA, opLDX, opILL,
	// c0  c1     c2     c3     c4     c5     c6     c7     c8     c9     ca     cb     cc     cd     ce     cf
	opCPY, opCMP, opILL, opILL, opCPY, opCMP, opDEC, opILL, opINY, opCMP, opDEX, opILL, opCPY, opCMP, opDEC, opILL,
	// d0  d1     d2     d3     d4     d5     d6     d7     d8     d9     da     db     dc     dd     de     df
	opBNE, opCMP, opILL, opILL, opILL, opCMP, opDEC, opILL, opCLD, opCMP, opILL, opILL, opILL, opCMP, opDEC, opILL,
	// e0  e1     e2     e3     e4     e5     e6     e7     e8     e9     ea     eb     ec     ed     ee     ef
	opCPX, opSBC, opILL, opILL, opCPX, opSBC, opINC, opILL, opINX, opSBC, opNOP, opILL, opCPX, opSBC, opINC, opILL,
	// f0  f1     f2     f3     f4     f5     f6     f7     f8     f9     fa     fb     fc     fd     fe     ff
	opBEQ, opSBC, opILL, opILL, opILL, opSBC, opINC, opILL, opSED, opSBC, opILL, opILL, opILL, opSBC, opINC, opILL,
}

// done 0,1,2,3,4,5,6,7

//-----------------------------------------------------------------------------

// New6502 returns a 6502 CPU in the powered-on and reset state.
func New6502() *M6502 {
	var m M6502
	m.Power(true)
	m.Reset()
	return &m
}

// Power on/off the 6502 CPU.
func (m *M6502) Power(state bool) {
	if state {
		m.pc = initialPC
		m.s = initialS
		m.p = initialP
		m.a = initialA
		m.x = initialX
		m.y = initialY
		m.irq = false
		m.nmi = false
	} else {
		m.pc = 0
		m.s = 0
		m.p = 0
		m.a = 0
		m.x = 0
		m.y = 0
		m.irq = false
		m.nmi = false
	}
}

// Reset the 6502 CPU.
func (m *M6502) Reset() {
	m.pc = m.readPointer(resetAddress)
	m.s = initialS
	m.p = initialP
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

		// Execute instruction and update consumed cycles
		opcode := m.read8(m.pc)
		clks += opcodeFunc[opcode](m)
	}

	return clks
}

//-----------------------------------------------------------------------------
