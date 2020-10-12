package protovm

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
	vm := NewVm(&buf, 10, 1000)
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
			// And ImmediateBool
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(And), F(ImmediateBool), Boolean(true), R(Reg(1)), R(Reg(2))},
				{Op(And), F(ImmediateBool), Boolean(false), R(Reg(1)), R(Reg(3))},
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
			// Or ImmediateBool
			ByteCode{
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(Or), F(ImmediateBool), Boolean(true), R(Reg(1)), R(Reg(2))},
				{Op(Or), F(ImmediateBool), Boolean(false), R(Reg(1)), R(Reg(3))},
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
			// Add ImmediateInt
			ByteCode{
				// store int 10 in reg 1,
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Add), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(2))},
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
			// Add ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Add), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(2))},
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
			// Sub ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Sub), F(ImmediateInt), Int64(30), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub IntImmediate
			ByteCode{
				{Op(Store), F(None), Int64(30), R(Reg(1))},
				{Op(Sub), F(IntImmediate), R(Reg(1)), Int64(10), R(Reg(2))},
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
			// Sub ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Sub), F(ImmediateFloat), Float64(30), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"20",
		},
		{
			// Sub FloatImmediate
			ByteCode{
				{Op(Store), F(None), Float64(30), R(Reg(1))},
				{Op(Sub), F(FloatImmediate), R(Reg(1)), Float64(10), R(Reg(2))},
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
			// Mul ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Mul), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(2))},
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
			// Mul ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Mul), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(2))},
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
			// Quo ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Quo), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo IntImmediate
			ByteCode{
				{Op(Store), F(None), Int64(20), R(Reg(1))},
				{Op(Quo), F(IntImmediate), R(Reg(1)), Int64(10), R(Reg(2))},
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
			// Quo ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Quo), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Quo FloatImmediate
			ByteCode{
				{Op(Store), F(None), Float64(20), R(Reg(1))},
				{Op(Quo), F(FloatImmediate), R(Reg(1)), Float64(10), R(Reg(2))},
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
			// Pow ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(3), R(Reg(1))},
				{Op(Pow), F(ImmediateInt), Int64(2), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow IntImmediate
			ByteCode{
				{Op(Store), F(None), Int64(2), R(Reg(1))},
				{Op(Pow), F(IntImmediate), R(Reg(1)), Int64(3), R(Reg(2))},
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
			// Pow ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(3), R(Reg(1))},
				{Op(Pow), F(ImmediateFloat), Float64(2), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"8",
		},
		{
			// Pow FloatImmediate
			ByteCode{
				{Op(Store), F(None), Float64(2), R(Reg(1))},
				{Op(Pow), F(FloatImmediate), R(Reg(1)), Float64(3), R(Reg(2))},
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
			// Rem Int
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
			// Rem ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(4), R(Reg(1))},
				{Op(Rem), F(ImmediateInt), Int64(10), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Rem IntImmediate
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Rem), F(IntImmediate), R(Reg(1)), Int64(4), R(Reg(2))},
				{Op(Print), F(Int), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Rem Float
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
			// Rem ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(4), R(Reg(1))},
				{Op(Rem), F(ImmediateFloat), Float64(10), R(Reg(1)), R(Reg(2))},
				{Op(Print), F(Float), R(Reg(2))},
				{Op(Exit)},
			},
			"2",
		},
		{
			// Rem FloatImmediate
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Rem), F(FloatImmediate), R(Reg(1)), Float64(4), R(Reg(2))},
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
			// Eq ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(Eq), F(ImmediateInt), Int64(11), R(Reg(1)), R(Reg(3))},
				{Op(Eq), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
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
			// Eq ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(Eq), F(ImmediateFloat), Float64(10), R(Reg(2)), R(Reg(3))},
				{Op(Eq), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
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
			// Eq ImmediateBool
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(Eq), F(ImmediateBool), Boolean(true), R(Reg(2)), R(Reg(3))},
				{Op(Eq), F(ImmediateBool), Boolean(false), R(Reg(2)), R(Reg(4))},
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
			// NEq ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(20), R(Reg(2))},
				{Op(NEq), F(ImmediateInt), Int64(11), R(Reg(1)), R(Reg(3))},
				{Op(NEq), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
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
			// NEq ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(20), R(Reg(2))},
				{Op(NEq), F(ImmediateFloat), Float64(10), R(Reg(2)), R(Reg(3))},
				{Op(NEq), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
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
			// NEq ImmediateBool
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(Store), F(None), Boolean(false), R(Reg(2))},
				{Op(NEq), F(ImmediateBool), Boolean(true), R(Reg(2)), R(Reg(3))},
				{Op(NEq), F(ImmediateBool), Boolean(false), R(Reg(2)), R(Reg(4))},
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
			// LT ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(LT), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(LT), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
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
			// LT ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(LT), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(LT), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
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
			// LTE ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(LTE), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(LTE), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(ImmediateInt), Int64(30), R(Reg(2)), R(Reg(5))},
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
			// LTE ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(LTE), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(LTE), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(LTE), F(ImmediateFloat), Float64(30), R(Reg(2)), R(Reg(5))},
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
			// GT ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(GT), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(GT), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
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
			// GT ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10),  R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(GT), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(GT), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
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
			// GTE ImmediateInt
			ByteCode{
				{Op(Store), F(None), Int64(10), R(Reg(1))},
				{Op(Store), F(None), Int64(30), R(Reg(2))},
				{Op(GTE), F(ImmediateInt), Int64(20), R(Reg(1)), R(Reg(3))},
				{Op(GTE), F(ImmediateInt), Int64(20), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(ImmediateInt), Int64(30), R(Reg(2)), R(Reg(5))},
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
			// GTE ImmediateFloat
			ByteCode{
				{Op(Store), F(None), Float64(10), R(Reg(1))},
				{Op(Store), F(None), Float64(30), R(Reg(2))},
				{Op(GTE), F(ImmediateFloat), Float64(20), R(Reg(1)), R(Reg(3))},
				{Op(GTE), F(ImmediateFloat), Float64(20), R(Reg(2)), R(Reg(4))},
				{Op(GTE), F(ImmediateFloat), Float64(30), R(Reg(2)), R(Reg(5))},
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
			// Print ImmediateInt
			ByteCode{
				{Op(Print), F(ImmediateInt), Int64(10)},
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
			// Print ImmediateFloat
			ByteCode{
				{Op(Print), F(ImmediateFloat), Float64(10)},
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
			// Print ImmediateBool
			ByteCode{
				{Op(Print), F(ImmediateBool), Boolean(false)},
				{Op(Print), F(ImmediateBool), Boolean(true)},
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
			// PrintLn ImmediateInt
			ByteCode{
				{Op(PrintLn), F(ImmediateInt), Int64(10)},
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
			// PrintLn ImmediateFloat
			ByteCode{
				{Op(PrintLn), F(ImmediateFloat), Float64(10)},
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
			// PrintLn ImmediateBool
			ByteCode{
				{Op(PrintLn), F(ImmediateBool), Boolean(false)},
				{Op(PrintLn), F(ImmediateBool), Boolean(true)},
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
			// Store bytes in stack from Immediate
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
			// Store bytes in stack from Reg
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
			// JMP Immediate
			ByteCode{
				{Op(JMP), F(Immediate), Uint64(2)},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Print), F(ImmediateInt), Int64(2)},
				{Op(JMP), F(Immediate), Uint64(5)},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
		{
			// JMP Reg
			ByteCode{
				{Op(Store), F(None), Uint64(4), R(Reg(1))},
				{Op(Store), F(None), Uint64(7), R(Reg(2))},
				{Op(JMP), F(None), R(Reg(1))},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Print), F(ImmediateInt), Int64(2)},
				{Op(JMP), F(None), R(Reg(2))},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_JMPEQ(t *testing.T) {
	test([]tcase{
		{
			// print reg 1
			ByteCode{
				{Op(Store), F(None), Boolean(true), R(Reg(1))},
				{Op(JMPEQ), F(ImmediateBool), Uint64(3), Boolean(true), R(Reg(1))},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Print), F(ImmediateInt), Int64(2)},
				{Op(Store), F(None), Boolean(false), R(Reg(1))},
				{Op(JMPEQ), F(ImmediateBool), Uint64(7), Boolean(false), R(Reg(1))},
				{Op(Print), F(ImmediateInt), Int64(1)},
				{Op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_UnhandledOp(t *testing.T) {
	buf := bytes.Buffer{}
	vm := NewVm(&buf, 2, 2)
	err := vm.Exec(ByteCode{{Uint64(1999)}})
	if err == nil {
		t.Error("expected error, got none")
	}
}