package vm

import (
	"fmt"
	"unsafe"
)

// Opcode represents an Op
type Opcode uint8

// Funct defines types of Instruction format
type Funct uint8

// Reg is a general purpose Register
type Reg uint64

// ByteCode is a slice of instructions
type ByteCode []Inst

// Inst represents an instruction
// type Inst [5][8]byte
type Inst struct {
	O Opcode
	F Funct
	X uint64
	Y uint64
	Z uint64
}

// Registers is defined as 20 8 byte arrays
type Registers [20][8]byte

// Stack defines the stack size
type Stack [100][8]byte

const (
	None Funct = iota
	Imm
	Int
	ImmI
	IImm
	Float
	ImmF
	FImm
	Bool
	ImmB
	SP
	SPR
)

var functs = map[Funct]string{
	None:  "",
	Imm:   "Imm",
	Int:   "Int",
	ImmI:  "ImmI",
	IImm:  "IImm",
	Float: "Float",
	ImmF:  "ImmF",
	FImm:  "FImm",
	Bool:  "Bool",
	ImmB:  "ImmB",
	SP:    "SP",
	SPR:   "SPR",
}

var __ [8]byte

const (
	NoOp Opcode = iota

	And
	Or
	Not

	Band
	Bor
	Bnot
	Bxor
	Bclr

	Add
	Sub
	Mul
	Quo
	Pow
	Rem

	Eq
	NEq
	LT
	LTE
	GT
	GTE

	Print
	PrintLn

	Load
	Store

	JMP
	JMPEQ
	JMPNEQ

	Exit
)

var opcodes = map[Opcode]string{
	NoOp:    "NoOp",
	And:     "And",
	Or:      "Or",
	Not:     "Not",
	Band:    "Band",
	Bor:     "Bor",
	Bnot:    "Bnot",
	Bxor:    "Bxor",
	Bclr:    "Bclr",
	Add:     "Add",
	Sub:     "Sub",
	Mul:     "Mul",
	Quo:     "Quo",
	Pow:     "Pow",
	Rem:     "Rem",
	Eq:      "Eq",
	NEq:     "NEq",
	LT:      "LT",
	LTE:     "LTE",
	GT:      "GT",
	GTE:     "GTE",
	Print:   "Print",
	PrintLn: "PrintLn",
	Load:    "Load",
	Store:   "Store",
	JMP:     "JMP",
	JMPEQ:   "JMPEQ",
	JMPNEQ:  "JMPNEQ",
	Exit:    "Exit",
}

func (f Funct) String() string {
	if n, ok := functs[f]; ok {
		return n
	}
	return "ERROR: func string not found"
}
func (o Opcode) String() string {
	if n, ok := opcodes[o]; ok {
		return n
	}
	return "ERROR: opcode string not found"
}
func (b ByteCode) String() string {
	s := ""
	for i, bc := range b {
		s += fmt.Sprintf("%d %d %d %d %d %d\n", i, bc.O, bc.F, bc.X, bc.Y, bc.Z)
	}
	return s
}

func R(r Reg) uint64 {
	return uint64(r)
}

func Int64(v int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&v))
}
func Float64(v float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&v))
}
func Boolean(v bool) uint64 {
	return *(*uint64)(unsafe.Pointer(&v))
}
