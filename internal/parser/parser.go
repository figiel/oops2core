package parser

import (
	"fmt"
	"regexp"
	"strconv"

	c "opensource.go.fig.lu/oops2core/internal/common"
)

type armRegisters c.ARMRegisters

var crashRegexp = regexp.MustCompile(`(?sm)` +
	` pc : \[<([[:xdigit:]]+)>\].*?` + // 1
	` lr : \[<([[:xdigit:]]+)>\].*?` +
	` psr: ([[:xdigit:]]+).*?` +
	` sp : ([[:xdigit:]]+).*?` +
	` ip : ([[:xdigit:]]+).*?` + // 5
	` fp : ([[:xdigit:]]+).*?` +
	` r10: ([[:xdigit:]]+).*?` +
	` r9 : ([[:xdigit:]]+).*?` +
	` r8 : ([[:xdigit:]]+).*?` +
	` r7 : ([[:xdigit:]]+).*?` + // 10
	` r6 : ([[:xdigit:]]+).*?` +
	` r5 : ([[:xdigit:]]+).*?` +
	` r4 : ([[:xdigit:]]+).*?` +
	` r3 : ([[:xdigit:]]+).*?` +
	` r2 : ([[:xdigit:]]+).*?` + // 15
	` r1 : ([[:xdigit:]]+).*?` +
	` r0 : ([[:xdigit:]]+).*?` +
	`Stack: \(.*?\)$\r?\n((^.* [[:xdigit:]]+: [[:xdigit:] ]+$(\r?\n)?)+)`)

const crashRegexpStackSubgroup = 18

var stackElementRegexp = regexp.MustCompile(`([[:xdigit:]]{8})( |\n)`)

func parseWord(text string) (uint32, error) {
	ret, err := strconv.ParseUint(text, 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(ret), nil
}

func parseRegisters(m []string) (c.ARMRegisters, error) {
	var ret c.ARMRegisters
	var err error

	orderInText := []*uint32{
		&ret.R[c.PC], &ret.R[c.LR], &ret.CPSR,
		&ret.R[c.SP], &ret.R[c.IP], &ret.R[c.FP],
		&ret.R[c.R10], &ret.R[c.R9], &ret.R[c.R8],
		&ret.R[c.R7], &ret.R[c.R6], &ret.R[c.R5],
		&ret.R[c.R4], &ret.R[c.R3], &ret.R[c.R2],
		&ret.R[c.R1], &ret.R[c.R0]}

	for i, reg := range orderInText {
		*reg, err = parseWord(m[i+1])
		if err != nil {
			return c.ARMRegisters{}, err
		}
	}
	return ret, nil
}

func parseStack(m []string) (c.Stack, error) {
	var ret c.Stack
	elems := stackElementRegexp.FindAllStringSubmatch(m[crashRegexpStackSubgroup], -1)
	for _, z := range elems {
		w, err := parseWord(z[1])
		if err != nil {
			return c.Stack{}, nil
		}
		ret = append(ret, w)
	}
	return ret, nil
}

// ParseCrash extracts information about registers
// and stack from the provided crash report text.
func ParseCrash(crashText string) (c.CrashInfo, error) {
	m := crashRegexp.FindStringSubmatch(crashText)

	if len(m) < crashRegexp.NumSubexp() {
		return c.CrashInfo{}, fmt.Errorf("can't parse crash text")
	}

	regs, err := parseRegisters(m)
	if err != nil {
		return c.CrashInfo{}, err
	}

	stack, err := parseStack(m)
	if err != nil {
		return c.CrashInfo{}, err
	}

	return c.CrashInfo{Regs: regs, Stack: stack}, nil
}
