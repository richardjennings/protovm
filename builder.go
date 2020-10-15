package protovm

import (
	"fmt"
	"unsafe"
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

func (b *Builder) Add(o Opcode, f Funct, r ...interface{}) {
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
	b.asm = append(b.asm, line{o: Exit})
}

func (b *Builder) BC() ByteCode {
	for _, v := range b.asm {
		inst := Inst{}
		inst[0] = op(v.o)
		inst[1] = f(v.f)
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
			b.bc[pos.i][pos.j] = ui64(b.labels[l])
		}
	}
	bc := b.bc
	b.bc = ByteCode{}
	return bc
}

func (b *Builder) operand(a interface{}, inst *Inst, i int) {
	switch a := a.(type) {
	case [8]byte:
		inst[i] = a
	case string:
		inst[i] = __
		if _, ok := b.labelledInsts[a]; !ok {
			b.labelledInsts[a] = []pos{}
		}
		b.labelledInsts[a] = append(b.labelledInsts[a], pos{len(b.bc), i})
	case R:
		inst[i] = ui64(uint64(a))
	case uint64:
		inst[i] = ui64(a)
	case int64:
		inst[i] = i64(a)
	case float64:
		inst[i] = f64(a)
	case bool:
		inst[i] = boolean(a)
	}
}

func (r R) String() string {
	return fmt.Sprintf("r%d", r)
}
func (n Nil) String() string {
	return " "
}
func (b ByteCode) String() string {
	s := ""
	for i, bc := range b {
		s += fmt.Sprintf("%d %d %d %d %d %d\n", i, bc[0], bc[1], bc[2], bc[3], bc[4])
	}
	return s
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


func op(v Opcode) [8]byte {
	return ui64(uint64(v))
}
func f(v Funct) [8]byte {
	return ui64(uint64(v))
}
func r(v Reg) [8]byte {
	return ui64(uint64(v))
}
func ui64(v uint64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func i64(v int64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func f64(v float64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}
func boolean(v bool) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&v))
}