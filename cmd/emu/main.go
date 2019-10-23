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
	ram [64 << 10]uint8
}

// Read8 reads a byte from memory.
func (m *memory) Read8(adr uint16) uint8 {
	return m.ram[adr]
}

// Write8 writes a byte to memory.
func (m *memory) Write8(adr uint16, val uint8) {
	m.ram[adr] = val
}

func newMemory() *memory {
	m := memory{}
	// all 0xffs
	for i := range m.ram {
		m.ram[i] = 0xff
	}
	return &m
}

// Load loads a sim6502 style executable file (See cc65).
func (m *memory) Load(filename string) (string, error) {

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

	// copy the code to the load address
	loadAdr := uint16(x[8]) | (uint16(x[9]) << 8)
	for i, v := range x[12:] {
		m.Write8(loadAdr+uint16(i), v)
	}
	endAdr := loadAdr + uint16(len(x[12:]))

	// setup the reset address
	rstAdr := uint16(x[10]) | (uint16(x[11]) << 8)
	m.Write8(cpu.RstAddress+0, x[10])
	m.Write8(cpu.RstAddress+1, x[11])

	return fmt.Sprintf("%s code %04x-%04x reset %04x", filename, loadAdr, endAdr, rstAdr), nil
}

//-----------------------------------------------------------------------------

// userApp is state associated with the user application.
type userApp struct {
	mem       *memory
	cpu       *cpu.M6502
}

// newUserApp returns a user application.
func newUserApp(fname string) (*userApp, error) {
	mem := newMemory()
	status, err := mem.Load(fname)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", status)
	cpu := cpu.New6502(mem)
	return &userApp{
		mem: mem,
		cpu: cpu,
	}, nil
}

// Put outputs a string to the user application.
func (user *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	// command line flags
	fname := flag.String("f", "out.bin", "executable file")
	flag.Parse()

	// create the application
	app, err := newUserApp(*fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
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
