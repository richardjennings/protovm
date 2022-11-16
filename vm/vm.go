package vm

import (
	"fmt"
	"io"
	"math"
	"unsafe"
)

type (
	VM struct {
		w  io.Writer
		pc uint64
		sp uint64
		r  Registers
		s  Stack
		h  []interface{}
	}
)

// NewVm Creates a new VM struct
func NewVm(writer io.Writer) *VM {
	v := &VM{w: writer}
	return v
}

// Exec Execute instructions
func (vm *VM) Exec(bc ByteCode) error {
	_ = vm.r[len(vm.r)-1] // bounds check elimination
	_ = vm.s[len(vm.s)-1] // bounds check elimination
	// set sp to top of stack
	vm.sp = uint64(len(vm.s) - 1)
	var i *Inst
	var o uint64

	for {
		i = &bc[vm.pc]
		o = *(*uint64)(unsafe.Pointer(&i[0]))
		switch o {
		case uint64(And):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) &&
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) &&
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(Or):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ||
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) ||
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(Not):
			*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
				!*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))

		case uint64(Add):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) +
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) +
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) +
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) +
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(Sub):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) -
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) -
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*float64)(unsafe.Pointer(&i[3]))
			}

		case uint64(Mul):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) *
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) *
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) *
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) *
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(Quo):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) /
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) /
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*float64)(unsafe.Pointer(&i[3]))
			}

		case uint64(Pow):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&i[2]))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))),
						float64(*(*int64)(unsafe.Pointer(&i[3])))))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&i[2])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&i[3])))
			}

		case uint64(Rem):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) %
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) %
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) %
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&i[2])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&i[3])))
			}

		case uint64(Eq):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) ==
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) ==
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) ==
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case uint64(NEq):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) !=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) !=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) !=
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(LT):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) <
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) <
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(LTE):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) <=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) <=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(GT):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) >
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) >
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(GTE):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) >=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) >=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}

		case uint64(Print):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmI:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&i[2])))
			case Float:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmF:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&i[2])))
			case Bool:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmB:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&i[2])))
			}

		case uint64(PrintLn):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmI:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&i[2])))
			case Float:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmF:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&i[2])))
			case Bool:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmB:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&i[2])))
			}

		case uint64(Load):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case SP:
				vm.sp++
				vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.s[vm.sp]
			}

		case uint64(Store):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case None:
				vm.r[*(*uint64)(unsafe.Pointer(&i[3]))] = i[2]
			case SP:
				vm.s[vm.sp] = i[2]
				vm.sp--
			case SPR:
				vm.s[vm.sp] = vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]
				vm.sp--
			}

		case uint64(JMP):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case None:
				vm.pc = *(*uint64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))
				continue
			case Imm:
				vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
				continue
			case SP:
				vm.sp++
				vm.pc = *(*uint64)(unsafe.Pointer(&vm.s[vm.sp]))
				continue
			}

		case uint64(JMPEQ):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case ImmB:
				if *(*bool)(unsafe.Pointer(&i[3])) == *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) {
					vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
					continue
				}
			case ImmI:
				if *(*int64)(unsafe.Pointer(&i[3])) == *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) {
					vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
					continue
				}
			}

		case uint64(JMPNEQ):
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case ImmB:
				if *(*bool)(unsafe.Pointer(&i[3])) != *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) {
					vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
					continue
				}
			case ImmI:
				if *(*int64)(unsafe.Pointer(&i[3])) != *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) {
					vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
					continue
				}
			}

		case uint64(Exit):
			return nil

		case uint64(NoOp):
			vm.pc++
			continue

		default:
			return fmt.Errorf("unhandled OP %d", o)
		}
		vm.pc++
	}
}
