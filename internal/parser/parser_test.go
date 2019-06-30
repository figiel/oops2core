// SPDX-License-Identifier: MIT
package parser

import (
	"testing"

	c "opensource.go.fig.lu/oops2core/internal/common"
)

const testCrashSnippet = `[    2.472742] pc : [<c052fa4c>]    lr : [<c006d464>]    psr: 60000113
[    2.474303] sp : cf9cfef0  ip : cf9cfe38  fp : cf9cff04
[    2.476136] r10: c0bc7d74  r9 : 99999999  r8 : 88888888
[    2.477646] r7 : cfdaa200  r6 : c0b43050  r5 : c0ba6f24  r4 : 00000001
[    2.479216] r3 : 33333333  r2 : c7d5b9c6  r1 : 11111111  r0 : 00000005
[    2.480871] Flags: nZCv  IRQs on  FIQs on  Mode SVC_32  ISA ARM  Segment none
[    2.482523] Control: 00093177  Table: 00004000  DAC: 00000053
[    2.484159] Process kworker/0:2 (pid: 17, stack limit = 0x(ptrval))
[    2.485854] Stack: (0xcf9cfef0 to 0xcf9d0000)
[    2.487561] fee0:                                     cf9c3120 c0ba6f24 cf9cff14 cf9cff08
[    2.490814] ff00: c052fa84 c052fa34 cf9cff4c cf9cff18 c0045338 c052fa70 c0b4bae0 ffffe000
[    2.495144] ff20: c0b43064 cf9c3120 cf9c3134 c0b43050 c0b4bae0 c0b43050 c0b43064 00000008
[    2.498513] ff40: cf9cff7c cf9cff50 c004610c c00450b4 cf915be0 cf97ff40 cf9c9020 cf9ce000
[    2.501931] ff60: cf9c3120 c0045df0 cf9cde90 cf97ff58 cf9cffac cf9cff80 c004b1d8 c0045dfc
[    2.515526] ff80: cf9ce000 cf9c9020 c004b0ac 00000000 00000000 00000000 00000000 00000000
[    2.519240] ffa0: 00000000 cf9cffb0 c00090e0 c004b0b8 00000000 00000000 00000000 00000000
[    2.522947] ffc0: 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000
[    2.526655] ffe0: 00000000 00000000 00000000 00000000 00000013 00000000 00000000 00000000
[    2.530815] [<c052fa4c>] (really_explode) from [<c052fa84>] (explode+0x20/0x24)
[    2.534683] [<c052fa84>] (explode) from [<c0045338>] (process_one_work+0x290/0x44c)
[    2.538505] [<c0045338>] (process_one_work) from [<c004610c>] (worker_thread+0x31c/0x4c8)
`

func testGetRegs(t *testing.T, snippet string) c.ARMRegisters {
	ret, err := ParseCrash(snippet)
	if err != nil {
		t.Fatal("Can't parse test snippet: ", err)
	}
	return ret.Regs
}

func testExpectFailure(t *testing.T, snippet string) {
	_, err := ParseCrash(snippet)

	if err == nil {
		t.Fatal("Expected parse failure but got success")
	}
}

func assertEq(t *testing.T, expected uint32, got uint32, s string) {
	if got != expected {
		t.Errorf("Assertion failed: expected %08x but got %08x (%s).\n",
			expected, got, s)
	}
}

func assertEqArr(t *testing.T, expected []uint32, got []uint32) {
	lenExpected := len(expected)
	lenGot := len(got)
	if len(expected) != len(got) {
		t.Errorf("Assertion failed: arrays of different length, len(expected)=%d, len(got)=%d\n",
			lenExpected, lenGot)
		return
	}
	for i := range expected {
		if expected[i] != got[i] {
			t.Errorf("Assertion failed: element %d differs, expected=%d, got=%d\n", i, expected[i], got[i])
		}
	}
}

