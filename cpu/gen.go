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
	if x.ins == "ill" {
		return "opXX"
	}
	return fmt.Sprintf("op%02X", code)
}

// generate the opcode function comment
func opcodeFuncComment(code uint8) string {
	x := opcodeLookup(code)
	s := make([]string, 0)
	s = append(s, opcodeFuncName(code))
	s = append(s, fmt.Sprintf("%s %s", strings.ToUpper(x.ins), insDescr[x.ins]))
	if x.mode != amNone && x.mode != amImpl {
		s = append(s, modeDescr[x.mode].descr)
	}
	return fmt.Sprintf("// %s", strings.Join(s, ", "))
}

// generate the opcode function template
func genOpcodeFunction(code uint8) string {

  s := make([]string, 0)
	s = append(s, opcodeFuncComment(code))
	s = append(s, fmt.Sprintf("func %s(m *M6502) uint {", opcodeFuncName(code)))
	s = append(s, fmt.Sprintf("panic(\"TODO\")"))

	n := insLength(code)
	if n == 1 {
		s = append(s, "m.PC ++")
	} else {
		s = append(s, fmt.Sprintf("m.PC += %d", n))
	}

	s = append(s, "return 0")
	s = append(s, "}")
	return strings.Join(s, "\n")
}

// genOpcodes generates the unique sorted list of opcodes.
func genOpcodes() []uint8 {

	// get the unique set of opcode mneumonics
	fset := make(map[string]uint8)
	for code := 0; code < 256; code++ {
		x := opcodeLookup(uint8(code))
		var s string
		if x.ins == "ill" {
			s = "ill"
		} else {
			s = fmt.Sprintf("%s%02x", x.ins, code)
		}
		fset[s] = uint8(code)
	}

	// sort the opcodes
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
	s[1] = "type opFunc func(m *M6502) uint"
	s[2] = genOpcodeTable()
	return strings.Join(s, "\n\n")

}

//-----------------------------------------------------------------------------
