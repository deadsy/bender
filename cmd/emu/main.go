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

func (u *userApp) loadSim6502(filename string) (string, error) {
	// get the file contents
	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

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

	u.cpu.Power(true)
	u.cpu.Reset()

	return fmt.Sprintf("%s code %04x-%04x reset %04x sp %02x", filename, loadAdr, endAdr, rstAdr, u.mem.spAdr), nil
}

func (u *userApp) loadRaw(filename string) (string, error) {

	// get the file contents
	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// copy the code to the load address
	var loadAdr uint16
	for i, v := range x {
		u.mem.Write8(loadAdr+uint16(i), v)
	}
	endAdr := loadAdr + uint16(len(x)) - 1

	u.cpu.Power(true)
	u.cpu.Reset()

	return fmt.Sprintf("%s code %04x-%04x", filename, loadAdr, endAdr), nil
}

// Put outputs a string to the user application.
func (u *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	sname := flag.String("s", "out.bin", "sim6502 binary file")
	rname := flag.String("r", "out.bin", "raw binary file")
	flag.Parse()

	_ = rname
	_ = sname

	// create the application
	app := newUserApp()

	//status, err := app.loadSim6502(*sname)
	status, err := app.loadRaw(*rname)

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

	// run the cli
	for c.Running() {
		c.Run()
	}

	// exit
	c.HistorySave(historyPath)
	os.Exit(0)
}

//-----------------------------------------------------------------------------
