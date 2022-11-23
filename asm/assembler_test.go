package asm

import (
	"bytes"
	"testing"
)

// End to End test conversion of src assembly to vm ByteCode
func TestAssembler_Parse(t *testing.T) {
	assembly := `
# text section start
section	.text

# entry point for 'linker'
global _start

# fib(n int) int function
fib:
    jne         $0, %r0, .LB0_1 
	# a comment
    store       $0, %r1
    jmp         %sp
.LB0_1:
    jne         $1, %r0, .LB0_2
    store       $1, %r1 # return 1
    jmp         %sp
.LB0_2:
    store       %r0, %sp
    sub         %r0, $1, %r0
    call        fib
    load        %sp, %r0
    store       %r1, %sp
    sub         %r0, $2, %r0
    call        fib
    load        %sp, %r2
    add         %r1, %r2, %r1
    jmp         %sp
_start:
    store       $35, %r0        # call fib with int 35
    call        fib
    println     %r1
    exit
#comment
	`
	assembler := NewAssembler()
	bc, err := assembler.Parse(bytes.NewBuffer([]byte(assembly)))
	if err != nil {
		t.Fatal(err)
	}
	expected := Asm{
		t: textSection{
			labels: map[string]int{
				".LB0_1": 3,
				".LB0_2": 6,
				"_start": 16,
			},
			functions: map[string]int{
				"fib": 0,
			},
			start: "_start",
			insts: []inst{
				{"jne", "$0", "%r0", ".LB0_1"},
				{"store", "$0", "%r1"},
				{"jmp", "%sp"},
				{"jne", "$1", "%r0", ".LB0_2"},
				{"store", "$1", "%r1"},
				{"jmp", "%sp"},
				{"store", "%r0", "%sp"},
				{"sub", "%r0", "$1", "%r0"},
				{"call", "fib"},
				{"load", "%sp", "%r0"},
				{"store", "%r1", "%sp"},
				{"sub", "%r0", "$2", "%r0"},
				{"call", "fib"},
				{"load", "%sp", "%r2"},
				{"add", "%r1", "%r2", "%r1"},
				{"jmp", "%sp"},
				{"store", "$35", "%r0"},
				{"call", "fib"},
				{"println", "%r1"},
				{"exit"},
			},
		},
	}
	testEqualAsm(assembler.parser.asm, &expected, t)
	_ = bc
}

func testEqualAsm(a *Asm, e *Asm, t *testing.T) {
	if e.t.start != "" {
		if a.t.start != e.t.start {
			t.Error(".text start different")
		}
	}
	// functions
	if len(a.t.functions) != len(e.t.functions) {
		t.Errorf("expected %d got %d for .text functions", len(e.t.functions), len(a.t.functions))
	}
	for v, i := range a.t.functions {
		if j, ok := e.t.functions[v]; !ok {
			t.Error("functions different")
		} else {
			if i != j {
				t.Error("function offset different")
			}
		}
	}
	// labels
	if len(a.t.labels) != len(e.t.labels) {
		t.Errorf("expected %d got %d for .text labels", len(e.t.labels), len(a.t.labels))
	}
	for v, i := range a.t.labels {
		if j, ok := e.t.labels[v]; !ok {
			t.Error("labels different")
		} else {
			if i != j {
				t.Error("function offset different")
			}
		}
	}
	// instructions
	if len(a.t.insts) != len(e.t.insts) {
		t.Error(".tex instructions different length")
	}
	for i, v := range a.t.insts {
		if e.t.insts[i] != v {
			t.Errorf("expected %v got %v", e.t.insts[i], v)
		}
	}
}

// Test conversion of parsed Asm to ByteCode
func TestAssembler_Compile(t *testing.T) {
	asm := Asm{
		t: textSection{
			labels: map[string]int{
				".LB0_1": 3,
				".LB0_2": 6,
				"_start": 16,
			},
			functions: map[string]int{
				"fib": 0,
			},
			start: "_start",
			insts: []inst{
				{"jne", ".LB0_1", "$0", "%r0"},
				{"store", "$0", "%r1"},
				{"jmp", "%sp"},
				{"jne", ".LB0_2", "$1", "%r0"},
				{"store", "$1", "%r1"},
				{"jmp", "%sp"},
				{"store", "%r0", "%sp"},
				{"sub", "%r0", "$1", "%r0"},
				{"call", "fib"},
				{"load", "%sp", "%r0"},
				{"store", "%r1", "%sp"},
				{"sub", "%r0", "$2", "%r0"},
				{"call", "fib"},
				{"load", "%sp", "%r2"},
				{"add", "%r1", "%r2", "%r1"},
				{"jmp", "%sp"},
				{"store", "$35", "%r0"},
				{"call", "fib"},
				{"println", "%r1"},
				{"exit"},
			},
		},
	}

	assembler := NewAssembler()
	bc, err := assembler.Compile(&asm)
	if err != nil {
		t.Error(err)
	}
	_ = bc
}

func TestParser_parse(t *testing.T) {
	for i, tc := range []struct{ src string }{
		{
			src: `
					section .text
					global _start
					_start:

				`,
		},
	} {
		p := newParser()
		err := p.parse(bytes.NewBuffer([]byte(tc.src)))
		if err != nil {
			t.Errorf("test %d: expected no error got %s", i, err)
		}
	}
}

func TestParser_parseError(t *testing.T) {
	for i, tc := range []struct {
		src string
		err string
	}{
		{
			src: "string", err: "0:7 expected section",
		},
		{
			src: "section ",
			err: "0:9 expected .text .data or .bss after section",
		},
		{
			src: `section .text`, //no new line
			err: "0:14 expected 10 got -1",
		},
	} {
		p := newParser()
		src := bytes.NewBuffer([]byte(tc.src))
		err := p.parse(src)
		if err == nil {
			t.Errorf("test %d: expected error '%s' got none", i, tc.err)
		}
		if err != nil && err.Error() != tc.err {
			t.Errorf("test %d: expected error '%s' got '%s'", i, tc.err, err.Error())
		}
	}
}
