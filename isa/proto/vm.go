package proto

import (
	"fmt"
	"io"
	"math"
	"unsafe"
)

type (
	VM struct {
		w io.Writer
		r Registers
		m RAM
	}
)

// NewVm Creates a new VM struct
func NewVm(writer io.Writer) *VM {
	v := &VM{w: writer, m: make(RAM, 2000)}
	pc := 0
	sp := 99
	v.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&pc))
	v.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&sp))

	return v
}

// Exec Execute instructions
func (vm *VM) Exec(bc ByteCode) error {
	_ = vm.r[len(vm.r)-1] // bounds check elimination
	_ = vm.m[len(vm.m)-1] // bounds check elimination

	var tmp uint64
	var i *Inst

	for {
		// pull instruction from bytecode at pc position
		i = &bc[*(*uint64)(unsafe.Pointer(&vm.r[ProgramCounter]))]

		switch i.O {

		case Call:
			// add return address to stack, jmp to offset
			tmp = *(*uint64)(unsafe.Pointer(&vm.r[ProgramCounter])) + 1
			vm.m[*(*uint64)(unsafe.Pointer(&vm.r[StackPointer]))] = *(*[8]byte)(unsafe.Pointer(&tmp))
			tmp = *(*uint64)(unsafe.Pointer(&vm.r[StackPointer])) - 1
			vm.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&tmp))
			vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
			continue

		case And:
			switch i.F {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) = *(*bool)(unsafe.Pointer(&vm.r[i.X])) && *(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) = *(*bool)(unsafe.Pointer(&i.X)) && *(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case Or:
			switch i.F {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) = *(*bool)(unsafe.Pointer(&vm.r[i.X])) || *(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) = *(*bool)(unsafe.Pointer(&i.X)) || *(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case Not:
			*(*bool)(unsafe.Pointer(&vm.r[i.Z])) = !*(*bool)(unsafe.Pointer(&vm.r[i.X]))

		case Band:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) & *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case Bor:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) | *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case Bxor:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) ^ *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case Bnot:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = ^*(*uint64)(unsafe.Pointer(&vm.r[i.X]))
			}
		case ShiftL:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) << *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case ShiftR:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) >> *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case Add:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) + *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) + *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&i.X)) + *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) + *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&i.X)) + *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case Sub:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) - *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) - *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&i.X)) - *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) - *(*int64)(unsafe.Pointer(&i.Y))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) - *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&i.X)) - *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) - *(*float64)(unsafe.Pointer(&i.Y))
			}

		case Mul:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) * *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) * *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&i.X)) * *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) * *(*int64)(unsafe.Pointer(&i.Y))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) * *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&i.X)) * *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) * *(*float64)(unsafe.Pointer(&i.Y))
			}

		case Quo:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) = *(*uint64)(unsafe.Pointer(&vm.r[i.X])) / *(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) / *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&i.X)) / *(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) = *(*int64)(unsafe.Pointer(&vm.r[i.X])) / *(*int64)(unsafe.Pointer(&i.Y))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) / *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&i.X)) / *(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) = *(*float64)(unsafe.Pointer(&vm.r[i.X])) / *(*float64)(unsafe.Pointer(&i.Y))
			}
		case Pow:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) =
					uint64(math.Pow(float64(*(*uint64)(unsafe.Pointer(&vm.r[i.X]))), float64(*(*uint64)(unsafe.Pointer(&vm.r[i.Y])))))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[i.X]))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[i.Y])))))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&i.X))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[i.Y])))))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[i.X]))),
						float64(*(*int64)(unsafe.Pointer(&i.Y)))))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[i.X])),
						*(*float64)(unsafe.Pointer(&vm.r[i.Y])))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Pow(*(*float64)(unsafe.Pointer(&i.X)),
						*(*float64)(unsafe.Pointer(&vm.r[i.Y])))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[i.X])),
						*(*float64)(unsafe.Pointer(&i.Y)))
			}

		case Rem:
			switch i.F {
			case None:
				*(*uint64)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) %
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) %
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) %
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case IImm:
				*(*int64)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) %
						*(*int64)(unsafe.Pointer(&i.Y))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[i.X])),
						*(*float64)(unsafe.Pointer(&vm.r[i.Y])))
			case ImmF:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&i.X)),
						*(*float64)(unsafe.Pointer(&vm.r[i.Y])))
			case FImm:
				*(*float64)(unsafe.Pointer(&vm.r[i.Z])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[i.X])),
						*(*float64)(unsafe.Pointer(&i.Y)))
			}

		case Eq:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) ==
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) ==
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) ==
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) ==
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) ==
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*bool)(unsafe.Pointer(&vm.r[i.X])) ==
						*(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*bool)(unsafe.Pointer(&i.X)) ==
						*(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			}
		case NEq:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) !=
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) !=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) !=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) !=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) !=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*bool)(unsafe.Pointer(&vm.r[i.X])) !=
						*(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmB:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*bool)(unsafe.Pointer(&i.X)) !=
						*(*bool)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case LT:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) <
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) <
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) <
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) <
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) <
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case LTE:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) <=
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) <=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) <=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) <=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) <=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case GT:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) >
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) >
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) >
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) >
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) >
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case GTE:
			switch i.F {
			case None:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*uint64)(unsafe.Pointer(&vm.r[i.X])) >=
						*(*uint64)(unsafe.Pointer(&vm.r[i.Y]))
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&vm.r[i.X])) >=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmI:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*int64)(unsafe.Pointer(&i.X)) >=
						*(*int64)(unsafe.Pointer(&vm.r[i.Y]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&vm.r[i.X])) >=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			case ImmF:
				*(*bool)(unsafe.Pointer(&vm.r[i.Z])) =
					*(*float64)(unsafe.Pointer(&i.X)) >=
						*(*float64)(unsafe.Pointer(&vm.r[i.Y]))
			}

		case Print:
			switch i.F {
			case None:
				_, _ = fmt.Fprint(vm.w, vm.r[i.X])
			case Int:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&vm.r[i.X])))
			case ImmI:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&i.X)))
			case Float:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&vm.r[i.X])))
			case ImmF:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&i.X)))
			case Bool:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&vm.r[i.X])))
			case ImmB:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&i.X)))
			}

		case PrintLn:
			switch i.F {
			case None:
				_, _ = fmt.Fprintln(vm.w, &vm.r[i.X])
			case Int:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&vm.r[i.X])))
			case ImmI:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&i.X)))
			case Float:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&vm.r[i.X])))
			case ImmF:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&i.X)))
			case Bool:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&vm.r[i.X])))
			case ImmB:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&i.X)))
			}

		case Load:
			switch i.F {
			case SP:
				tmp = *(*uint64)(unsafe.Pointer(&vm.r[StackPointer])) + 1
				vm.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&tmp))
				vm.r[i.Y] = vm.m[tmp]
			}

		case Store:
			switch i.F {
			case None:
				vm.r[i.Y] = *(*[8]byte)(unsafe.Pointer(&i.X))
			case SP:
				tmp = *(*uint64)(unsafe.Pointer(&vm.r[StackPointer])) - 1
				vm.m[tmp+1] = *(*[8]byte)(unsafe.Pointer(&i.X))
				vm.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&tmp))
			case SPR:
				tmp = *(*uint64)(unsafe.Pointer(&vm.r[StackPointer])) - 1
				vm.m[tmp+1] = vm.r[i.X]
				vm.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&tmp))
			}

		case JMP:
			switch i.F {
			case None:
				vm.r[ProgramCounter] = vm.r[i.X]
				continue
			case Imm:
				vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
				continue
			case SP:
				tmp = *(*uint64)(unsafe.Pointer(&vm.r[StackPointer])) + 1
				vm.r[StackPointer] = *(*[8]byte)(unsafe.Pointer(&tmp))
				vm.r[ProgramCounter] = vm.m[tmp]
				continue
			}

		case JMPEQ:
			switch i.F {
			case None:
				if i.Y == *(*uint64)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}
			case ImmI:
				if *(*int64)(unsafe.Pointer(&i.Y)) == *(*int64)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}
			case ImmB:
				if *(*bool)(unsafe.Pointer(&i.Y)) == *(*bool)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}

			}

		case JMPNEQ:
			switch i.F {
			case None:
				if i.Y != *(*uint64)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}
			case ImmI:
				if *(*int64)(unsafe.Pointer(&i.Y)) != *(*int64)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}
			case ImmB:
				if *(*bool)(unsafe.Pointer(&i.Y)) != *(*bool)(unsafe.Pointer(&vm.r[i.Z])) {
					vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&i.X))
					continue
				}
			}

		case Exit:
			return nil

		case NoOp:
			npc := *(*uint64)(unsafe.Pointer(&vm.r[ProgramCounter])) + 1
			vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&npc))
			continue

		default:
			return fmt.Errorf("unhandled OP %d", i.O)
		}
		npc := *(*uint64)(unsafe.Pointer(&vm.r[ProgramCounter])) + 1
		vm.r[ProgramCounter] = *(*[8]byte)(unsafe.Pointer(&npc))
	}
}
