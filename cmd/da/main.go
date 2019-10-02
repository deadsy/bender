//-----------------------------------------------------------------------------
/*

6502 Disassembler

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"

	"github.com/deadsy/bender/cpu"
)

//-----------------------------------------------------------------------------

func memString(adr uint16, mem []byte) string {
	s := make([]string, len(mem))
	for i, v := range mem {
		s[i] = fmt.Sprintf("%02x", v)
	}
	return fmt.Sprintf("%04x: %-15s", uint16(adr), strings.Join(s, " "))
}

//-----------------------------------------------------------------------------

func main() {

	symbols := map[uint16]string{
		0x0000: "fmpnt",
		0x0002: "topnt",
		0x0004: "cntr",
		0x0005: "tsign",
		0x0006: "signs",
		0x0007: "fplswe",
		//0x0008: "fpacce",
		//0x0008: "fplsw",
		0x0009: "fpnsw",
		0x000a: "fpmsw",
		0x000c: "mcando",
		0x000d: "mcand1",
		0x000e: "mcand2",
		0x000f: "folswe",
		0x0010: "foplsw",
		0x0011: "fopnsw",
		0x0012: "fopmsw",
		0x0013: "fopexp",
		0x0014: "work0",
		0x0015: "work1",
		0x0016: "work2",
		0x0017: "work3",
		0x0018: "work4",
		0x0019: "works",
		0x001a: "work6",
		0x001b: "work7",
		0x001c: "inmtas",
		0x001d: "inexps",
		0x001e: "inprdi",
		0x001f: "iolsw",
		0x0020: "ionsw",
		0x0021: "iomsw",
		0x0022: "ioexp",
		0x0023: "iostr",
		0x0024: "iostr1",
		0x0025: "iostr2",
		0x0026: "iostr3",
		0x0027: "ioexpd",
		0x0028: "tplsw",
		0x0029: "tpnsw",
		0x002a: "tpmsw",
		0x002b: "tpexp",
		0x002c: "temp1",
		0x0200: "clrmem",
		0x0203: "clrm1",
		//0x0203: "compl",
		0x020a: "rotatl",
		0x020b: "rotl",
		0x0211: "morrtl",
		0x0215: "rotatr",
		0x0216: "rotr",
		0x021c: "morrtr",
		0x0220: "complm",
		0x022e: "adder",
		0x0231: "addr1",
		0x023c: "subber",
		0x023f: "subb1",
		0x024a: "movind",
		0x024c: "movin1",
		0x0255: "fpnorm",
		0x0262: "accmin",
		0x026b: "aczert",
		0x026f: "looko",
		0x0279: "normex",
		0x027a: "acnonz",
		0x028a: "accset",
		0x029a: "fdadd",
		0x029e: "movop",
		0x02b1: "nonzac",
		0x02b6: "ckeqex",
		0x02ce: "skpneg",
		0x02da: "lineup",
		0x02e2: "moracc",
		0x02ed: "shifto",
		0x02f5: "shacop",
		0x0303: "negop",
		0x0315: "shloop",
		0x031b: "fshift",
		0x0326: "bring1",
		0x032a: "rescnt",
		0x032d: "fpsub",
		0x0337: "fpmult",
		0x033a: "addexp",
		0x0343: "setmct",
		0x0347: "multip",
		0x0350: "adoppp",
		//0x0350: "nadopp",
		0x037f: "cround",
		0x0389: "prexfr",
		0x0393: "exmldv",
		0x03a4: "multex",
		0x03a5: "cksign",
		0x03c5: "negfpa",
		0x03ce: "opsgnt",
		0x03dc: "fpdiv",
		0x03e3: "subexp",
		0x03ec: "setdct",
		0x03f0: "divide",
		0x0401: "execho",
		0x0404: "expinp",
		0x0406: "derror",
		0x0407: "noexps",
		0x040b: "nogo",
		0x040c: "quorot",
		0x0441: "dvexit",
		0x044e: "setsub",
		0x0468: "subr1",
		0x0475: "fpinp",
		0x0492: "secho",
		0x0495: "ninput",
		0x0498: "notplm",
		0x049c: "erase",
		0x04a7: "serase",
		0x04ab: "period",
		0x04b1: "per1",
		0x04bd: "spriod",
		0x04c1: "fndexp",
		0x04dd: "island",
		0x0503: "sfndxp",
		0x0535: "endinp",
		0x0540: "finput",
		0x0564: "posexp",
		0x056d: "expok",
		0x0577: "expfix",
		0x057d: "fpx10",
		0x0591: "minexp",
		0x0597: "fpd10",
		0x05ab: "decbin",
		0x05de: "fpout",
		0x05ea: "outneg",
		0x05f3: "ahead1",
		0x0602: "decext",
		0x060e: "decrep",
		0x0613: "decexd",
		0x0619: "decout",
		0x0634: "compen",
		0x0642: "outdig",
		0x064a: "outdgs",
		0x0651: "decrdg",
		0x065b: "zerodg",
		0x066f: "expout",
		0x067d: "exoutn",
		0x0685: "ahead2",
		0x068c: "sub12",
		0x0697: "tomuch",
		0x06a4: "fpcont",
		0x06c1: "nvalid",
		0x06d1: "notadd",
		0x06de: "notsub",
		0x06eb: "notmul",
		0x06f5: "final",
		0x06fb: "notdiv",
		0x0701: "operat",
		0x0780: "input",
		0x07c0: "spaces",
		0x07c5: "echo",
	}

	_ = symbols

	code := []byte{
		0xa9, 0x00, 0xa8, 0x91, 0x02, 0xc8, 0xca, 0xd0,
		0xfa, 0x60, 0x18, 0x36, 0x00, 0x88, 0xd0, 0x01,
		0x60, 0xe8, 0x4c, 0x0b, 0x02, 0x18, 0x76, 0x00,
		0x88, 0xd0, 0x01, 0x60, 0xca, 0x4c, 0x16, 0x02,
		0x38, 0xa9, 0xff, 0x55, 0x00, 0x69, 0x00, 0x95,
		0x00, 0xe8, 0x88, 0xd0, 0xf4, 0x60, 0x18, 0xa0,
		0x00, 0xb1, 0x02, 0x71, 0x00, 0x91, 0x02, 0xc8,
		0xca, 0xd0, 0xf6, 0x60, 0x38, 0xa0, 0x00, 0xb1,
		0x02, 0xf1, 0x00, 0x91, 0x02, 0xc8, 0xca, 0xd0,
		0xf6, 0x60, 0xa0, 0x00, 0xb1, 0x00, 0x91, 0x02,
		0xc8, 0xca, 0xd0, 0xf8, 0x60, 0xa2, 0x05, 0xa5,
		0x0a, 0x30, 0x07, 0xa0, 0x00, 0x94, 0x00, 0x4c,
		0x6b, 0x02, 0x95, 0x00, 0xa0, 0x04, 0xa2, 0x07,
		0x20, 0x20, 0x02, 0xa2, 0x0a, 0xa0, 0x04, 0xb5,
		0x00, 0xd0, 0x07, 0xca, 0x88, 0xd0, 0xf8, 0x84,
		0x0b, 0x60, 0xa2, 0x07, 0xa0, 0x04, 0x20, 0x0a,
		0x02, 0xb5, 0x00, 0x30, 0x05, 0xc6, 0x0b, 0x4c,
		0x7a, 0x02, 0xa2, 0x0a, 0xa0, 0x03, 0x20, 0x15,
		0x02, 0xa5, 0x05, 0xf0, 0xe4, 0xa0, 0x03, 0x4c,
		0x20, 0x02, 0xa5, 0x0a, 0xd0, 0x13, 0xa2, 0x10,
		0x86, 0x00, 0xa2, 0x08, 0x86, 0x02, 0xa9, 0x00,
		0x85, 0x01, 0x85, 0x03, 0xa2, 0x04, 0x4c, 0x4a,
		0x02, 0xa5, 0x12, 0xd0, 0x01, 0x60, 0xa2, 0x0b,
		0xb5, 0x00, 0xc5, 0x13, 0xf0, 0x37, 0x38, 0xa9,
		0x00, 0xf5, 0x00, 0x65, 0x13, 0x10, 0x07, 0x38,
		0x85, 0x2c, 0xa9, 0x00, 0xe5, 0x2c, 0xc9, 0x18,
		0x30, 0x08, 0x38, 0xa5, 0x13, 0xf5, 0x00, 0x10,
		0xc5, 0x60, 0xa5, 0x13, 0x38, 0xf5, 0x00, 0xa8,
		0x30, 0x0b, 0xa2, 0x0b, 0x20, 0x15, 0x03, 0x88,
		0xd0, 0xf8, 0x4c, 0xf5, 0x02, 0xa2, 0x13, 0x20,
		0x15, 0x03, 0xc8, 0xd0, 0xf8, 0xa9, 0x00, 0x85,
		0x07, 0x85, 0x0f, 0xa2, 0x0b, 0x20, 0x15, 0x03,
		0xa2, 0x13, 0x20, 0x15, 0x03, 0xa2, 0x0f, 0x86,
		0x00, 0xa2, 0x07, 0x86, 0x02, 0xa2, 0x04, 0x20,
		0x2e, 0x02, 0x4c, 0x55, 0x02, 0xf6, 0x00, 0xca,
		0x98, 0xa0, 0x04, 0x48, 0xb5, 0x00, 0x30, 0x06,
		0x20, 0x15, 0x02, 0x4c, 0x2a, 0x03, 0x38, 0x20,
		0x16, 0x02, 0x68, 0xa8, 0x60, 0xa2, 0x08, 0xa0,
		0x03, 0x20, 0x20, 0x02, 0x4c, 0x9a, 0x02, 0x20,
		0xa5, 0x03, 0xa5, 0x13, 0x18, 0x65, 0x0b, 0x85,
		0x0b, 0xe6, 0x0b, 0xa9, 0x17, 0x85, 0x04, 0xa2,
		0x0a, 0xa0, 0x03, 0x20, 0x15, 0x02, 0x90, 0x0d,
		0xa2, 0x0d, 0x86, 0x00, 0xa2, 0x15, 0x86, 0x02,
		0xa2, 0x06, 0x20, 0x2e, 0x02, 0xa2, 0x1a, 0xa0,
		0x06, 0x20, 0x15, 0x02, 0xc6, 0x04, 0xd0, 0xdf,
		0xa2, 0x1a, 0xa0, 0x06, 0x20, 0x15, 0x02, 0xa6,
		0x17, 0xb5, 0x00, 0x2a, 0x10, 0x13, 0x18, 0xa0,
		0x03, 0xa9, 0x40, 0x75, 0x00, 0x85, 0x17, 0xa9,
		0x00, 0x75, 0x00, 0x95, 0x00, 0xe8, 0x88, 0xd0,
		0xf6, 0xa2, 0x07, 0x86, 0x02, 0xa2, 0x17, 0x86,
		0x00, 0xa2, 0x04, 0x20, 0x4a, 0x02, 0x20, 0x55,
		0x02, 0xa5, 0x06, 0xd0, 0x07, 0xa2, 0x08, 0xa0,
		0x03, 0x20, 0x20, 0x02, 0x60, 0xa9, 0x00, 0x85,
		0x03, 0x85, 0x01, 0xa9, 0x14, 0x85, 0x02, 0xa2,
		0x08, 0x20, 0x00, 0x02, 0xa9, 0x0c, 0x85, 0x02,
		0xa2, 0x04, 0x20, 0x00, 0x02, 0xa9, 0x01, 0x85,
		0x06, 0xa5, 0x0a, 0x10, 0x09, 0xc6, 0x06, 0xa2,
		0x08, 0xa0, 0x03, 0x20, 0x20, 0x02, 0xa5, 0x12,
		0x30, 0x01, 0x60, 0xc6, 0x06, 0xa2, 0x10, 0xa0,
		0x03, 0x4c, 0x20, 0x02, 0x20, 0xa5, 0x03, 0xa5,
		0x0a, 0xf0, 0x23, 0xa5, 0x13, 0x38, 0xe5, 0x0b,
		0x85, 0x0b, 0xe6, 0x0b, 0xa9, 0x17, 0x85, 0x04,
		0x20, 0x4e, 0x04, 0x30, 0x16, 0xa2, 0x10, 0x86,
		0x02, 0xa2, 0x14, 0x86, 0x00, 0xa2, 0x03, 0x20,
		0x4a, 0x02, 0x38, 0x4c, 0x0c, 0x04, 0xa9, 0xbf,
		0x4c, 0x06, 0x04, 0x18, 0xa2, 0x18, 0xa0, 0x03,
		0x20, 0x0b, 0x02, 0xa2, 0x10, 0xa0, 0x03, 0x20,
		0x0a, 0x02, 0xc6, 0x04, 0xd0, 0xd2, 0x20, 0x4e,
		0x04, 0x30, 0x1e, 0xa9, 0x01, 0x18, 0x65, 0x18,
		0x85, 0x18, 0xa9, 0x00, 0x65, 0x19, 0x85, 0x19,
		0xa9, 0x00, 0x65, 0x1a, 0x85, 0x1a, 0x10, 0x09,
		0xa2, 0x17, 0xa0, 0x03, 0x20, 0x15, 0x02, 0xe6,
		0x0b, 0xa2, 0x07, 0x86, 0x02, 0xa2, 0x17, 0x86,
		0x00, 0xa2, 0x04, 0x4c, 0x93, 0x03, 0xa2, 0x14,
		0x86, 0x02, 0xa2, 0x08, 0x86, 0x00, 0xa2, 0x03,
		0x20, 0x4a, 0x02, 0xa2, 0x14, 0x86, 0x02, 0xa2,
		0x10, 0x86, 0x00, 0xa0, 0x00, 0xa2, 0x03, 0x38,
		0xb1, 0x00, 0xf1, 0x02, 0x91, 0x02, 0xc8, 0xca,
		0xd0, 0xf6, 0xa5, 0x16, 0x60, 0xa9, 0x00, 0x85,
		0x01, 0x85, 0x03, 0xd8, 0xa2, 0x1c, 0x86, 0x02,
		0xa2, 0x0c, 0x20, 0x00, 0x02, 0x20, 0x80, 0x07,
		0xc9, 0xab, 0xf0, 0x06, 0xc9, 0xad, 0xd0, 0x08,
		0x85, 0x1c, 0x20, 0xc5, 0x07, 0x20, 0x80, 0x07,
		0xc9, 0x8f, 0xd0, 0x0b, 0xa9, 0xbc, 0x20, 0xc5,
		0x07, 0x20, 0xc0, 0x07, 0x4c, 0x75, 0x04, 0xc9,
		0xae, 0xd0, 0x12, 0x24, 0x1e, 0x10, 0x02, 0x30,
		0x2c, 0x85, 0x1e, 0xa0, 0x00, 0x84, 0x04, 0x20,
		0xc5, 0x07, 0x4c, 0x95, 0x04, 0xc9, 0xc5, 0xd0,
		0x42, 0x20, 0xc5, 0x07, 0x20, 0x80, 0x07, 0xc9,
		0xab, 0xf0, 0x06, 0xc9, 0xad, 0xd0, 0x08, 0x85,
		0x1d, 0x20, 0xc5, 0x07, 0x20, 0x80, 0x07, 0xc9,
		0x8f, 0xf0, 0xc1, 0xc9, 0xb0, 0x30, 0x56, 0xc9,
		0xba, 0x10, 0x52, 0x29, 0x0f, 0x85, 0x2c, 0xa2,
		0x27, 0xa9, 0x03, 0xd5, 0x00, 0x30, 0x46, 0xb5,
		0x00, 0x18, 0x36, 0x00, 0x36, 0x00, 0x75, 0x00,
		0x2a, 0x65, 0x2c, 0x95, 0x00, 0xa9, 0xb0, 0x05,
		0x2c, 0xd0, 0xce, 0xc9, 0xb0, 0x30, 0x2e, 0xc9,
		0xba, 0x10, 0x2a, 0xa8, 0xa9, 0xf8, 0x24, 0x25,
		0xd0, 0x83, 0x98, 0x20, 0xc5, 0x07, 0xe6, 0x04,
		0x29, 0x0f, 0x48, 0x20, 0xab, 0x05, 0xa2, 0x23,
		0x68, 0x18, 0x75, 0x00, 0x95, 0x00, 0xa9, 0x00,
		0x75, 0x01, 0x95, 0x01, 0xa9, 0x00, 0x75, 0x02,
		0x95, 0x02, 0x4c, 0x95, 0x04, 0xa5, 0x1c, 0xf0,
		0x07, 0xa2, 0x23, 0xa0, 0x03, 0x20, 0x20, 0x02,
		0xa9, 0x00, 0x85, 0x22, 0xa9, 0x07, 0x85, 0x02,
		0xa9, 0x22, 0x85, 0x00, 0xa2, 0x04, 0x20, 0x4a,
		0x02, 0xa0, 0x17, 0x84, 0x0b, 0x20, 0x55, 0x02,
		0xa5, 0x1d, 0xf0, 0x08, 0xa9, 0xff, 0x45, 0x27,
		0x85, 0x27, 0xe6, 0x27, 0xa5, 0x1e, 0xf0, 0x05,
		0xa9, 0x00, 0x38, 0xe5, 0x04, 0x18, 0x65, 0x27,
		0x85, 0x27, 0x30, 0x1d, 0xd0, 0x01, 0x60, 0x20,
		0x7d, 0x05, 0xd0, 0xfb, 0x60, 0xa9, 0x04, 0x85,
		0x13, 0xa9, 0x50, 0x85, 0x12, 0xa9, 0x00, 0x85,
		0x11, 0x85, 0x10, 0x20, 0x37, 0x03, 0xc6, 0x27,
		0x60, 0x20, 0x97, 0x05, 0xd0, 0xfb, 0x60, 0xa9,
		0xfd, 0x85, 0x13, 0xa9, 0x66, 0x85, 0x12, 0x85,
		0x11, 0xa9, 0x67, 0x85, 0x10, 0x20, 0x37, 0x03,
		0xe6, 0x27, 0x60, 0xa9, 0x00, 0x85, 0x26, 0xa2,
		0x1f, 0x86, 0x02, 0xa2, 0x23, 0x86, 0x00, 0xa2,
		0x04, 0x20, 0x4a, 0x02, 0xa2, 0x23, 0xa0, 0x04,
		0x20, 0x0a, 0x02, 0xa2, 0x23, 0xa0, 0x04, 0x20,
		0x0a, 0x02, 0xa2, 0x1f, 0x86, 0x00, 0xa2, 0x23,
		0x86, 0x02, 0xa2, 0x04, 0x20, 0x2e, 0x02, 0xa2,
		0x23, 0xa0, 0x04, 0x4c, 0x0a, 0x02, 0xa9, 0x00,
		0x85, 0x27, 0xa5, 0x0a, 0x30, 0x04, 0xa9, 0xab,
		0xd0, 0x09, 0xa2, 0x08, 0xa0, 0x03, 0x20, 0x20,
		0x02, 0xa9, 0xad, 0x20, 0xc5, 0x07, 0xa9, 0xb0,
		0x20, 0xc5, 0x07, 0xa9, 0xae, 0x20, 0xc5, 0x07,
		0xc6, 0x0b, 0x10, 0x0f, 0xa9, 0x04, 0x18, 0x65,
		0x0b, 0x10, 0x0e, 0x20, 0x7d, 0x05, 0xa5, 0x0b,
		0x4c, 0x02, 0x06, 0x20, 0x97, 0x05, 0x4c, 0x0e,
		0x06, 0xa2, 0x23, 0x86, 0x02, 0xa2, 0x08, 0x86,
		0x00, 0xa2, 0x03, 0x20, 0x4a, 0x02, 0xa9, 0x00,
		0x85, 0x26, 0xa2, 0x23, 0xa0, 0x03, 0x20, 0x0a,
		0x02, 0x20, 0xab, 0x05, 0xe6, 0x0b, 0xf0, 0x0a,
		0xa2, 0x26, 0xa0, 0x04, 0x20, 0x15, 0x02, 0x4c,
		0x34, 0x06, 0xa9, 0x07, 0x85, 0x04, 0xa5, 0x26,
		0xf0, 0x11, 0xa5, 0x26, 0x09, 0xb0, 0x20, 0xc5,
		0x07, 0xc6, 0x04, 0xf0, 0x1a, 0x20, 0xab, 0x05,
		0x4c, 0x4a, 0x06, 0xc6, 0x27, 0xa5, 0x25, 0xd0,
		0xf0, 0xa5, 0x24, 0xd0, 0xec, 0xa5, 0x23, 0xd0,
		0xe8, 0xa9, 0x00, 0x85, 0x27, 0xf0, 0xe2, 0xa9,
		0xc5, 0x20, 0xc5, 0x07, 0xa5, 0x27, 0x30, 0x05,
		0xa9, 0xab, 0x4c, 0x85, 0x06, 0x49, 0xff, 0x85,
		0x27, 0xe6, 0x27, 0xa9, 0xad, 0x20, 0xc5, 0x07,
		0xa0, 0x00, 0xa5, 0x27, 0x38, 0xe9, 0x0a, 0x30,
		0x06, 0x85, 0x27, 0xc8, 0x4c, 0x8c, 0x06, 0x98,
		0x09, 0xb0, 0x20, 0xc5, 0x07, 0xa5, 0x27, 0x09,
		0xb0, 0x4c, 0xc5, 0x07, 0xa9, 0x8d, 0x20, 0xc5,
		0x07, 0xa9, 0x8a, 0x20, 0xc5, 0x07, 0x20, 0x75,
		0x04, 0x20, 0xc0, 0x07, 0xa2, 0x28, 0x86, 0x02,
		0xa2, 0x08, 0x86, 0x00, 0xa2, 0x04, 0x20, 0x4a,
		0x02, 0x20, 0x80, 0x07, 0xc9, 0xab, 0xd0, 0x09,
		0x20, 0x01, 0x07, 0x20, 0x9a, 0x02, 0x4c, 0xf5,
		0x06, 0xc9, 0xad, 0xd0, 0x09, 0x20, 0x01, 0x07,
		0x20, 0x2d, 0x03, 0x4c, 0xf5, 0x06, 0xc9, 0xd8,
		0xd0, 0x09, 0x20, 0x01, 0x07, 0x20, 0x37, 0x03,
		0x4c, 0xf5, 0x06, 0xc9, 0xaf, 0xd0, 0x0c, 0x20,
		0x01, 0x07, 0x20, 0xdc, 0x03, 0x20, 0xde, 0x05,
		0x4c, 0xa4, 0x06, 0xc9, 0x8f, 0xd0, 0xc2, 0xf0,
		0xa3, 0x20, 0xc5, 0x07, 0x20, 0xc0, 0x07, 0x20,
		0x75, 0x04, 0x20, 0xc0, 0x07, 0xa9, 0xbd, 0x20,
		0xc5, 0x07, 0x20, 0xc0, 0x07, 0xa2, 0x10, 0x86,
		0x02, 0xa2, 0x28, 0x86, 0x00, 0xa2, 0x04, 0x4c,
		0x4a, 0x02,
	}

	adr := uint16(0x200)
	var ofs int

	for ofs < len(code) {
		sIns, n := cpu.Disassemble(adr, code[ofs:])
		sMem := memString(adr, code[ofs:ofs+n])
		sSym := fmt.Sprintf("%-8s", symbols[adr])
		fmt.Printf("%s %s %s\n", sMem, sSym, sIns)
		ofs += n
		adr += uint16(n)
	}

}

//-----------------------------------------------------------------------------
