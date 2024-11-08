package proto

import (
	"errors"
	"fmt"
)

// An easier way to construct bytecode format programmatically, with labeled jumps to negate
// annoying hard code address requirement when constructing manually.
type (
	Builder struct {
		// labels is a map of label strings and positions enabling
		// pass to update with offsets
		labels  map[string]uint64
		rLabels map[uint64]string

		// comments
		comments map[uint64]string

		// map of label to instruction and op positions that require updating
		labelledInsts map[string][]pos

		// the ByteCode
		bc ByteCode

		// lines
		asm asm
	}

	pos struct {
		i int
		j int
	}

	asm  []line
	line struct {
		o Opcode
		f Funct
		x interface{}
		y interface{}
		z interface{}
	}
	// register, float, bool, stack pointer and none values
	RV   uint64
	FV   uint64
	BV   uint64
	SPV  uint64
	NilV [8]byte
)

func NewBuilder() *Builder {
	b := Builder{}
	b.bc = nil
	b.labels = make(map[string]uint64)
	b.rLabels = make(map[uint64]string)
	b.comments = make(map[uint64]string)
	b.labelledInsts = make(map[string][]pos)
	return &b
}

func (b *Builder) Add(o Opcode, f Funct, r ...interface{}) error {
	var l line
	switch len(r) {
	case 0:
		l = line{o, f, NilV{}, NilV{}, NilV{}}
	case 1:
		l = line{o, f, r[0], NilV{}, NilV{}}
	case 2:
		l = line{o, f, r[0], r[1], NilV{}}
	case 3:
		l = line{o, f, r[0], r[1], r[2]}
	default:
		return errors.New("too many arguments")
	}
	b.asm = append(b.asm, l)
	return nil
}

func (b *Builder) LabelAt(l string, line int) {
	b.labels[l] = uint64(line)
	b.rLabels[uint64(line)] = l
}

func (b *Builder) Label(l string) {
	b.labels[l] = uint64(len(b.asm))
	b.rLabels[uint64(len(b.asm))] = l
}

func (b *Builder) Comment(c string) {
	b.comments[uint64(len(b.asm))] = c
}

func (b *Builder) Exit() {
	b.asm = append(b.asm, line{o: Exit, f: None, x: NilV{}, y: NilV{}, z: NilV{}})
}

func (b *Builder) BC() (ByteCode, error) {
	for _, v := range b.asm {
		inst := Inst{}
		inst.O = v.o
		inst.F = v.f
		b.operand(v.x, &inst, 2)
		b.operand(v.y, &inst, 3)
		b.operand(v.z, &inst, 4)
		b.bc = append(b.bc, inst)
	}
	for l, positions := range b.labelledInsts {
		if _, ok := b.labels[l]; !ok {
			return nil, fmt.Errorf("label %s not found", l)
		}
		for _, pos := range positions {
			if pos.j == 2 {
				b.bc[pos.i].X = b.labels[l]
			} else if pos.j == 3 {
				b.bc[pos.i].Y = b.labels[l]
			} else {
				b.bc[pos.i].Z = b.labels[l]
			}
		}
	}
	bc := b.bc
	b.bc = ByteCode{}
	return bc, nil
}

func (b *Builder) operand(a interface{}, inst *Inst, i int) {
	var prop *uint64
	switch i {
	case 2:
		prop = &inst.X
	case 3:
		prop = &inst.Y
	case 4:
		prop = &inst.Z
	}
	switch a := a.(type) {
	case string:
		*prop = 0
		if _, ok := b.labelledInsts[a]; !ok {
			b.labelledInsts[a] = []pos{}
		}
		b.labelledInsts[a] = append(b.labelledInsts[a], pos{len(b.bc), i})
	case RV:
		*prop = uint64(a)
	case uint64:
		*prop = a
	case int64:
		*prop = Int64(a)
	case float64:
		*prop = Float64(a)
	case bool:
		*prop = Boolean(a)
	}
}

func (r RV) String() string {
	return fmt.Sprintf("r%d", r)
}
func (n NilV) String() string {
	return " "
}

func (b *Builder) String() string {
	var s string
	for l, i := range b.asm {
		if label, ok := b.rLabels[uint64(l)]; ok {
			s += fmt.Sprintf("%s:\n", label)
		}
		s += fmt.Sprintf("%10d %-8v %-8v %-8v %-8v %-8v //%s\n", l, i.o, i.f, i.x, i.y, i.z, b.comments[uint64(l)])
	}
	return s
}

func (b *Builder) Funct(op Opcode, x interface{}, y interface{}, z interface{}) (Funct, error) {

	// modifies inst to reflect line
	// determines Funct from op and args
	var funct Funct

	switch op {
	case NoOp:
		funct = None
	case And, Or:
		switch x.(type) {
		case bool:
			funct = ImmB
		case RV:
			funct = Bool
		default:
			return None, errors.New("invalid")
		}
	case Not:
		funct = None
	case Add, Mul, LT, LTE, GT, GTE:
		switch x.(type) {
		case int64:
			funct = ImmI
		case float64:
			funct = ImmF
		case RV:
			funct = Int
		case FV:
			funct = Float
		default:
			return None, errors.New("invalid")
		}
	case Sub, Quo, Pow, Rem:
		switch x.(type) {
		case int64:
			funct = ImmI
		case float64:
			funct = ImmF
		case RV:
			funct = Int
		case FV:
			funct = Float
		default:
			return None, errors.New("invalid")
		}
		switch y.(type) {
		case int64:
			funct = IImm
		case float64:
			funct = FImm
		case RV, FV:
		default:
			return None, errors.New("invalid")
		}
	case Eq, NEq, Print, PrintLn:
		switch x.(type) {
		case int64:
			funct = ImmI
		case float64:
			funct = ImmF
		case bool:
			funct = ImmB
		case RV:
			funct = Int
		case FV:
			funct = Float
		case BV:
			funct = Bool
		default:
			return None, errors.New("invalid")
		}
		//case vm.Print:
		//case vm.PrintLn:
	case Load:
		switch x.(type) {
		case SPV:
			funct = SP
		default:
			return None, errors.New("invalid")
		}
	case Store:
		switch x.(type) {
		case uint64, int64, float64, bool:
			funct = None
		case RV:
			switch y.(type) {
			case SPV:
				funct = SPR
			default:
				return None, errors.New("invalid")
			}
			//case SPRV:
			//	Funct = vm.SPR
			//case SP:
			//	Funct = vm.SP
		default:
			return None, errors.New("invalid")
		}
	case JMP:
		switch x.(type) {
		case string:
			funct = Imm
		case SPV:
			funct = SP
		case RV:
			funct = None
		default:
			return None, errors.New("invalid")
		}
	case JMPEQ, JMPNEQ:
		switch y.(type) {
		case bool:
			funct = ImmB
		case int64:
			funct = ImmI
		default:
			return None, errors.New("invalid")
		}
	case Exit:
		funct = None

	}
	return funct, nil
}
