package proto

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

const ProgramCounter = int(20)
const StackPointer = 21

// Registers is defined as 20 8 byte arrays
type Registers [22][8]byte

// RAM defines the stack size
type RAM [][8]byte

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
	Invalid Opcode = iota

	Call
	NoOp

	And
	Or
	Not

	Band
	Bor
	Bnot
	Bxor

	ShiftL
	ShiftR

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
	NoOp:    "noop",
	Call:    "call",
	And:     "and",
	Or:      "or",
	Not:     "not",
	Band:    "band",
	Bor:     "bor",
	Bnot:    "bnot",
	Bxor:    "bxor",
	ShiftL:  "shiftl",
	ShiftR:  "shiftr",
	Add:     "add",
	Sub:     "sub",
	Mul:     "mul",
	Quo:     "quo",
	Pow:     "pow",
	Rem:     "rem",
	Eq:      "eq",
	NEq:     "neq",
	LT:      "lt",
	LTE:     "lte",
	GT:      "gt",
	GTE:     "gte",
	Print:   "print",
	PrintLn: "println",
	Load:    "load",
	Store:   "store",
	JMP:     "jmp",
	JMPEQ:   "je",
	JMPNEQ:  "jne",
	Exit:    "exit",
}

func GetOpcode(o string) Opcode {
	for i, v := range opcodes {
		if v == o {
			return i
		}
	}
	return Invalid
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
