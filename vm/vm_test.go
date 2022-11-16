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
				{O: Exit},
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
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(2))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(4))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(5))},
				{O: And, F: Bool, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: And, F: Bool, X: R(Reg(4)), Y: R(Reg(5)), Z: R(Reg(6))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// And ImmB
			ByteCode{
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: And, F: ImmB, X: Boolean(true), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: And, F: ImmB, X: Boolean(false), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(2))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Exit},
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
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(2))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(4))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(5))},
				{O: Or, F: Bool, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Or, F: Bool, X: R(Reg(4)), Y: R(Reg(5)), Z: R(Reg(6))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// Or ImmB
			ByteCode{
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: Or, F: ImmB, X: Boolean(true), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Or, F: ImmB, X: Boolean(false), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(2))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Exit},
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
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(2))},
				{O: Not, F: Bool, X: R(Reg(1)), Y: 0, Z: R(Reg(3))},
				{O: Not, F: Bool, X: R(Reg(2)), Y: 0, Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
	}, t)
}

func Test_Band(t *testing.T) {
	test([]tcase{
		{
			// Not
			ByteCode{
				{O: Store, F: None, X: 0b1011, Y: R(Reg(1))},
				{O: Store, F: None, X: 0b0011, Y: R(Reg(2))},
				{O: Band, F: None, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: None, X: R(Reg(1))},
				{O: Print, F: None, X: R(Reg(2))},
				{O: Print, F: None, X: R(Reg(3))},
				{O: Exit},
			},
			"[11 0 0 0 0 0 0 0][3 0 0 0 0 0 0 0][3 0 0 0 0 0 0 0]",
		},
	}, t)
}

func Test_Bor(t *testing.T) {
	test([]tcase{
		{
			// Not
			ByteCode{
				{O: Store, F: None, X: 0b1011, Y: R(Reg(1))},
				{O: Store, F: None, X: 0b0011, Y: R(Reg(2))},
				{O: Bor, F: None, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: None, X: R(Reg(1))},
				{O: Print, F: None, X: R(Reg(2))},
				{O: Print, F: None, X: R(Reg(3))},
				{O: Exit},
			},
			"[11 0 0 0 0 0 0 0][3 0 0 0 0 0 0 0][11 0 0 0 0 0 0 0]",
		},
	}, t)
}

func Test_Bnot(t *testing.T) {
	test([]tcase{
		{
			// Not
			ByteCode{
				{O: Store, F: None, X: 0b1011, Y: R(Reg(1))},
				{O: Bnot, F: None, X: R(Reg(1)), Z: R(Reg(3))},
				{O: Print, F: None, X: R(Reg(1))},
				{O: Print, F: None, X: R(Reg(3))},
				{O: Exit},
			},
			"[11 0 0 0 0 0 0 0][244 255 255 255 255 255 255 255]",
		},
	}, t)
}

func Test_Bxor(t *testing.T) {
	test([]tcase{
		{
			// Not
			ByteCode{
				{O: Store, F: None, X: 0b1011, Y: R(Reg(1))},
				{O: Store, F: None, X: 0b0011, Y: R(Reg(2))},
				{O: Bxor, F: None, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: None, X: R(Reg(1))},
				{O: Print, F: None, X: R(Reg(2))},
				{O: Print, F: None, X: R(Reg(3))},
				{O: Exit},
			},
			"[11 0 0 0 0 0 0 0][3 0 0 0 0 0 0 0][8 0 0 0 0 0 0 0]",
		},
	}, t)
}

