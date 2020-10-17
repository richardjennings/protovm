package vm

import (
	"bytes"
	"fmt"
	"testing"
)

type tcase struct {
	bc ByteCode
	o  string
}

func testcase(t tcase) error {
	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	err := vm.Exec(t.bc)
	if err != nil {
		return err
	}
	out := string(buf.Bytes())
	if out != t.o {
		return fmt.Errorf("expected %s got %s", t.o, out)
	}
	return nil
}

func test(tc []tcase, t *testing.T) {
	for i, c := range tc {
		err := testcase(c)
		if err != nil {
			t.Error(err, fmt.Sprintf("test case: %d", i))
		}
	}
}

func Test_NoOp(t *testing.T) {
	test([]tcase{
		{
			// NoOp does nothing other than increase PC
			ByteCode{
				{},
				{Op(Exit)},
			},
			"",
		},
	}, t)
}


func Test_And(t *testing.T) {
	test([]tcase{
		{
			// And Bool
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(true), R(Reg(2))},
				{Op(Store), F(None), Boolean(false), R(Reg(4))},
				{Op(Store), F(None), Boolean(true), R(Reg(5))},
				{Op(And), F(Bool), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(And), F(Bool), R(Reg(4)), R(Reg(5)), R(Reg(6))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// And ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(And), F(ImmB), Boolean(true), R(Reg(1)), R(Reg(2))},
				{Op(And), F(ImmB), Boolean(false), R(Reg(1)), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(2))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Exit)},
			},
			"truefalse",
		},
	}, t)
}

func Test_Or(t *testing.T) {
	test([]tcase{
		{
			// Or Bool
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(Store), F(None), Boolean(false), R(Reg(4))},
				{Op(Store), F(None), Boolean(false), R(Reg(5))},
				{Op(Or), F(Bool), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Or), F(Bool), R(Reg(4)), R(Reg(5)), R(Reg(6))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// Or ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Or), F(ImmB), Boolean(true), R(Reg(1)), R(Reg(2))},
				{Op(Or), F(ImmB), Boolean(false), R(Reg(1)), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(2))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Exit)},
			},
			"truefalse",
		},
	}, t)
}

