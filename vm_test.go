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

func Test_And(t *testing.T) {
	test([]tcase{
		{
			// And Bool
			ByteCode{
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(Store), f(None), boolean(true), r(Reg(2))},
				{op(Store), f(None), boolean(false), r(Reg(4))},
				{op(Store), f(None), boolean(true), r(Reg(5))},
				{op(And), f(Bool), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(And), f(Bool), r(Reg(4)), r(Reg(5)), r(Reg(6))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// And ImmB
			ByteCode{
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(And), f(ImmB), boolean(true), r(Reg(1)), r(Reg(2))},
				{op(And), f(ImmB), boolean(false), r(Reg(1)), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(2))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Exit)},
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
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(Store), f(None), boolean(false), r(Reg(2))},
				{op(Store), f(None), boolean(false), r(Reg(4))},
				{op(Store), f(None), boolean(false), r(Reg(5))},
				{op(Or), f(Bool), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Or), f(Bool), r(Reg(4)), r(Reg(5)), r(Reg(6))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// Or ImmB
			ByteCode{
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(Or), f(ImmB), boolean(true), r(Reg(1)), r(Reg(2))},
				{op(Or), f(ImmB), boolean(false), r(Reg(1)), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(2))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Exit)},
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
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(Store), f(None), boolean(false), r(Reg(2))},
				{op(Not), f(Bool), r(Reg(1)), [8]byte{}, r(Reg(3))},
				{op(Not), f(Bool), r(Reg(2)), [8]byte{}, r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Add), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"30",
		},
		{
			// Add ImmI
			ByteCode{
				// store int 10 in reg 1,
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Add), f(ImmI), i64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"30",
		},
		{
			// Add Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(Add), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"30",
		},
		{
			// Add ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Add), f(ImmF), f64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(30), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Sub), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"20",
		},
		{
			// Sub ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Sub), f(ImmI), i64(30), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"20",
		},
		{
			// Sub IImm
			ByteCode{
				{op(Store), f(None), i64(30), r(Reg(1))},
				{op(Sub), f(IImm), r(Reg(1)), i64(10), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"20",
		},
		{
			// Sub Float
			ByteCode{
				{op(Store), f(None), f64(30), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(Sub), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"20",
		},
		{
			// Sub ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Sub), f(ImmF), f64(30), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
			},
			"20",
		},
		{
			// Sub FImm
			ByteCode{
				{op(Store), f(None), f64(30), r(Reg(1))},
				{op(Sub), f(FImm), r(Reg(1)), f64(10), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Mul), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"200",
		},
		{
			// Mul ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Mul), f(ImmI), i64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"200",
		},
		{
			// Mul Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(Mul), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"200",
		},
		{
			// Mul ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Mul), f(ImmF), f64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Quo), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Quo ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Quo), f(ImmI), i64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Quo IImm
			ByteCode{
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Quo), f(IImm), r(Reg(1)), i64(10), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Quo Float
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(Quo), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Quo ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Quo), f(ImmF), f64(20), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Quo FImm
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Quo), f(FImm), r(Reg(1)), f64(10), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(2), r(Reg(1))},
				{op(Store), f(None), i64(3), r(Reg(2))},
				{op(Pow), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"8",
		},
		{
			// Pow ImmI
			ByteCode{
				{op(Store), f(None), i64(3), r(Reg(1))},
				{op(Pow), f(ImmI), i64(2), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"8",
		},
		{
			// Pow IImm
			ByteCode{
				{op(Store), f(None), i64(2), r(Reg(1))},
				{op(Pow), f(IImm), r(Reg(1)), i64(3), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"8",
		},
		{
			// Pow Float
			ByteCode{
				{op(Store), f(None), f64(2), r(Reg(1))},
				{op(Store), f(None), f64(3), r(Reg(2))},
				{op(Pow), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"8",
		},
		{
			// Pow ImmF
			ByteCode{
				{op(Store), f(None), f64(3), r(Reg(1))},
				{op(Pow), f(ImmF), f64(2), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
			},
			"8",
		},
		{
			// Pow FImm
			ByteCode{
				{op(Store), f(None), f64(2), r(Reg(1))},
				{op(Pow), f(FImm), r(Reg(1)), f64(3), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(4), r(Reg(2))},
				{op(Rem), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Rem ImmI
			ByteCode{
				{op(Store), f(None), i64(4), r(Reg(1))},
				{op(Rem), f(ImmI), i64(10), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Rem IImm
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Rem), f(IImm), r(Reg(1)), i64(4), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Rem Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(4), r(Reg(2))},
				{op(Rem), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(Print), f(Float), r(Reg(3))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Rem ImmF
			ByteCode{
				{op(Store), f(None), f64(4), r(Reg(1))},
				{op(Rem), f(ImmF), f64(10), r(Reg(1)), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
			},
			"2",
		},
		{
			// Rem FImm
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Rem), f(FImm), r(Reg(1)), f64(4), r(Reg(2))},
				{op(Print), f(Float), r(Reg(2))},
				{op(Exit)},
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
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(20), r(Reg(2))},
				{op(Store), f(None), i64(20), r(Reg(3))},
				{op(Eq), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(Eq), f(Int), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(20), r(Reg(2))},
				{op(Eq), f(ImmI), i64(11), r(Reg(1)), r(Reg(3))},
				{op(Eq), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(Store), f(None), f64(20), r(Reg(3))},
				{op(Eq), f(Float), r(Reg(1)), r(Reg(3)), r(Reg(4))},
				{op(Eq), f(Float), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(Eq), f(ImmF), f64(10), r(Reg(2)), r(Reg(3))},
				{op(Eq), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq Bool
			ByteCode{
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(Store), f(None), boolean(true), r(Reg(2))},
				{op(Store), f(None), boolean(true), r(Reg(3))},
				{op(Eq), f(Bool), r(Reg(1)), r(Reg(3)), r(Reg(4))},
				{op(Eq), f(Bool), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Eq ImmB
			ByteCode{
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(Store), f(None), boolean(false), r(Reg(2))},
				{op(Eq), f(ImmB), boolean(true), r(Reg(2)), r(Reg(3))},
				{op(Eq), f(ImmB), boolean(false), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
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
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(20), r(Reg(2))},
				{op(Store), f(None), i64(20), r(Reg(3))},
				{op(NEq), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(NEq), f(Int), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(20), r(Reg(2))},
				{op(NEq), f(ImmI), i64(11), r(Reg(1)), r(Reg(3))},
				{op(NEq), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(Store), f(None), f64(20), r(Reg(3))},
				{op(NEq), f(Float), r(Reg(1)), r(Reg(3)), r(Reg(4))},
				{op(NEq), f(Float), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(20), r(Reg(2))},
				{op(NEq), f(ImmF), f64(10), r(Reg(2)), r(Reg(3))},
				{op(NEq), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq Bool
			ByteCode{
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(Store), f(None), boolean(true), r(Reg(2))},
				{op(Store), f(None), boolean(true), r(Reg(3))},
				{op(NEq), f(Bool), r(Reg(1)), r(Reg(3)), r(Reg(4))},
				{op(NEq), f(Bool), r(Reg(2)), r(Reg(3)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// NEq ImmB
			ByteCode{
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(Store), f(None), boolean(false), r(Reg(2))},
				{op(NEq), f(ImmB), boolean(true), r(Reg(2)), r(Reg(3))},
				{op(NEq), f(ImmB), boolean(false), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(LT), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(LT), f(Int), r(Reg(2)), r(Reg(1)), r(Reg(4))},
				// print reg 3
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(30), r(Reg(2))},
				{op(LT), f(ImmI), i64(20), r(Reg(1)), r(Reg(3))},
				{op(LT), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT Float
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(LT), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(LT), f(Float), r(Reg(2)), r(Reg(1)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// LT ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(30), r(Reg(2))},
				{op(LT), f(ImmF), f64(20), r(Reg(1)), r(Reg(3))},
				{op(LT), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Store), f(None), i64(10), r(Reg(3))},
				{op(LTE), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(LTE), f(Int), r(Reg(2)), r(Reg(1)), r(Reg(5))},
				{op(LTE), f(Int), r(Reg(3)), r(Reg(3)), r(Reg(6))},
				// print reg 3
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(30), r(Reg(2))},
				{op(LTE), f(ImmI), i64(20), r(Reg(1)), r(Reg(3))},
				{op(LTE), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(LTE), f(ImmI), i64(30), r(Reg(2)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE Float
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(Store), f(None), f64(10), r(Reg(3))},
				{op(LTE), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(LTE), f(Float), r(Reg(2)), r(Reg(1)), r(Reg(5))},
				{op(LTE), f(Float), r(Reg(2)), r(Reg(3)), r(Reg(6))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(30), r(Reg(2))},
				{op(LTE), f(ImmF), f64(20), r(Reg(1)), r(Reg(3))},
				{op(LTE), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(LTE), f(ImmF), f64(30), r(Reg(2)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(GT), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(GT), f(Int), r(Reg(2)), r(Reg(1)), r(Reg(4))},
				// print reg 3
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// GT ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(30), r(Reg(2))},
				{op(GT), f(ImmI), i64(20), r(Reg(1)), r(Reg(3))},
				{op(GT), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// GT Float
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(GT), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(3))},
				{op(GT), f(Float), r(Reg(2)), r(Reg(1)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
			},
			"truefalse",
		},
		{
			// GT ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(30), r(Reg(2))},
				{op(GT), f(ImmF), f64(20), r(Reg(1)), r(Reg(3))},
				{op(GT), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Exit)},
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
				{op(Store), f(None), i64(20), r(Reg(1))},
				{op(Store), f(None), i64(10), r(Reg(2))},
				{op(Store), f(None), i64(10), r(Reg(3))},
				{op(GTE), f(Int), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(GTE), f(Int), r(Reg(2)), r(Reg(1)), r(Reg(5))},
				{op(GTE), f(Int), r(Reg(3)), r(Reg(3)), r(Reg(6))},
				// print reg 3
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmI
			ByteCode{
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Store), f(None), i64(30), r(Reg(2))},
				{op(GTE), f(ImmI), i64(20), r(Reg(1)), r(Reg(3))},
				{op(GTE), f(ImmI), i64(20), r(Reg(2)), r(Reg(4))},
				{op(GTE), f(ImmI), i64(30), r(Reg(2)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE Float
			ByteCode{
				{op(Store), f(None), f64(20), r(Reg(1))},
				{op(Store), f(None), f64(10), r(Reg(2))},
				{op(Store), f(None), f64(10), r(Reg(3))},
				{op(GTE), f(Float), r(Reg(1)), r(Reg(2)), r(Reg(4))},
				{op(GTE), f(Float), r(Reg(2)), r(Reg(1)), r(Reg(5))},
				{op(GTE), f(Float), r(Reg(2)), r(Reg(3)), r(Reg(6))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(6))},
				{op(Exit)},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmF
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Store), f(None), f64(30), r(Reg(2))},
				{op(GTE), f(ImmF), f64(20), r(Reg(1)), r(Reg(3))},
				{op(GTE), f(ImmF), f64(20), r(Reg(2)), r(Reg(4))},
				{op(GTE), f(ImmF), f64(30), r(Reg(2)), r(Reg(5))},
				{op(Print), f(Bool), r(Reg(3))},
				{op(Print), f(Bool), r(Reg(4))},
				{op(Print), f(Bool), r(Reg(5))},
				{op(Exit)},
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
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(Print), f(Int), r(Reg(1))},
				{op(Exit)},
			},
			"10",
		},
		{
			// Print ImmI
			ByteCode{
				{op(Print), f(ImmI), i64(10)},
				{op(Exit)},
			},
			"10",
		},
		{
			// Print Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(Print), f(Float), r(Reg(1))},
				{op(Exit)},
			},
			"10",
		},
		{
			// Print ImmF
			ByteCode{
				{op(Print), f(ImmF), f64(10)},
				{op(Exit)},
			},
			"10",
		},
		{
			// Print Bool
			ByteCode{
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(Store), f(None), boolean(true), r(Reg(2))},
				{op(Print), f(Bool), r(Reg(1))},
				{op(Print), f(Bool), r(Reg(2))},
				{op(Exit)},
			},
			"falsetrue",
		},
		{
			// Print ImmB
			ByteCode{
				{op(Print), f(ImmB), boolean(false)},
				{op(Print), f(ImmB), boolean(true)},
				{op(Exit)},
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
				{op(Store), f(None), i64(10), r(Reg(1))},
				{op(PrintLn), f(Int), r(Reg(1))},
				{op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn ImmI
			ByteCode{
				{op(PrintLn), f(ImmI), i64(10)},
				{op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn Float
			ByteCode{
				{op(Store), f(None), f64(10), r(Reg(1))},
				{op(PrintLn), f(Float), r(Reg(1))},
				{op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn ImmF
			ByteCode{
				{op(PrintLn), f(ImmF), f64(10)},
				{op(Exit)},
			},
			"10\n",
		},
		{
			// PrintLn Bool
			ByteCode{
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(Store), f(None), boolean(true), r(Reg(2))},
				{op(PrintLn), f(Bool), r(Reg(1))},
				{op(PrintLn), f(Bool), r(Reg(2))},
				{op(Exit)},
			},
			"false\ntrue\n",
		},
		{
			// PrintLn ImmB
			ByteCode{
				{op(PrintLn), f(ImmB), boolean(false)},
				{op(PrintLn), f(ImmB), boolean(true)},
				{op(Exit)},
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
				{op(Store), f(None), r(Reg(9)), r(Reg(1))},
				{op(Store), f(None), i64(3), r(Reg(9))},
				{op(Load), f(Ptr), r(Reg(2)), r(Reg(1))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
			},
			"3",
		},
		{
			// Load from Stack
			ByteCode{
				{op(Store), f(SP), i64(3)},
				{op(Load), f(SP), r(Reg(1))},
				{op(Print), f(Int), r(Reg(1))},
				{op(Exit)},
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
				{op(Store), f(None), i64(1), r(Reg(1))},
				{op(Print), f(Int), r(Reg(1))},
				{op(Exit)},
			},
			"1",
		},
		{
			// Store bytes in stack from Imm
			ByteCode{
				{op(Store), f(SP), i64(1)},
				{op(Store), f(SP), i64(2)},
				{op(Store), f(SP), i64(3)},
				{op(Load), f(SP), r(Reg(1))},
				{op(Load), f(SP), r(Reg(2))},
				{op(Load), f(SP), r(Reg(3))},
				{op(Print), f(Int), r(Reg(1))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Print), f(Int), r(Reg(3))},
				{op(Exit)},
			},
			"321",
		},
		{
			// Store bytes in stack from Reg
			ByteCode{
				{op(Store), f(None), i64(1), r(Reg(1))},
				{op(Store), f(SPR), r(Reg(1))},
				{op(Load), f(SP), r(Reg(2))},
				{op(Print), f(Int), r(Reg(2))},
				{op(Exit)},
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
				{op(JMP), f(Imm), ui64(2)},
				{op(Print), f(ImmI), i64(1)},
				{op(Print), f(ImmI), i64(2)},
				{op(JMP), f(Imm), ui64(5)},
				{op(Print), f(ImmI), i64(1)},
				{op(Exit)},
			},
			"2",
		},
		{
			// JMP Reg
			ByteCode{
				{op(Store), f(None), ui64(4), r(Reg(1))},
				{op(Store), f(None), ui64(7), r(Reg(2))},
				{op(JMP), f(None), r(Reg(1))},
				{op(Print), f(ImmI), i64(1)},
				{op(Print), f(ImmI), i64(2)},
				{op(JMP), f(None), r(Reg(2))},
				{op(Print), f(ImmI), i64(1)},
				{op(Exit)},
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
				{op(Store), f(None), boolean(true), r(Reg(1))},
				{op(JMPEQ), f(ImmB), ui64(3), boolean(true), r(Reg(1))},
				{op(Print), f(ImmI), i64(1)},
				{op(Print), f(ImmI), i64(2)},
				{op(Store), f(None), boolean(false), r(Reg(1))},
				{op(JMPEQ), f(ImmB), ui64(7), boolean(false), r(Reg(1))},
				{op(Print), f(ImmI), i64(1)},
				{op(Exit)},
			},
			"2",
		},
		{
			// ImmI
			ByteCode{
				{op(Store), f(None), i64(2), r(Reg(1))},
				{op(JMPEQ), f(ImmI), ui64(3), i64(2), r(Reg(1))},
				{op(Print), f(ImmI), i64(1)},
				{op(Print), f(ImmI), i64(2)},
				{op(Store), f(None), i64(3), r(Reg(1))},
				{op(JMPEQ), f(ImmI), ui64(7), ui64(3), r(Reg(1))},
				{op(Print), f(ImmI), i64(1)},
				{op(Exit)},
			},
			"2",
		},
	}, t)
}

func Test_UnhandledOp(t *testing.T) {
	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	err := vm.Exec(ByteCode{{ui64(1999)}})
	if err == nil {
		t.Error("expected error, got none")
	}
}