func TestParseRegistersOK(t *testing.T) {

	r := testGetRegs(t, testCrashSnippet)

	assertEq(t, 0x00000005, r.R[c.R0], "R0")
	assertEq(t, 0x11111111, r.R[c.R1], "R1")
	assertEq(t, 0xc7d5b9c6, r.R[c.R2], "R2")
	assertEq(t, 0x33333333, r.R[c.R3], "R3")
	assertEq(t, 0x00000001, r.R[c.R4], "R4")
	assertEq(t, 0xc0ba6f24, r.R[c.R5], "R5")
	assertEq(t, 0xc0b43050, r.R[c.R6], "R6")
	assertEq(t, 0xcfdaa200, r.R[c.R7], "R7")
	assertEq(t, 0x88888888, r.R[c.R8], "R8")
	assertEq(t, 0x99999999, r.R[c.R9], "R9")
	assertEq(t, 0xc0bc7d74, r.R[c.R10], "R10")
	assertEq(t, 0xcf9cff04, r.R[c.FP], "FP")
	assertEq(t, 0xcf9cfe38, r.R[c.IP], "IP")
	assertEq(t, 0xcf9cfef0, r.R[c.SP], "SP")
	assertEq(t, 0xc006d464, r.R[c.LR], "LR")
	assertEq(t, 0xc052fa4c, r.R[c.PC], "PC")
	assertEq(t, 0x60000113, r.CPSR, "CPSR")
}

func TestParseRegistersNoMatch(t *testing.T) {
	const testRegsSnippet = "[    2.472742] pc : [<c052fa4c>]    lr : [<c006d464>]    psr: 60000113\n" +
		"[    2.474303] sp : cf9cfef0  ip : cf9cfe38  fp : cf9cff04\n" +
		"[    2.476136] r10: c0bc7d74  r9 : 99999999  r8 : 88888888\n"

	testExpectFailure(t, testRegsSnippet)
}

func TestParseRegistersCantParseValue(t *testing.T) {
	const testRegsSnippet = "[    2.472742] pc : [<aaa>]    lr : [<bbb>]    psr: ccc\n" +
		"[    2.474303] sp : 99999999999  ip : 00000000  fp : 00000000\n" +
		"[    2.476136] r10: 00000000  r9 : 00000000  r8 : 00000000\n" +
		"[    2.477646] r7 : 00000000  r6 : 00000000  r5 : 00000000  r4 : 00000000\n" +
		"[    2.479216] r3 : 00000000  r2 : 00000000  r1 : 00000000  r0 : 00000000\n" +
		"[    2.485854] Stack: (0xcf9cfef0 to 0xcf9d0000)\n" +
		"[    2.487561] fee0: cf9c3120 c0ba6f24 cf9cff14 cf9cff08\n"

	testExpectFailure(t, testRegsSnippet)
}

func TestParseStack(t *testing.T) {
	crashInfo, err := ParseCrash(testCrashSnippet)
	expected := []uint32{
		0xcf9c3120, 0xc0ba6f24, 0xcf9cff14, 0xcf9cff08,
		0xc052fa84, 0xc052fa34, 0xcf9cff4c, 0xcf9cff18, 0xc0045338, 0xc052fa70, 0xc0b4bae0, 0xffffe000,
		0xc0b43064, 0xcf9c3120, 0xcf9c3134, 0xc0b43050, 0xc0b4bae0, 0xc0b43050, 0xc0b43064, 0x00000008,
		0xcf9cff7c, 0xcf9cff50, 0xc004610c, 0xc00450b4, 0xcf915be0, 0xcf97ff40, 0xcf9c9020, 0xcf9ce000,
		0xcf9c3120, 0xc0045df0, 0xcf9cde90, 0xcf97ff58, 0xcf9cffac, 0xcf9cff80, 0xc004b1d8, 0xc0045dfc,
		0xcf9ce000, 0xcf9c9020, 0xc004b0ac, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0xcf9cffb0, 0xc00090e0, 0xc004b0b8, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000000,
		0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00000013, 0x00000000, 0x00000000, 0x00000000,
	}
	if err != nil {
		t.Fatal("ParseCrash failed: ", err)
	}
	assertEqArr(t, expected, []uint32(crashInfo.Stack))
}