func Test_Add(t *testing.T) {
	test([]tcase{
		{
			// Add Int
			ByteCode{
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Add, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"30",
		},
		{
			// Add ImmI
			ByteCode{
				// store int 10 in reg 1,
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Add, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"30",
		},
		{
			// Add Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: Add, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"30",
		},
		{
			// Add ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Add, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(30), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Sub, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"20",
		},
		{
			// Sub ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Sub, F: ImmI, X: Int64(30), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"20",
		},
		{
			// Sub IImm
			ByteCode{
				{O: Store, F: None, X: Int64(30), Y: R(Reg(1))},
				{O: Sub, F: IImm, X: R(Reg(1)), Y: Int64(10), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"20",
		},
		{
			// Sub Float
			ByteCode{
				{O: Store, F: None, X: Float64(30), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: Sub, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"20",
		},
		{
			// Sub ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Sub, F: ImmF, X: Float64(30), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
			},
			"20",
		},
		{
			// Sub FImm
			ByteCode{
				{O: Store, F: None, X: Float64(30), Y: R(Reg(1))},
				{O: Sub, F: FImm, X: R(Reg(1)), Y: Float64(10), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Mul, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"200",
		},
		{
			// Mul ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Mul, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"200",
		},
		{
			// Mul Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: Mul, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"200",
		},
		{
			// Mul ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Mul, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Quo, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"2",
		},
		{
			// Quo ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Quo, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// Quo IImm
			ByteCode{
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Quo, F: IImm, X: R(Reg(1)), Y: Int64(10), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// Quo Float
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: Quo, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"2",
		},
		{
			// Quo ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Quo, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// Quo FImm
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Quo, F: FImm, X: R(Reg(1)), Y: Float64(10), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(2), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(3), Y: R(Reg(2))},
				{O: Pow, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"8",
		},
		{
			// Pow ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(3), Y: R(Reg(1))},
				{O: Pow, F: ImmI, X: Int64(2), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"8",
		},
		{
			// Pow IImm
			ByteCode{
				{O: Store, F: None, X: Int64(2), Y: R(Reg(1))},
				{O: Pow, F: IImm, X: R(Reg(1)), Y: Int64(3), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"8",
		},
		{
			// Pow Float
			ByteCode{
				{O: Store, F: None, X: Float64(2), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(3), Y: R(Reg(2))},
				{O: Pow, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"8",
		},
		{
			// Pow ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(3), Y: R(Reg(1))},
				{O: Pow, F: ImmF, X: Float64(2), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
			},
			"8",
		},
		{
			// Pow FImm
			ByteCode{
				{O: Store, F: None, X: Float64(2), Y: R(Reg(1))},
				{O: Pow, F: FImm, X: R(Reg(1)), Y: Float64(3), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(4), Y: R(Reg(2))},
				{O: Rem, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"2",
		},
		{
			// asm.Rem ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(4), Y: R(Reg(1))},
				{O: Rem, F: ImmI, X: Int64(10), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// asm.Rem IImm
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Rem, F: IImm, X: R(Reg(1)), Y: Int64(4), Z: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// asm.Rem Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(4), Y: R(Reg(2))},
				{O: Rem, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Print, F: Float, X: R(Reg(3))},
				{O: Exit},
			},
			"2",
		},
		{
			// asm.Rem ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(4), Y: R(Reg(1))},
				{O: Rem, F: ImmF, X: Float64(10), Y: R(Reg(1)), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
			},
			"2",
		},
		{
			// asm.Rem FImm
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Rem, F: FImm, X: R(Reg(1)), Y: Float64(4), Z: R(Reg(2))},
				{O: Print, F: Float, X: R(Reg(2))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(2))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(3))},
				{O: Eq, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Eq, F: Int, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Eq ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(2))},
				{O: Eq, F: ImmI, X: Int64(11), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: Eq, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Eq Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(3))},
				{O: Eq, F: Float, X: R(Reg(1)), Y: R(Reg(3)), Z: R(Reg(4))},
				{O: Eq, F: Float, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Eq ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: Eq, F: ImmF, X: Float64(10), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Eq, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Eq Bool
			ByteCode{
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(2))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(3))},
				{O: Eq, F: Bool, X: R(Reg(1)), Y: R(Reg(3)), Z: R(Reg(4))},
				{O: Eq, F: Bool, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Eq ImmB
			ByteCode{
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(2))},
				{O: Eq, F: ImmB, X: Boolean(true), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: Eq, F: ImmB, X: Boolean(false), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(2))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(3))},
				{O: NEq, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: NEq, F: Int, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// NEq ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(20), Y: R(Reg(2))},
				{O: NEq, F: ImmI, X: Int64(11), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: NEq, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// NEq Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(3))},
				{O: NEq, F: Float, X: R(Reg(1)), Y: R(Reg(3)), Z: R(Reg(4))},
				{O: NEq, F: Float, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// NEq ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(20), Y: R(Reg(2))},
				{O: NEq, F: ImmF, X: Float64(10), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: NEq, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// NEq Bool
			ByteCode{
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(2))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(3))},
				{O: NEq, F: Bool, X: R(Reg(1)), Y: R(Reg(3)), Z: R(Reg(4))},
				{O: NEq, F: Bool, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// NEq ImmB
			ByteCode{
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(2))},
				{O: NEq, F: ImmB, X: Boolean(true), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: NEq, F: ImmB, X: Boolean(false), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: LT, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: LT, F: Int, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(4))},
				// print reg 3
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// LT ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(30), Y: R(Reg(2))},
				{O: LT, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: LT, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// LT Float
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: LT, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: LT, F: Float, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// LT ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(30), Y: R(Reg(2))},
				{O: LT, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: LT, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(3))},
				{O: LTE, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: LTE, F: Int, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(5))},
				{O: LTE, F: Int, X: R(Reg(3)), Y: R(Reg(3)), Z: R(Reg(6))},
				// print reg 3
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(30), Y: R(Reg(2))},
				{O: LTE, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: LTE, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: LTE, F: ImmI, X: Int64(30), Y: R(Reg(2)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"falsetruetrue",
		},
		{
			// LTE Float
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(3))},
				{O: LTE, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: LTE, F: Float, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(5))},
				{O: LTE, F: Float, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(6))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"falsetruetrue",
		},
		{
			// LTE ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(30), Y: R(Reg(2))},
				{O: LTE, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: LTE, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: LTE, F: ImmF, X: Float64(30), Y: R(Reg(2)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: GT, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: GT, F: Int, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(4))},
				// print reg 3
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// GT ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(30), Y: R(Reg(2))},
				{O: GT, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: GT, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// GT Float
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: GT, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(3))},
				{O: GT, F: Float, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
			},
			"truefalse",
		},
		{
			// GT ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(30), Y: R(Reg(2))},
				{O: GT, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: GT, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(2))},
				{O: Store, F: None, X: Int64(10), Y: R(Reg(3))},
				{O: GTE, F: Int, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: GTE, F: Int, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(5))},
				{O: GTE, F: Int, X: R(Reg(3)), Y: R(Reg(3)), Z: R(Reg(6))},
				// print reg 3
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Int64(30), Y: R(Reg(2))},
				{O: GTE, F: ImmI, X: Int64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: GTE, F: ImmI, X: Int64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: GTE, F: ImmI, X: Int64(30), Y: R(Reg(2)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
			},
			"truefalsetrue",
		},
		{
			// GTE Float
			ByteCode{
				{O: Store, F: None, X: Float64(20), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(2))},
				{O: Store, F: None, X: Float64(10), Y: R(Reg(3))},
				{O: GTE, F: Float, X: R(Reg(1)), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: GTE, F: Float, X: R(Reg(2)), Y: R(Reg(1)), Z: R(Reg(5))},
				{O: GTE, F: Float, X: R(Reg(2)), Y: R(Reg(3)), Z: R(Reg(6))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(6))},
				{O: Exit},
			},
			"truefalsetrue",
		},
		{
			// GTE ImmF
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Store, F: None, X: Float64(30), Y: R(Reg(2))},
				{O: GTE, F: ImmF, X: Float64(20), Y: R(Reg(1)), Z: R(Reg(3))},
				{O: GTE, F: ImmF, X: Float64(20), Y: R(Reg(2)), Z: R(Reg(4))},
				{O: GTE, F: ImmF, X: Float64(30), Y: R(Reg(2)), Z: R(Reg(5))},
				{O: Print, F: Bool, X: R(Reg(3))},
				{O: Print, F: Bool, X: R(Reg(4))},
				{O: Print, F: Bool, X: R(Reg(5))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: Print, F: Int, X: R(Reg(1))},
				{O: Exit},
			},
			"10",
		},
		{
			// Print ImmI
			ByteCode{
				{O: Print, F: ImmI, X: Int64(10)},
				{O: Exit},
			},
			"10",
		},
		{
			// Print Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: Print, F: Float, X: R(Reg(1))},
				{O: Exit},
			},
			"10",
		},
		{
			// Print ImmF
			ByteCode{
				{O: Print, F: ImmF, X: Float64(10)},
				{O: Exit},
			},
			"10",
		},
		{
			// Print Bool
			ByteCode{
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(2))},
				{O: Print, F: Bool, X: R(Reg(1))},
				{O: Print, F: Bool, X: R(Reg(2))},
				{O: Exit},
			},
			"falsetrue",
		},
		{
			// Print ImmB
			ByteCode{
				{O: Print, F: ImmB, X: Boolean(false)},
				{O: Print, F: ImmB, X: Boolean(true)},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(10), Y: R(Reg(1))},
				{O: PrintLn, F: Int, X: R(Reg(1))},
				{O: Exit},
			},
			"10\n",
		},
		{
			// PrintLn ImmI
			ByteCode{
				{O: PrintLn, F: ImmI, X: Int64(10)},
				{O: Exit},
			},
			"10\n",
		},
		{
			// PrintLn Float
			ByteCode{
				{O: Store, F: None, X: Float64(10), Y: R(Reg(1))},
				{O: PrintLn, F: Float, X: R(Reg(1))},
				{O: Exit},
			},
			"10\n",
		},
		{
			// PrintLn ImmF
			ByteCode{
				{O: PrintLn, F: ImmF, X: Float64(10)},
				{O: Exit},
			},
			"10\n",
		},
		{
			// PrintLn Bool
			ByteCode{
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(2))},
				{O: PrintLn, F: Bool, X: R(Reg(1))},
				{O: PrintLn, F: Bool, X: R(Reg(2))},
				{O: Exit},
			},
			"false\ntrue\n",
		},
		{
			// PrintLn ImmB
			ByteCode{
				{O: PrintLn, F: ImmB, X: Boolean(false)},
				{O: PrintLn, F: ImmB, X: Boolean(true)},
				{O: Exit},
			},
			"false\ntrue\n",
		},
	}, t)
}

