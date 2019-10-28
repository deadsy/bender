//-----------------------------------------------------------------------------
/*

6502 Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deadsy/bender/cpu"
	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

const historyPath = "history.txt"

//-----------------------------------------------------------------------------
// target memory

type memory struct {
	ram   [64 << 10]uint8
	spAdr uint8 // sim6502: zero page stack pointer address
}

// Read8 reads a byte from memory.
func (m *memory) Read8(adr uint16) uint8 {
	return m.ram[adr]
}

// Write8 writes a byte to memory.
func (m *memory) Write8(adr uint16, val uint8) {
	m.ram[adr] = val
}

func (m *memory) read16(adr uint16) uint16 {
	l := uint16(m.Read8(adr))
	h := uint16(m.Read8(adr + 1))
	return (h << 8) | l
}

func (m *memory) read16zp(adr uint8) uint16 {
	l := uint16(m.Read8(uint16(adr)))
	h := uint16(m.Read8(uint16(adr + 1)))
	return (h << 8) | l
}

func (m *memory) write16(adr uint16, val uint16) {
	m.Write8(adr, uint8(val))
	m.Write8(adr+1, uint8(val>>8))
}

func newMemory() *memory {
	m := memory{}
	// all 0xffs
	for i := range m.ram {
		m.ram[i] = 0xff
	}
	return &m
}

//-----------------------------------------------------------------------------

// userApp is state associated with the user application.
type userApp struct {
	mem *memory
	cpu *cpu.M6502
}

// newUserApp returns a user application.
func newUserApp() *userApp {
	mem := newMemory()
	cpu := cpu.New6502(mem)
	return &userApp{
		mem: mem,
		cpu: cpu,
	}
}

//-----------------------------------------------------------------------------
// file loading

// loadSim6502 loads a sim6502 binary file as produced by the cc65 tools.
func (u *userApp) loadSim6502(filename string, x []uint8) (string, error) {

	// validate the header
	if string(x[0:5]) != "sim65" {
		return "", fmt.Errorf("%s: bad magic", filename)
	}
	if x[5] != 2 {
		return "", fmt.Errorf("%s: bad version", filename)
	}
	if x[6] != 0 {
		return "", fmt.Errorf("%s: bad cpu type", filename)
	}

	// zero page stack pointer address (virtual subroutine abi)
	u.mem.spAdr = x[7]

	// copy the code to the load address
	loadAdr := uint16(x[8]) | (uint16(x[9]) << 8)
	for i, v := range x[12:] {
		u.mem.Write8(loadAdr+uint16(i), v)
	}
	endAdr := loadAdr + uint16(len(x[12:])) - 1

	// setup the reset address
	rstAdr := uint16(x[10]) | (uint16(x[11]) << 8)
	u.mem.write16(cpu.RstAddress, rstAdr)

	// Add the sim6502 VSRs
	u.cpu.AddVSR(0xfff4, vsrOpen)
	u.cpu.AddVSR(0xfff5, vsrClose)
	u.cpu.AddVSR(0xfff6, vsrRead)
	u.cpu.AddVSR(0xfff7, vsrWrite)
	u.cpu.AddVSR(0xfff8, vsrArgs)
	u.cpu.AddVSR(0xfff9, vsrExit)

	return fmt.Sprintf("%s code %04x-%04x reset %04x sp %02x", filename, loadAdr, endAdr, rstAdr, u.mem.spAdr), nil
}

// loadRaw loads a raw binary file.
func (u *userApp) loadRaw(filename string, x []uint8) (string, error) {

	// copy the code to the load address
	var loadAdr uint16
	for i, v := range x {
		u.mem.Write8(loadAdr+uint16(i), v)
	}
	endAdr := loadAdr + uint16(len(x)) - 1

	return fmt.Sprintf("%s code %04x-%04x", filename, loadAdr, endAdr), nil
}

func (u *userApp) loadFile(filename string) (string, error) {

	// get the file contents
	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// identify the file type
	if string(x[0:5]) == "sim65" {
		return u.loadSim6502(filename, x)
	}

	return u.loadRaw(filename, x)
}

//-----------------------------------------------------------------------------

// Put outputs a string to the user application.
func (u *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	fname := flag.String("f", "out.bin", "file to load (sim6502 or raw)")
	flag.Parse()

	// create the application
	app := newUserApp()

	// load the file
	status, err := app.loadFile(*fname)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", status)
	}

	// create the cli
	c := cli.NewCLI(app)
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt("emu> ")

	// reset the cpu
	app.cpu.Power(false)
	app.cpu.Power(true)
	app.cpu.Reset()

	// run the cli
	for c.Running() {
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	os.Exit(0)
}

//-----------------------------------------------------------------------------
