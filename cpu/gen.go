//-----------------------------------------------------------------------------
/*

6502 CPU Code Generation

*/
//-----------------------------------------------------------------------------

package cpu

import (
	"fmt"
	"sort"
	"strings"
)

//-----------------------------------------------------------------------------

func opcodeFuncName(code uint8) string {
	x := opcodeLookup(code)
	return fmt.Sprintf("op%s%s", x.ins, modeDescr[x.mode].suffix)
}

func genOpcodeFunc() string {

	// get the set of unique opcode function names
	fs := make(map[string]bool)
	for code := 0; code < 256; code++ {
		fs[opcodeFuncName(uint8(code))] = true
	}
	f := make([]string, len(fs))
	i := 0
	for name := range fs {
		f[i] = name
		i++
	}
	sort.Strings(f)
	return strings.Join(f, "\n")
}

func genOpcodeTable() string {

	f := make([]string, 16)
	for i := 0; i < 16; i++ {
		l := make([]string, 16)
		for j := 0; j < 16; j++ {
			code := uint8((i * 16) + j)
			l[j] = opcodeFuncName(code)
		}
		f[i] = strings.Join(l, ",")
	}

	s := make([]string, 3)
	s[0] = "var opcodeTable = [256]opFunc{"
	s[1] = strings.Join(f, "\n")
	s[2] = "}"
	return strings.Join(s, "\n")
}

// GenOpcodeFunctions generates opcode table and template functions.
func GenOpcodeFunctions() string {

	s := make([]string, 2)

	s[0] = genOpcodeFunc()
	s[1] = genOpcodeTable()

	return strings.Join(s, "\n")

}

//-----------------------------------------------------------------------------