func Test_Not(t *testing.T) {
	test([]tcase{
		{
			// Not
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(Not), F(Bool), R(Reg(1)), [8]byte{}, R(Reg(3))},
				{Op(Not), F(Bool), R(Reg(2)), [8]byte{}, R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
	}, t)
}

func Test_Add(t *testing.T) {
	test([]tcase{
		{
			// Add Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Add), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"30",
		},
		{
			// Add ImmI
			ByteCode{
				// store int 10 in reg 1,
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Add), F(ImmI), Int64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"30",
		},
		{
			// Add Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Add), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"30",
		},
		{
			// Add ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Add), F(ImmF), Float64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"30",
		},
	}, t)
}

func Test_Sub(t *testing.T) {
	test([]tcase{
		{
			// Sub Int
			ByteCode{
				{Op(Store), F(None), Int64(30), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Sub), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Sub), F(ImmI), Int64(30), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub IImm
			ByteCode{
				{Op(Store), F(None), Int64(30), R(Reg(1))},
				{Op(Sub), F(IImm), R(Reg(1)), Int64(10), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub Float
			ByteCode{
				{Op(Store), F(None), Float64(30), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(Sub), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Sub), F(ImmF), Float64(30), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub FImm
			ByteCode{
				{Op(Store), F(None), Float64(30), R(Reg(1))},
				{Op(Sub), F(FImm), R(Reg(1)), Float64(10), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
	}, t)
}

func Test_Mul(t *testing.T) {
	test([]tcase{
		{
			// Mul Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Mul), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"200",
		},
		{
			// Mul ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Mul), F(ImmI), Int64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"200",
		},
		{
			// Mul Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Mul), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"200",
		},
		{
			// Mul ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Mul), F(ImmF), Float64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"200",
		},
	}, t)
}

func Test_Quo(t *testing.T) {
	test([]tcase{
		{
			// Quo Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Quo), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Quo), F(ImmI), Int64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo IImm
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Quo), F(IImm), R(Reg(1)), Int64(10), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo Float
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(Quo), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Quo), F(ImmF), Float64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo FImm
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Quo), F(FImm), R(Reg(1)), Float64(10), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_Pow(t *testing.T) {
	test([]tcase{
		{
			// Pow Int
			ByteCode{
				{Op(Store), F(None), Int64(2), R(Reg(1))},
				{Op(Store), F(None), Int64(3), R(Reg(2))},
				{Op(Pow), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow ImmI
			ByteCode{
				{Op(Store), F(None), Int64(3), R(Reg(1))},
				{Op(Pow), F(ImmI), Int64(2), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow IImm
			ByteCode{
				{Op(Store), F(None), Int64(2), R(Reg(1))},
				{Op(Pow), F(IImm), R(Reg(1)), Int64(3), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow Float
			ByteCode{
				{Op(Store), F(None), Float64(2), R(Reg(1))},
				{Op(Store), F(None), Float64(3), R(Reg(2))},
				{Op(Pow), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow ImmF
			ByteCode{
				{Op(Store), F(None), Float64(3), R(Reg(1))},
				{Op(Pow), F(ImmF), Float64(2), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow FImm
			ByteCode{
				{Op(Store), F(None), Float64(2), R(Reg(1))},
				{Op(Pow), F(FImm), R(Reg(1)), Float64(3), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
	}, t)
}

func Test_Rem(t *testing.T) {
	test([]tcase{
		{
			// asm.Rem Int
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(4), R(Reg(2))},
				{Op(Rem), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// asm.Rem ImmI
			ByteCode{
				{Op(Store), F(None), Int64(4), R(Reg(1))},
				{Op(Rem), F(ImmI), Int64(10), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// asm.Rem IImm
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Rem), F(IImm), R(Reg(1)), Int64(4), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// asm.Rem Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(4), R(Reg(2))},
				{Op(Rem), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(Print), F(Float), R(Reg(3))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// asm.Rem ImmF
			ByteCode{
				{Op(Store), F(None), Float64(4), R(Reg(1))},
				{Op(Rem), F(ImmF), Float64(10), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// asm.Rem FImm
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Rem), F(FImm), R(Reg(1)), Float64(4), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_Eq(t *testing.T) {
	test([]tcase{
		{
			// Eq Int
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(Store), F(None), Int64(20), R(Reg(3))},
				{Op(Eq), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(Eq), F(Int), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(Eq), F(ImmI), Int64(11), R(Reg(1)), R(Reg(3))},
				{Op(Eq), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Store), F(None), Float64(20), R(Reg(3))},
				{Op(Eq), F(Float), R(Reg(1)), R(Reg(3)), R(Reg(4))},
				{Op(Eq), F(Float), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Eq), F(ImmF), Float64(10), R(Reg(2)), R(Reg(3))},
				{Op(Eq), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq Bool
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Store), F(None), Boolean(true), R(Reg(2))},
				{Op(Store), F(None), Boolean(true), R(Reg(3))},
				{Op(Eq), F(Bool), R(Reg(1)), R(Reg(3)), R(Reg(4))},
				{Op(Eq), F(Bool), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(Eq), F(ImmB), Boolean(true), R(Reg(2)), R(Reg(3))},
				{Op(Eq), F(ImmB), Boolean(false), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
	}, t)
}

func Test_NEq(t *testing.T) {
	test([]tcase{
		{
			// NEq Int
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(Store), F(None), Int64(20), R(Reg(3))},
				{Op(NEq), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(NEq), F(Int), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(NEq), F(ImmI), Int64(11), R(Reg(1)), R(Reg(3))},
				{Op(NEq), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Store), F(None), Float64(20), R(Reg(3))},
				{Op(NEq), F(Float), R(Reg(1)), R(Reg(3)), R(Reg(4))},
				{Op(NEq), F(Float), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(NEq), F(ImmF), Float64(10), R(Reg(2)), R(Reg(3))},
				{Op(NEq), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq Bool
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Store), F(None), Boolean(true), R(Reg(2))},
				{Op(Store), F(None), Boolean(true), R(Reg(3))},
				{Op(NEq), F(Bool), R(Reg(1)), R(Reg(3)), R(Reg(4))},
				{Op(NEq), F(Bool), R(Reg(2)), R(Reg(3)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(NEq), F(ImmB), Boolean(true), R(Reg(2)), R(Reg(3))},
				{Op(NEq), F(ImmB), Boolean(false), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
	}, t)
}

func Test_LT(t *testing.T) {
	test([]tcase{
		{
			// LT Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(LT), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(LT), F(Int), R(Reg(2)), R(Reg(1)), R(Reg(4))},
				// print reg 3
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(LT), F(ImmI), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(LT), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT Float
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(LT), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(LT), F(Float), R(Reg(2)), R(Reg(1)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(LT), F(ImmF), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(LT), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"falsetrue",
		},
	}, t)
}

func Test_LTE(t *testing.T) {
	test([]tcase{
		{
			// LTE Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Store), F(None), Int64(10), R(Reg(3))},
				{Op(LTE), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(Int), R(Reg(2)), R(Reg(1)), R(Reg(5))},
				{Op(LTE), F(Int), R(Reg(3)), R(Reg(3)), R(Reg(6))},
				// print reg 3
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(LTE), F(ImmI), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(LTE), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(ImmI), Int64(30), R(Reg(2)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE Float
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(Store), F(None), Float64(10), R(Reg(3))},
				{Op(LTE), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(Float), R(Reg(2)), R(Reg(1)), R(Reg(5))},
				{Op(LTE), F(Float), R(Reg(2)), R(Reg(3)), R(Reg(6))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(LTE), F(ImmF), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(LTE), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(ImmF), Float64(30), R(Reg(2)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"falsetruetrue",
		},
	}, t)
}

func Test_GT(t *testing.T) {
	test([]tcase{
		{
			// GT Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(GT), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(GT), F(Int), R(Reg(2)), R(Reg(1)), R(Reg(4))},
				// print reg 3
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// GT ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(GT), F(ImmI), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(GT), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// GT Float
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(GT), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(3))},
				{Op(GT), F(Float), R(Reg(2)), R(Reg(1)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
		{
			// GT ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(GT), F(ImmF), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(GT), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Exit)},
			},
			"truefalse",
		},
	}, t)
}

func Test_GTE(t *testing.T) {
	test([]tcase{
		{
			// GTE Int
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Store), F(None), Int64(10), R(Reg(2))},
				{Op(Store), F(None), Int64(10), R(Reg(3))},
				{Op(GTE), F(Int), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(Int), R(Reg(2)), R(Reg(1)), R(Reg(5))},
				{Op(GTE), F(Int), R(Reg(3)), R(Reg(3)), R(Reg(6))},
				// print reg 3
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmI
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(GTE), F(ImmI), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(GTE), F(ImmI), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(ImmI), Int64(30), R(Reg(2)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE Float
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Store), F(None), Float64(10), R(Reg(2))},
				{Op(Store), F(None), Float64(10), R(Reg(3))},
				{Op(GTE), F(Float), R(Reg(1)), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(Float), R(Reg(2)), R(Reg(1)), R(Reg(5))},
				{Op(GTE), F(Float), R(Reg(2)), R(Reg(3)), R(Reg(6))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(6))},
				{Op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmF
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(GTE), F(ImmF), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(GTE), F(ImmF), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(ImmF), Float64(30), R(Reg(2)), R(Reg(5))},
				{Op(Print), F(Bool), R(Reg(3))},
				{Op(Print), F(Bool), R(Reg(4))},
				{Op(Print), F(Bool), R(Reg(5))},
				{Op(Exit)},
			},
			"truefalsetrue",
		},
	}, t)
}

func Test_Print(t *testing.T) {
	test([]tcase{
		{
			// Print Int
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Print), F(Int), R(Reg(1))},
				{Op(Exit)},
			},
			"10",
		},
		{
			// Print ImmI
			ByteCode{
				{Op(Print), F(ImmI), Int64(10)},
				{Op(Exit)},
			},
			"10",
		},
		{
			// Print Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Print), F(Float), R(Reg(1))},
				{Op(Exit)},
			},
			"10",
		},
		{
			// Print ImmF
			ByteCode{
				{Op(Print), F(ImmF), Float64(10)},
				{Op(Exit)},
			},
			"10",
		},
		{
			// Print Bool
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Store), F(None), Boolean(true), R(Reg(2))},
				{Op(Print), F(Bool), R(Reg(1))},
				{Op(Print), F(Bool), R(Reg(2))},
				{Op(Exit)},
			},
			"falsetrue",
		},
		{
			// Print ImmB
			ByteCode{
				{Op(Print), F(ImmB), Boolean(false)},
				{Op(Print), F(ImmB), Boolean(true)},
				{Op(Exit)},
			},
			"falsetrue",
		},
	}, t)
}

func Test_PrintLn(t *testing.T) {
	test([]tcase{
		{
			// PrintLn Int
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(PrintLn), F(Int), R(Reg(1))},
				{Op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn ImmI
			ByteCode{
				{Op(PrintLn), F(ImmI), Int64(10)},
				{Op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn Float
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(PrintLn), F(Float), R(Reg(1))},
				{Op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn ImmF
			ByteCode{
				{Op(PrintLn), F(ImmF), Float64(10)},
				{Op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn Bool
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Store), F(None), Boolean(true), R(Reg(2))},
				{Op(PrintLn), F(Bool), R(Reg(1))},
				{Op(PrintLn), F(Bool), R(Reg(2))},
				{Op(Exit)},
			},
			"false\ntrue\n",
		},
		{
			// PrintLn ImmB
			ByteCode{
				{Op(PrintLn), F(ImmB), Boolean(false)},
				{Op(PrintLn), F(ImmB), Boolean(true)},
				{Op(Exit)},
			},
			"false\ntrue\n",
		},
	}, t)
}

func Test_Load(t *testing.T) {
	test([]tcase{
		{
			// Load from Ptr
			ByteCode{
				{Op(Store), F(None), R(Reg(9)), R(Reg(1))},
				{Op(Store), F(None), Int64(3), R(Reg(9))},
				{Op(Load), F(Ptr), R(Reg(2)), R(Reg(1))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"3",
		},
		{
			// Load from Stack
			ByteCode{
				{Op(Store), F(SP), Int64(3)},
				{Op(Load), F(SP), R(Reg(1))},
				{Op(Print), F(Int), R(Reg(1))},
				{Op(Exit)},
			},
			"3",
		},
	}, t)
}

func Test_Store(t *testing.T) {
	test([]tcase{
		{
			// Store bytes in register
			ByteCode{
				{Op(Store), F(None), Int64(1), R(Reg(1))},
				{Op(Print), F(Int), R(Reg(1))},
				{Op(Exit)},
			},
			"1",
		},
		{
			// Store bytes in stack from Imm
			ByteCode{
				{Op(Store), F(SP), Int64(1)},
				{Op(Store), F(SP), Int64(2)},
				{Op(Store), F(SP), Int64(3)},
				{Op(Load), F(SP), R(Reg(1))},
				{Op(Load), F(SP), R(Reg(2))},
				{Op(Load), F(SP), R(Reg(3))},
				{Op(Print), F(Int), R(Reg(1))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(3))},
				{Op(Exit)},
			},
			"321",
		},
		{
			// Store bytes in stack from asm.Reg
			ByteCode{
				{Op(Store), F(None), Int64(1), R(Reg(1))},
				{Op(Store), F(SPR), R(Reg(1))},
				{Op(Load), F(SP), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"1",
		},
	}, t)
}

func Test_JMP(t *testing.T) {
	test([]tcase{
		{
			// JMP Imm
			ByteCode{
				{Op(JMP), F(Imm), Uint64(2)},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(JMP), F(Imm), Uint64(5)},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
		{
			// JMP asm.Reg
			ByteCode{
				{Op(Store), F(None), Uint64(4), R(Reg(1))},
				{Op(Store), F(None), Uint64(7), R(Reg(2))},
				{Op(JMP), F(None), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(JMP), F(None), R(Reg(2))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
		{
			// JMP Stack Pointer
			ByteCode{
				{Op(Store), F(SP), Uint64(7)},
				{Op(Store), F(SP), Uint64(4)},
				{Op(JMP), F(SP)},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(JMP), F(SP)},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_JMPEQ(t *testing.T) {
	test([]tcase{
		{
			// ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(JMPEQ), F(ImmB), Uint64(3), Boolean(true), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(JMPEQ), F(ImmB), Uint64(7), Boolean(false), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
		{
			// ImmI
			ByteCode{
				{Op(Store), F(None), Int64(2), R(Reg(1))},
				{Op(JMPEQ), F(ImmI), Uint64(3), Int64(2), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(Store), F(None), Int64(3), R(Reg(1))},
				{Op(JMPEQ), F(ImmI), Uint64(7), Uint64(3), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_JMPNEQ(t *testing.T) {
	test([]tcase{
		{
			// ImmB
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(JMPNEQ), F(ImmB), Uint64(3), Boolean(true), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(JMPNEQ), F(ImmB), Uint64(7), Boolean(false), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
		{
			// ImmI
			ByteCode{
				{Op(Store), F(None), Int64(1), R(Reg(1))},
				{Op(JMPNEQ), F(ImmI), Uint64(3), Int64(2), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Print), F(ImmI), Int64(2)},
				{Op(Store), F(None), Int64(4), R(Reg(1))},
				{Op(JMPNEQ), F(ImmI), Uint64(7), Uint64(3), R(Reg(1))},
				{Op(Print), F(ImmI), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_UnhandledOp(t *testing.T) {
	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	err := vm.Exec(ByteCode{{Uint64(1999)}})
	if err == nil {
		t.Error("expected error, got none")
	}
}


func BenchmarkFibRecursive_0(b *testing.B) {
	benchmarkFibRecursive(0, "0\n", b)
}

func BenchmarkFibRecursive_1(b *testing.B) {
	benchmarkFibRecursive(1, "1\n", b)
}

func BenchmarkFibRecursive_5(b *testing.B) {
	benchmarkFibRecursive(5, "5\n", b)
}

func BenchmarkFibRecursive_20(b *testing.B) {
	benchmarkFibRecursive(20, "6765\n", b)
}

func BenchmarkFibRecursive_30(b *testing.B) {
	benchmarkFibRecursive(30, "832040\n", b)
}

func BenchmarkFibRecursive_35(b *testing.B) {
	benchmarkFibRecursive(35, "9227465\n", b)
}

func benchmarkFibRecursive(n int, e string, b *testing.B) {
	bc := ByteCode{
		{Op(Store), F(SP), Uint64(3)}, // add end to stack
		{Op(Store), F(None), Int64(int64(n)), R(Reg(0))},
		{Op(JMP), F(Imm), Uint64(5)},
		{Op(PrintLn), F(Int), R(1)}, //end
		{Op(Exit)},
		{Op(JMPNEQ), F(ImmI), Uint64(8), Int64(0), R(0)},
		{Op(Store), F(None), Uint64(0), R(1)},
		{Op(JMP), F(SP)},
		{Op(JMPNEQ), F(ImmI), Uint64(11), Uint64(1), R(0)},
		{Op(Store), F(None), Uint64(1), R(1)},
		{Op(JMP), F(SP)},
		{Op(Store), F(SPR), R(0)},
		{Op(Sub), F(IImm), R(0), Uint64(1), R(0)},
		{Op(Store), F(SP), Uint64(15)},
		{Op(JMP), F(Imm), Uint64(5)}, // return value from call should be in asm.R1
		{Op(Load), F(SP), R(0)},
		{Op(Store), F(SPR), R(1)},
		{Op(Sub), F(IImm), R(0), Int64(2), R(0)},
		{Op(Store), F(SP), Uint64(20)},
		{Op(JMP), F(Imm), Uint64(5)},
		{Op(Load), F(SP), R(2)},
		{Op(Add), F(Int), R(1), R(2), R(1)},
		{Op(JMP), F(SP)},
	}

	for i := 0; i < b.N; i++ {
		buf := bytes.Buffer{}
		vm := NewVm(&buf)
		err := vm.Exec(bc)
		if err != nil {
			b.Error(err)
		}
		if string(buf.Bytes()[:]) != e {
			b.Errorf("expected %s got %s", e, string(buf.Bytes()[:]))
		}
	}

}
