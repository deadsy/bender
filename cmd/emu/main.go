//-----------------------------------------------------------------------------
/*

6502 Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"crypto/rand"
	"fmt"
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

func newMemory() *memory {
	m := memory{}
	rand.Read(m.ram[:])
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
