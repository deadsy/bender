//-----------------------------------------------------------------------------
/*

6502 Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"errors"
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

func (m *memory) Read8(adr uint16) uint8 {
	return m.ram[adr]
}

func (m *memory) Write8(adr uint16, val uint8) {
	m.ram[adr] = val
}

func (m *memory) Load(filename string) error {

	x, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// validate the header
	if string(x[0:5]) != "sim65" {
		return errors.New("bad magic")
	}
	if x[5] != 2 {
		return errors.New("bad version")
	}
	if x[6] != 0 {
		return errors.New("bad cpu")
	}

	// load the code
	load := uint16(x[8]) | (uint16(x[9]) << 8)
	for i, v := range x[12:] {
		m.Write8(load+uint16(i), v)
	}

	// setup the reset address
	m.Write8(cpu.RstAddress+0, x[10])
	m.Write8(cpu.RstAddress+1, x[11])

	return nil
}

func newMemory() *memory {
	m := memory{}
	for i := range m.ram {
		m.ram[i] = 0xff
	}
	return &m
}

//-----------------------------------------------------------------------------

// userApp is state associated with the user application.
type userApp struct {
	mem       *memory
	cpu       *cpu.M6502
	savedRegs *cpu.Registers
}

// newUserApp returns a user application.
func newUserApp() *userApp {
	mem := newMemory()
	err := mem.Load("./test1")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	cpu := cpu.New6502(mem)
	return &userApp{
		mem: mem,
		cpu: cpu,
	}
}

// Put outputs a string to the user application.
func (user *userApp) Put(s string) {
	fmt.Printf("%s", s)
}

//-----------------------------------------------------------------------------

func main() {
	c := cli.NewCLI(newUserApp())
	c.HistoryLoad(historyPath)
	c.SetRoot(menuRoot)
	c.SetPrompt("emu> ")
	for c.Running() {
		c.Run()
	}
	c.HistorySave(historyPath)
	os.Exit(0)

}

//-----------------------------------------------------------------------------
