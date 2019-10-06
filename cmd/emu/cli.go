//-----------------------------------------------------------------------------
/*

6502 Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------
// cli related leaf functions

var cmdHelp = cli.Leaf{
	Descr: "general help",
	F: func(c *cli.CLI, args []string) {
		c.GeneralHelp()
	},
}

var cmdHistory = cli.Leaf{
	Descr: "command history",
	F: func(c *cli.CLI, args []string) {
		c.SetLine(c.DisplayHistory(args))
	},
}

var cmdExit = cli.Leaf{
	Descr: "exit application",
	F: func(c *cli.CLI, args []string) {
		c.Exit()
	},
}

//-----------------------------------------------------------------------------
// memory functions

// memArgs converts memory arguments to an (address, size) tuple.
func memArgs(args []string) (uint16, uint16, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	adr := 0
	if len(args) >= 1 {
		adr, err = cli.IntArg(args[0], [2]int{0, 0xffff}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// size
	size := 0x40 // default size
	if len(args) >= 2 {
		size, err = cli.IntArg(args[1], [2]int{1, 0xffff}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return uint16(adr), uint16(size), nil
}

var cmdMemDisplay = cli.Leaf{
	Descr: "display memory",
	F: func(c *cli.CLI, args []string) {

		adr, size, err := memArgs(args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}

		fmt.Printf("%x %x\n", adr, size)

	},
}

var cmdMemToFile = cli.Leaf{
	Descr: "read from memory, write to file",
	F: func(c *cli.CLI, args []string) {
	},
}

var cmdMemFromFile = cli.Leaf{
	Descr: "read from file, write to memory",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdMemRead8 = cli.Leaf{
	Descr: "read 8 bits",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdMemRead16 = cli.Leaf{
	Descr: "read 16 bits",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdMemWrite8 = cli.Leaf{
	Descr: "write 8 bits",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdMemWrite16 = cli.Leaf{
	Descr: "write 16 bits",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdMemVerify = cli.Leaf{
	Descr: "verify memory against file",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var helpMemDisplay = []cli.Help{
	{"<adr> [len]", "address (hex)"},
	{"", "length (hex) - default is 0x40"},
}

var helpMemToFile = []cli.Help{
	{"<adr> <len> [file]", "address (hex)"},
	{"", "length (hex)"},
	{"", "filename - default is \"mem.bin\""},
}

var helpMemFromFile = []cli.Help{
	{"<adr> [file] [len]", "address (hex)"},
	{"", "filename - default is \"mem.bin\""},
	{"", "length (hex) - default is file length"},
}

var helpMemRead = []cli.Help{
	{"<adr>", "address (hex)"},
}

var helpMemWrite = []cli.Help{
	{"<adr> <val>", "address (hex)"},
	{"", "value (hex)"},
}

// memory submenu items
var memoryMenu = cli.Menu{
	{"display", cmdMemDisplay, helpMemDisplay},
	{">file", cmdMemToFile, helpMemToFile},
	{"<file", cmdMemFromFile, helpMemFromFile},
	{"r8", cmdMemRead8, helpMemRead},
	{"r16", cmdMemRead16, helpMemRead},
	{"verify", cmdMemVerify, helpMemToFile},
	{"w8", cmdMemWrite8, helpMemWrite},
	{"w16", cmdMemWrite16, helpMemWrite},
}

//-----------------------------------------------------------------------------

var cmdDisassemble = cli.Leaf{
	Descr: "disassemble memory",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdRegisters = cli.Leaf{
	Descr: "display cpu registers",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdRun = cli.Leaf{
	Descr: "run the emulation",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

var cmdStep = cli.Leaf{
	Descr: "single step the emulation",
	F: func(c *cli.CLI, args []string) {
		//c.Exit()
	},
}

//-----------------------------------------------------------------------------

var helpDisassemble = []cli.Help{
	{"[adr] [len]", "address (hex) - default is current pc"},
	{"", "length (hex) - default is 0x10"},
}

// root menu
var menuRoot = cli.Menu{
	{"da", cmdDisassemble, helpDisassemble},
	{"exit", cmdExit},
	{"help", cmdHelp},
	{"history", cmdHistory, cli.HistoryHelp},
	{"memory", memoryMenu, "memory functions"},
	{"regs", cmdRegisters},
	{"run", cmdRun},
	{"step", cmdStep},
}

//-----------------------------------------------------------------------------