func Test_Load(t *testing.T) {
	test([]tcase{
		{
			// Load from Stack
			ByteCode{
				{O: Store, F: SP, X: Int64(3)},
				{O: Load, F: SP, X: R(Reg(1))},
				{O: Print, F: Int, X: R(Reg(1))},
				{O: Exit},
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
				{O: Store, F: None, X: Int64(1), Y: R(Reg(1))},
				{O: Print, F: Int, X: R(Reg(1))},
				{O: Exit},
			},
			"1",
		},
		{
			// Store bytes in stack from Imm
			ByteCode{
				{O: Store, F: SP, X: Int64(1)},
				{O: Store, F: SP, X: Int64(2)},
				{O: Store, F: SP, X: Int64(3)},
				{O: Load, F: SP, X: R(Reg(1))},
				{O: Load, F: SP, X: R(Reg(2))},
				{O: Load, F: SP, X: R(Reg(3))},
				{O: Print, F: Int, X: R(Reg(1))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(3))},
				{O: Exit},
			},
			"321",
		},
		{
			// Store bytes in stack from asm.Reg
			ByteCode{
				{O: Store, F: None, X: Int64(1), Y: R(Reg(1))},
				{O: Store, F: SPR, X: R(Reg(1))},
				{O: Load, F: SP, X: R(Reg(2))},
				{O: Print, F: Int, X: R(Reg(2))},
				{O: Exit},
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
				{O: JMP, F: Imm, X: 2},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: JMP, F: Imm, X: 5},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
			},
			"2",
		},
		{
			// JMP asm.Reg
			ByteCode{
				{O: Store, F: None, X: 4, Y: R(Reg(1))},
				{O: Store, F: None, X: 7, Y: R(Reg(2))},
				{O: JMP, F: None, X: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: JMP, F: None, X: R(Reg(2))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
			},
			"2",
		},
		{
			// JMP Stack Pointer
			ByteCode{
				{O: Store, F: SP, X: 7},
				{O: Store, F: SP, X: 4},
				{O: JMP, F: SP},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: JMP, F: SP},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
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
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: JMPEQ, F: ImmB, X: 3, Y: Boolean(true), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: JMPEQ, F: ImmB, X: 7, Y: Boolean(false), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
			},
			"2",
		},
		{
			// ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(2), Y: R(Reg(1))},
				{O: JMPEQ, F: ImmI, X: 3, Y: Int64(2), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: Store, F: None, X: Int64(3), Y: R(Reg(1))},
				{O: JMPEQ, F: ImmI, X: 7, Y: 3, Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
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
				{O: Store, F: None, X: Boolean(false), Y: R(Reg(1))},
				{O: JMPNEQ, F: ImmB, X: 3, Y: Boolean(true), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: Store, F: None, X: Boolean(true), Y: R(Reg(1))},
				{O: JMPNEQ, F: ImmB, X: 7, Y: Boolean(false), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
			},
			"2",
		},
		{
			// ImmI
			ByteCode{
				{O: Store, F: None, X: Int64(1), Y: R(Reg(1))},
				{O: JMPNEQ, F: ImmI, X: 3, Y: Int64(2), Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Print, F: ImmI, X: Int64(2)},
				{O: Store, F: None, X: Int64(4), Y: R(Reg(1))},
				{O: JMPNEQ, F: ImmI, X: 7, Y: 3, Z: R(Reg(1))},
				{O: Print, F: ImmI, X: Int64(1)},
				{O: Exit},
			},
			"2",
		},
	}, t)
}

