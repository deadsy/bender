//-----------------------------------------------------------------------------
/*

6502 Emulator

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"

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
func memArgs(args []string) (uint16, uint, error) {
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
		size, err = cli.IntArg(args[1], [2]int{1, 0x10000}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return uint16(adr), uint(size), nil
}

var cmdMemDisplay = cli.Leaf{
	Descr: "display memory",
	F: func(c *cli.CLI, args []string) {
		adr, size, err := memArgs(args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		// round down address to 16 byte boundary
		adr &= ^uint16(15)
		// round up n to an integral multiple of 16 bytes
		size = (size + 15) & ^uint(15)
		// print the header
		c.User.Put("addr  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F\n")
		// read and print the data
		for i := 0; i < int(size>>4); i++ {
			// read 16 bytes per line
			var data [16]string
			var ascii [16]string
			for j := 0; j < 16; j++ {
				x := c.User.(*userApp).mem.Read8(adr + uint16(j))
				data[j] = fmt.Sprintf("%02x", x)
				if x >= 32 && x <= 126 {
					ascii[j] = fmt.Sprintf("%c", x)
				} else {
					ascii[j] = "."
				}
			}
			dataStr := strings.Join(data[:], " ")
			asciiStr := strings.Join(ascii[:], "")
			c.User.Put(fmt.Sprintf("%04x  %s  %s\n", adr, dataStr, asciiStr))
			adr += 16
		}
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

// daArgs converts disassembly arguments to an (address, size) tuple.
func daArgs(pc uint16, args []string) (uint16, uint, error) {
	err := cli.CheckArgc(args, []int{0, 1, 2})
	if err != nil {
		return 0, 0, err
	}
	// address
	adr := int(pc) // default address
	if len(args) >= 1 {
		adr, err = cli.IntArg(args[0], [2]int{0, 0xffff}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	// size
	size := 16 // default size
	if len(args) >= 2 {
		size, err = cli.IntArg(args[1], [2]int{1, 2048}, 16)
		if err != nil {
			return 0, 0, err
		}
	}
	return uint16(adr), uint(size), nil
}

var cmdDisassemble = cli.Leaf{
	Descr: "disassemble memory",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		adr, size, err := daArgs(m.ReadPC(), args)
		if err != nil {
			c.User.Put(fmt.Sprintf("%s\n", err))
			return
		}
		c.User.Put(fmt.Sprintf("%s\n", m.Disassemble(adr, int(size))))
	},
}

var cmdRegisters = cli.Leaf{
	Descr: "display cpu registers",
	F: func(c *cli.CLI, args []string) {
		m := c.User.(*userApp).cpu
		c.User.Put(fmt.Sprintf("%s\n", m.Dump()))
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
