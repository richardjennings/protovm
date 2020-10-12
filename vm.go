package protovm

import (
	"fmt"
	"io"
	"math"
	"unsafe"
)

type (
	ByteCode  []Inst
	Inst      [5][8]byte
	Registers [][8]byte
	VM        struct {
		w  io.Writer
		pc uint64
		sp uint64
		r  Registers
		s  [][8]byte
		h  []interface{}
	}
)

// Create a new VM struct
func NewVm(writer io.Writer, registers int, stack int) *VM {
	v := &VM{w: writer}
	v.r = make([][8]byte, registers, registers)
	v.s = make([][8]byte, stack, stack)
	return v
}

// Execute instructions
func (vm *VM) Exec(bc ByteCode) error {
	var i Inst
	for {
		i = bc[vm.pc]
		switch Opcode(*(*uint64)(unsafe.Pointer(&i[0]))) {

		// boolean
		case And:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) &&
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateBool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) &&
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case Or:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ||
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateBool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) ||
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case Not:
			*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
				!*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))

		//arithmetic
		case Add:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) +
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) +
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) +
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) +
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case Sub:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) -
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IntImmediate:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) -
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case FloatImmediate:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) -
						*(*float64)(unsafe.Pointer(&i[3]))
			}
		case Mul:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) *
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) *
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) *
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) *
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case Quo:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) /
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IntImmediate:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) /
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case FloatImmediate:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) /
						*(*float64)(unsafe.Pointer(&i[3]))
			}
		case Pow:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&i[2]))),
						float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))))
			case IntImmediate:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					int64(math.Pow(float64(*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))),
						float64(*(*int64)(unsafe.Pointer(&i[3])))))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&i[2])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case FloatImmediate:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Pow(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&i[3])))
			}
		case Rem:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) %
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) %
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case IntImmediate:
				*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) %
						*(*int64)(unsafe.Pointer(&i[3]))
			case Float:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case ImmediateFloat:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&i[2])),
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))])))
			case FloatImmediate:
				*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					math.Remainder(*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])),
						*(*float64)(unsafe.Pointer(&i[3])))
			}

			// equality
		case Eq:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) ==
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) ==
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) ==
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateBool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) ==
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case NEq:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) !=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) !=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Bool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) !=
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateBool:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*bool)(unsafe.Pointer(&i[2])) !=
						*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case LT:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) <
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) <
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case LTE:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) <=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) <=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) <=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case GT:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) >
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) >
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case GTE:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateInt:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*int64)(unsafe.Pointer(&i[2])) >=
						*(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case Float:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])) >=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			case ImmediateFloat:
				*(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) =
					*(*float64)(unsafe.Pointer(&i[2])) >=
						*(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))
			}
		case Print:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateInt:
				_, _ = fmt.Fprint(vm.w, *(*int64)(unsafe.Pointer(&i[2])))
			case Float:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateFloat:
				_, _ = fmt.Fprint(vm.w, *(*float64)(unsafe.Pointer(&i[2])))
			case Bool:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateBool:
				_, _ = fmt.Fprint(vm.w, *(*bool)(unsafe.Pointer(&i[2])))
			}
		case PrintLn:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Int:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateInt:
				_, _ = fmt.Fprintln(vm.w, *(*int64)(unsafe.Pointer(&i[2])))
			case Float:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateFloat:
				_, _ = fmt.Fprintln(vm.w, *(*float64)(unsafe.Pointer(&i[2])))
			case Bool:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))])))
			case ImmediateBool:
				_, _ = fmt.Fprintln(vm.w, *(*bool)(unsafe.Pointer(&i[2])))
			}
		case Load:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case Ptr:
				vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.r[*(*uint64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[3]))]))]
			case SP:
				vm.sp--
				vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.s[vm.sp]
			}
		case Store:
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
		case JMP:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case None:
				vm.pc = *(*uint64)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[2]))]))
				continue
			case Immediate:
				vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
				continue
			}

		case JMPEQ:
			switch Funct(*(*uint64)(unsafe.Pointer(&i[1]))) {
			case ImmediateBool:
				if *(*bool)(unsafe.Pointer(&i[3])) == *(*bool)(unsafe.Pointer(&vm.r[*(*uint64)(unsafe.Pointer(&i[4]))])) {
					vm.pc = *(*uint64)(unsafe.Pointer(&i[2]))
					continue
				}
			}
		case Exit:
			return nil
		default:
			return fmt.Errorf("unhandled Op %b", i[0])
		}
		vm.pc++
	}
}
