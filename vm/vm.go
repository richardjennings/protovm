package vm

import (
	"errors"
	"fmt"
	"io"
	"math"
	"unsafe"
)

type (
	ByteCode  []Inst
	Inst      [5][8]byte
	Registers [20][8]byte
	VM        struct {
		w  io.Writer
		pc uint64
		sp uint64
		r  Registers
		s  [100][8]byte
		h  []interface{}
	}
)

// Create a new VM struct
func NewVm(writer io.Writer) *VM {
	v := &VM{w: writer}
	return v
}

// Execute instructions
func (vm *VM) Exec(bc ByteCode) error {
	_ = vm.r[len(vm.r)-1] // bounds check elimination
	_ = vm.s[len(vm.s)-1] // bounds check elimination
	var i *Inst
	var o uint64

	for {
		i = &bc[vm.pc]
		o = *(*uint64)(unsafe.Pointer(&i[0]))
		if o > 11 {
			if o > 17 {
				if o > 20 {
					if o == uint64(JMPEQ) {
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
					} else if o == uint64(JMPNEQ) {
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
					} else if o == uint64(Exit) {
						return nil
					} else {
						return errors.New("unhandled Op")
					}
				} else {
					if o == uint64(Load) {
						switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
						case Ptr:
							vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.r[*(*uint64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))]
						case SP:
							vm.sp--
							vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.s[vm.sp]
						}
					} else if o == uint64(Store) {
						switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
						case None:
							vm.r[*(*uint64)(unsafe.Pointer(&i[3]))] = i[2]
						case SP:
							vm.s[vm.sp] = i[2] //@todo store register
							vm.sp++
						case SPR:
							vm.s[vm.sp] = vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]
							vm.sp++
						}
					} else if o == uint64(JMP) {
						switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
						case None:
							vm.pc = *(*uint64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))
							continue
						case Imm:
							vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
							continue
						case SP:
							vm.sp--
							vm.pc = *(*uint64)(unsafe.Pointer(&vm.s[vm.sp]))
							continue
						}
					}
				}
			} else {
				if o > 14 {
					if o == uint64(GTE) {
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
					} else if o == uint64(Print) {
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
					} else if o == uint64(PrintLn) {
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
					}
				} else {
					if o == uint64(LT) {
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
					} else if o == uint64(LTE) {
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
					} else if o == uint64(GT) {
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
					}
				}
			}
		} else {
			if o > 5 {
				if o > 8 {
					if o == uint64(Rem) {
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
					} else if o == uint64(Eq) {
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
					} else if o == uint64(NEq) {
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
					}
				} else {
					if o == uint64(Mul) {
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
					} else if o == uint64(Quo) {
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
					} else if o == uint64(Pow) {
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
					}
				}
			} else {
				if o > 2 {
					if o == uint64(Not) {
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
							!*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))
					} else if o == uint64(Add) {
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
					} else if o == uint64(Sub) {
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
					}
				} else {
					if o == uint64(NoOp) {
						vm.pc++
						continue
					} else if o == uint64(And) {
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
					} else if o == uint64(Or) {
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
					}
				}
			}
		}
		vm.pc++
	}
}
