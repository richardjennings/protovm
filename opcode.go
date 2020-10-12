package protovm

import "unsafe"

// Operations
type Opcode uint64
type Funct uint64
type Reg uint64

const (
	None Funct = iota
	Immediate
	Int
	ImmediateInt
	IntImmediate
	Float
	ImmediateFloat
	FloatImmediate
	Bool
	ImmediateBool
	Ptr
	SP
	SPR
)

const (
	NoOp Opcode = iota

	And
	Or
	Not

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

	Exit
)

func Op(v Opcode) [8]byte {
	return Uint64(uint64(v))
}
func F(v Funct) [8]byte {
	return Uint64(uint64(v))
}
func R(v Reg) [8]byte {
	return Uint64(uint64(v))
}
func Uint64(v uint64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func Int64(v int64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func Float64(v float64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func Boolean(v bool) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
