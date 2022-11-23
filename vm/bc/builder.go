package bc

import (
	"errors"
	"fmt"
	"github.com/richardjennings/protovm/vm"
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
		bc vm.ByteCode

		// lines
		asm asm
	}

	pos struct {
		i int
		j int
	}

	asm  []line
	line struct {
		o vm.Opcode
		f vm.Funct
		x interface{}
		y interface{}
		z interface{}
	}
	R   uint64
	F   uint64
	B   uint64
	SP  uint64
	Nil [8]byte
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

func (b *Builder) Add(o vm.Opcode, f vm.Funct, r ...interface{}) error {
	var l line
	switch len(r) {
	case 0:
		l = line{o, f, Nil{}, Nil{}, Nil{}}
	case 1:
		l = line{o, f, r[0], Nil{}, Nil{}}
	case 2:
		l = line{o, f, r[0], r[1], Nil{}}
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
	b.asm = append(b.asm, line{o: vm.Exit, f: vm.None, x: Nil{}, y: Nil{}, z: Nil{}})
}

func (b *Builder) BC() (vm.ByteCode, error) {
	for _, v := range b.asm {
		inst := vm.Inst{}
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
	b.bc = vm.ByteCode{}
	return bc, nil
}

func (b *Builder) operand(a interface{}, inst *vm.Inst, i int) {
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
	case R:
		*prop = uint64(a)
	case uint64:
		*prop = a
	case int64:
		*prop = vm.Int64(a)
	case float64:
		*prop = vm.Float64(a)
	case bool:
		*prop = vm.Boolean(a)
	}
}

func (r R) String() string {
	return fmt.Sprintf("r%d", r)
}
func (n Nil) String() string {
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

func (b *Builder) Funct(op vm.Opcode, x interface{}, y interface{}, z interface{}) (vm.Funct, error) {

	// modifies inst to reflect line
	// determines Funct from op and args
	var funct vm.Funct

	switch op {
	case vm.NoOp:
		funct = vm.None
	case vm.And, vm.Or:
		switch x.(type) {
		case bool:
			funct = vm.ImmB
		case R:
			funct = vm.Bool
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.Not:
		funct = vm.None
	case vm.Add, vm.Mul, vm.LT, vm.LTE, vm.GT, vm.GTE:
		switch x.(type) {
		case int64:
			funct = vm.ImmI
		case float64:
			funct = vm.ImmF
		case R:
			funct = vm.Int
		case F:
			funct = vm.Float
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.Sub, vm.Quo, vm.Pow, vm.Rem:
		switch x.(type) {
		case int64:
			funct = vm.ImmI
		case float64:
			funct = vm.ImmF
		case R:
			funct = vm.Int
		case F:
			funct = vm.Float
		default:
			return vm.None, errors.New("invalid")
		}
		switch y.(type) {
		case int64:
			funct = vm.IImm
		case float64:
			funct = vm.FImm
		case R, F:
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.Eq, vm.NEq, vm.Print, vm.PrintLn:
		switch x.(type) {
		case int64:
			funct = vm.ImmI
		case float64:
			funct = vm.ImmF
		case bool:
			funct = vm.ImmB
		case R:
			funct = vm.Int
		case F:
			funct = vm.Float
		case B:
			funct = vm.Bool
		default:
			return vm.None, errors.New("invalid")
		}
		//case vm.Print:
		//case vm.PrintLn:
	case vm.Load:
		switch x.(type) {
		case SP:
			funct = vm.SP
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.Store:
		switch x.(type) {
		case uint64, int64, float64, bool:
			funct = vm.None
		case R:
			switch y.(type) {
			case SP:
				funct = vm.SPR
			default:
				return vm.None, errors.New("invalid")
			}
			//case SPR:
			//	Funct = vm.SPR
			//case SP:
			//	Funct = vm.SP
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.JMP:
		switch x.(type) {
		case string:
			funct = vm.Imm
		case SP:
			funct = vm.SP
		case R:
			funct = vm.None
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.JMPEQ, vm.JMPNEQ:
		switch y.(type) {
		case bool:
			funct = vm.ImmB
		case int64:
			funct = vm.ImmI
		default:
			return vm.None, errors.New("invalid")
		}
	case vm.Exit:
		funct = vm.None

	}
	return funct, nil
}
