package common

// ARMRegisters hold contents of CPU registers at the moment
// of the crash.
type ARMRegisters struct {
	R    [16]uint32
	CPSR uint32
}

// Stack contains stack contents at the moment of the crash.
// First element of the array is at the stack pointer (on top
// of the stack). Subsequent elements are at addresses sp+4, sp+8,
// etc. (i.e. downwards in stack).
type Stack []uint32

// CrashInfo is not documented.
type CrashInfo struct {
	Regs  ARMRegisters
	Stack Stack
}

// RegisterIndex is not documented.
type RegisterIndex int

// Definition of register indices.
const (
	R0 RegisterIndex = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	FP
	IP
	SP
	LR
	PC
)
