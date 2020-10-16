package asm

import (
	"fmt"
	vm "github.com/richardjennings/proto/vm"
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

func (b *Builder) Add(o vm.Opcode, f vm.Funct, r ...interface{}) {
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
		panic("too many arguments")
	}
	b.asm = append(b.asm, l)
}

func (b *Builder) Label(l string) {
	b.labels[l] = uint64(len(b.asm))
	b.rLabels[uint64(len(b.asm))] = l
}

func (b *Builder) Comment(c string) {
	b.comments[uint64(len(b.asm))] = c
}

func (b *Builder) Exit() {
	b.asm = append(b.asm, line{o: vm.Exit})
}

func (b *Builder) BC() vm.ByteCode {
	for _, v := range b.asm {
		inst := vm.Inst{}
		inst[0] = vm.Op(v.o)
		inst[1] = vm.F(v.f)
		b.operand(v.x, &inst, 2)
		b.operand(v.y, &inst, 3)
		b.operand(v.z, &inst, 4)
		b.bc = append(b.bc, inst)
	}
	for l, positions := range b.labelledInsts {
		if _, ok := b.labels[l]; !ok {
			panic(fmt.Sprintf("label %s not found", l))
		}
		for _, pos := range positions {
			b.bc[pos.i][pos.j] = vm.Uint64(b.labels[l])
		}
	}
	bc := b.bc
	b.bc = vm.ByteCode{}
	return bc
}

func (b *Builder) operand(a interface{}, inst *vm.Inst, i int) {
	switch a := a.(type) {
	case [8]byte:
		inst[i] = a
	case string:
		inst[i] = [8]byte{}
		if _, ok := b.labelledInsts[a]; !ok {
			b.labelledInsts[a] = []pos{}
		}
		b.labelledInsts[a] = append(b.labelledInsts[a], pos{len(b.bc), i})
	case R:
		inst[i] = vm.Uint64(uint64(a))
	case uint64:
		inst[i] = vm.Uint64(a)
	case int64:
		inst[i] = vm.Int64(a)
	case float64:
		inst[i] = vm.Float64(a)
	case bool:
		inst[i] = vm.Boolean(a)
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