func Test_UnhandledOp(t *testing.T) {
	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	err := vm.Exec(ByteCode{{O: 127}})
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
		{O: Store, F: SP, X: 3}, // add end to stack
		{O: Store, F: None, X: Int64(int64(n)), Y: R(Reg(0))},
		{O: JMP, F: Imm, X: 5},
		{O: PrintLn, F: Int, X: R(1)}, //end
		{O: Exit},
		{O: JMPNEQ, F: ImmI, X: 8, Y: Int64(0), Z: R(0)},
		{O: Store, F: None, X: 0, Y: R(1)},
		{O: JMP, F: SP},
		{O: JMPNEQ, F: ImmI, X: 11, Y: 1, Z: R(0)},
		{O: Store, F: None, X: 1, Y: R(1)},
		{O: JMP, F: SP},
		{O: Store, F: SPR, X: R(0)},
		{O: Sub, F: IImm, X: R(0), Y: 1, Z: R(0)},
		{O: Store, F: SP, X: 15},
		{O: JMP, F: Imm, X: 5}, // return value from call should be in asm.R1
		{O: Load, F: SP, X: R(0)},
		{O: Store, F: SPR, X: R(1)},
		{O: Sub, F: IImm, X: R(0), Y: Int64(2), Z: R(0)},
		{O: Store, F: SP, X: 20},
		{O: JMP, F: Imm, X: 5},
		{O: Load, F: SP, X: R(2)},
		{O: Add, F: Int, X: R(1), Y: R(2), Z: R(1)},
		{O: JMP, F: SP},
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
