package protovm

// Operations
type Opcode uint64
type Funct uint64
type Reg uint64

const (
	None Funct = iota
	Imm
	Int
	ImmI
	IImm
	Float
	ImmF
	FImm
	Bool
	ImmB
	Ptr
	SP
	SPR
)

var functs = map[Funct]string{
	None:  "",
	Imm:   "Imm",
	Int:   "Int",
	ImmI:  "ImmI",
	IImm:  "IImm",
	Float: "Float",
	ImmF:  "ImmF",
	FImm:  "FImm",
	Bool:  "Bool",
	ImmB:  "ImmB",
	Ptr:   "Ptr",
	SP:    "SP",
	SPR:   "SPR",
}

var __ [8]byte

const (
	NoOp Opcode = iota

	And
	Or
	Not

	Add
	Sub
	Mul
	Quo
	Pow
	Rem

	Eq
	NEq
	LT
	LTE
	GT
	GTE

	Print
	PrintLn

	Load
	Store

	JMP
	JMPEQ
	JMPNEQ

	Exit
)

var opcodes = map[Opcode]string{
	NoOp:    "NoOp",
	And:     "And",
	Or:      "Or",
	Not:     "Not",
	Add:     "Add",
	Sub:     "Sub",
	Mul:     "Mul",
	Quo:     "Quo",
	Pow:     "Pow",
	Rem:     "Rem",
	Eq:      "Eq",
	NEq:     "NEq",
	LT:      "LT",
	LTE:     "LTE",
	GT:      "GT",
	GTE:     "GTE",
	Print:   "Print",
	PrintLn: "PrintLn",
	Load:    "Load",
	Store:   "Store",
	JMP:     "JMP",
	JMPEQ:   "JMPEQ",
	JMPNEQ:  "JMPNEQ",
	Exit:    "Exit",
}

func (f Funct) String() string {
	if n, ok := functs[f]; ok {
		return n
	}
	panic("func string not found")
}

func (o Opcode) String() string {
	if n, ok := opcodes[o]; ok {
		return n
	}
	panic("opcode string not found")
}

