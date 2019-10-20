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

// generate the opcode function name
func opcodeFuncName(code uint8) string {
	x := opcodeLookup(code)
	ins := strings.ToUpper(x.ins)
	return fmt.Sprintf("op%s", ins)
}

// generate the opcode function comment
func opcodeFuncComment(code uint8) string {
	x := opcodeLookup(code)
	name := opcodeFuncName(code)
	descr := insDescr[x.ins]
	return fmt.Sprintf("// %s, %s", name, descr)
}

// generate the opcode function template
func genOpcodeFunction(code uint8) string {
	name := opcodeFuncName(code)
	s := make([]string, 5)
	s[0] = opcodeFuncComment(code)
	s[1] = fmt.Sprintf("func %s(m *M6502, op uint8) uint {", name)
	s[2] = fmt.Sprintf("emuTODO()")
	s[3] = fmt.Sprintf("return %d", 0)
	s[4] = "}"
	return strings.Join(s, "\n")
}

// genOpcodes generates the unique sorted list of opcodes.
func genOpcodes() []uint8 {

	// get the unique set of opcode function names
	fset := make(map[string]uint8)
	for code := 0; code < 256; code++ {
		fset[opcodeFuncName(uint8(code))] = uint8(code)
	}

	// sort the opcode function names
	flist := make([]string, len(fset))
	i := 0
	for name := range fset {
		flist[i] = name
		i++
	}
	sort.Strings(flist)

	// return the sorted opcodes
	code := make([]uint8, len(flist))
	for i := range flist {
		code[i] = fset[flist[i]]
	}
	return code
}

func genOpcodeFunctions() string {
	opcodes := genOpcodes()
	s := make([]string, len(opcodes))
	for i, code := range opcodes {
		s[i] = genOpcodeFunction(code)
	}
	return strings.Join(s, "\n\n")
}

func genOpcodeTable() string {

	f := make([]string, 16)
	for i := 0; i < 16; i++ {
		l := make([]string, 16)
		for j := 0; j < 16; j++ {
			code := uint8((i * 16) + j)
			l[j] = opcodeFuncName(code)
		}
		f[i] = fmt.Sprintf("%s,", strings.Join(l, ","))
	}

	s := make([]string, 3)
	s[0] = "var opcodeTable = [256]opFunc{"
	s[1] = strings.Join(f, "\n")
	s[2] = "}"
	return strings.Join(s, "\n")
}

// GenOpcodeFunctions generates opcode table and template functions.
func GenOpcodeFunctions() string {
	s := make([]string, 3)
	s[0] = genOpcodeFunctions()
	s[1] = "type opFunc func(m *M6502, op uint8) uint"
	s[2] = genOpcodeTable()
	return strings.Join(s, "\n\n")

}

//-----------------------------------------------------------------------------
