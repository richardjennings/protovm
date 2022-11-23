package asm

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/richardjennings/protovm/isa/proto"
	"strconv"
)

type (
	Assembler struct {
		parser *Parser
		bc     proto.ByteCode
		asm    *Asm
	}

	Parser struct {
		scanner Scanner
		section sectionType
		asm     *Asm
	}

	Asm struct {
		t textSection
		d dataSection
		b bssSection
	}

	dataSection struct {
	}

	bssSection struct {
	}

	textSection struct {
		labels    map[string]int
		functions map[string]int
		insts     []inst
		start     string
	}

	inst [4]string

	sectionType uint
)

const (
	None sectionType = iota
	Text
	Data
	Bss
)

func NewAssembler() *Assembler {
	parser := newParser()
	asm := Assembler{parser: parser}
	return &asm
}

func (a *Assembler) Assemble(src *bytes.Buffer) (proto.ByteCode, error) {
	asm, err := a.Parse(src)
	if err != nil {
		return nil, err
	}
	return a.Compile(asm)
}

func (a *Assembler) Parse(src *bytes.Buffer) (*Asm, error) {
	return a.parser.asm, a.parser.parse(src)
}

func (a *Assembler) Compile(asm *Asm) (proto.ByteCode, error) {
	var funct proto.Funct
	var x, y, z interface{}
	var err error
	var op proto.Opcode
	a.asm = asm
	b := proto.NewBuilder()
	for _, v := range asm.t.insts {
		op, err = a.Op(v[0])
		if err != nil {
			return nil, err
		}
		x, err = a.Operand(v[1])
		if err != nil {
			return nil, err
		}
		y, err = a.Operand(v[2])
		if err != nil {
			return nil, err
		}
		z, err = a.Operand(v[3])
		if err != nil {
			return nil, err
		}
		funct, err = b.Funct(op, x, y, z)
		if err != nil {
			return nil, err
		}
		if err := b.Add(op, funct, x, y, z); err != nil {
			return nil, err
		}
	}
	for l, v := range asm.t.labels {
		b.LabelAt(l, v)
	}
	for l, v := range asm.t.functions { //@todo can be duplicates ...
		b.LabelAt(l, v)
	}
	return b.BC()
}

func (a *Assembler) Op(v string) (proto.Opcode, error) {
	op := proto.GetOpcode(v)
	if op == proto.Invalid {
		return proto.Invalid, fmt.Errorf("invalid op %s", v)
	}
	return op, nil
}

func (a *Assembler) Operand(v string) (interface{}, error) {
	if v == "" {
		return nil, nil
	}
	var i int
	var err error

	switch rune(v[0]) {
	case '$':
		i, err = strconv.Atoi(v[1:])
		if err != nil {
			return nil, fmt.Errorf("could not parse %s", v)
		}
		return int64(i), nil
	case '%':
		switch rune(v[1]) {
		case 'r':
			i, err = strconv.Atoi(v[2:])
			if err != nil {
				return nil, fmt.Errorf("could not parse %s", v)
			}
			return proto.RV(uint64(i)), nil
		default:
			if v == "%sp" {
				return proto.SPV(0), nil
			}
		}
	case '.':
		_, ok := a.asm.t.labels[v]
		if ok {
			return v, nil
		}
		return nil, errors.New("label not found")
	default:
		_, ok := a.asm.t.functions[v]
		if ok {
			return v, nil
		}
		return nil, errors.New("function not found")
	}
	return proto.NilV{}, nil
}

func (a *Asm) TextSection() {
	a.t = textSection{}
	a.t.functions = make(map[string]int)
	a.t.labels = make(map[string]int)
}
