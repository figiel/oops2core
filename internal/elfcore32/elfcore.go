package elfcore32

import (
	"bytes"
	"debug/elf"
	"encoding/binary"

	c "opensource.go.fig.lu/oops2core/internal/common"
)

const armRegsCount = 18

// PRStatus reflects struct elf_prstatus from Linux UAPI.
// This contains architecture-specific pr_reg (Regs field).
type PRStatus struct {
	SigInfo struct {
		SISigno int32
		SICode  int32
		SIErrno int32
	}
	CurSig  int16
	Pad0    int16
	SigPend uint32
	SigHold uint32
	PID     int32
	PPID    int32
	PGrp    int32
	PSID    int32
	UTime   uint64
	STime   uint64
	CUTime  uint64
	CSTime  uint64
	Regs    [armRegsCount]uint32
	FPValid int32
}

const prog32Size = 32 // bytes

func getElfHeader() elf.Header32 {
	var h elf.Header32

	copy(h.Ident[0:4], elf.ELFMAG[0:4])
	h.Ident[elf.EI_CLASS] = byte(elf.ELFCLASS32)
	h.Ident[elf.EI_DATA] = byte(elf.ELFDATA2LSB)
	h.Ident[elf.EI_VERSION] = byte(elf.EV_CURRENT)
	h.Type = uint16(elf.ET_CORE)
	h.Machine = uint16(elf.EM_ARM)
	h.Version = uint32(elf.EV_CURRENT)
	h.Ehsize = uint16(binary.Size(h))
	h.Phoff = uint32(h.Ehsize)
	h.Phentsize = prog32Size
	h.Phnum = 2 // notes and the stack

	return h
}

func getPRStatus(r *c.ARMRegisters) PRStatus {
	var pr PRStatus

	for i := 0; i < 16; i++ {
		pr.Regs[i] = r.R[int(c.R0)+i]
	}
	pr.Regs[16] = r.CPSR
	pr.Regs[17] = r.R[c.R0]
	pr.CurSig = 11
	pr.SigInfo.SISigno = 11

	return pr
}

type elf32Writer bytes.Buffer

func (b *elf32Writer) mustAdd(what interface{}) {
	err := binary.Write((*bytes.Buffer)(b), binary.LittleEndian, what)
	if err != nil {
		panic(err)
	}
}

// NewElfcore returns an ELF binary buffer with the core file
// built from the provided crashInfo.
func NewElfcore(crashInfo c.CrashInfo) *bytes.Buffer {
	var noteProgHeader elf.Prog32
	var stackProgHeader elf.Prog32
	var elfcore elf32Writer

	elfHeader := getElfHeader()
	prStatus := getPRStatus(&crashInfo.Regs)

	// This structure is roughly documented in Sys V ABI,
	// Exact values in this structure mimic Linux kernel core dumper.
	noteInfo := struct {
		namesz uint32
		descsz uint32
		ntype  uint32
		name   [8]byte // has to be aligned to 4 bytes
	}{5, uint32(binary.Size(prStatus)), uint32(elf.NT_PRSTATUS), [8]byte{'C', 'O', 'R', 'E', 0}}

	// Notes segment is just after the program headers and contains
	// note information and the PRStatus structure
	notesDataOffset := uint32(elfHeader.Ehsize + elfHeader.Phentsize*elfHeader.Phnum)
	notesSegmentSize := uint32(binary.Size(noteInfo) + binary.Size(prStatus))

	// Stack segment is directly after the notes segment
	stackDataOffset := notesDataOffset + notesSegmentSize
	stackSegmentSize := uint32(binary.Size(crashInfo.Stack))

	noteProgHeader.Type = uint32(elf.PT_NOTE)
	noteProgHeader.Off = notesDataOffset
	noteProgHeader.Filesz = notesSegmentSize

	stackProgHeader.Type = uint32(elf.PT_LOAD)
	stackProgHeader.Off = stackDataOffset
	stackProgHeader.Filesz = stackSegmentSize
	stackProgHeader.Memsz = stackSegmentSize
	stackProgHeader.Vaddr = crashInfo.Regs.R[c.SP]

	elfcore.mustAdd(elfHeader)
	elfcore.mustAdd(noteProgHeader)
	elfcore.mustAdd(stackProgHeader)
	elfcore.mustAdd(noteInfo)
	elfcore.mustAdd(prStatus)
	elfcore.mustAdd(crashInfo.Stack)

	return (*bytes.Buffer)(&elfcore)
}
