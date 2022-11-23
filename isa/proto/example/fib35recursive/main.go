package main

import (
	"fmt"
	"github.com/richardjennings/protovm/isa/proto"
	"log"
	"os"
)

func main() {
	b := proto.NewBuilder()
	b.Comment("add return position to stack")
	b.Add(proto.Store, proto.SP, "END")
	b.Comment("set function arg n as R0")
	b.Add(proto.Store, proto.None, uint64(35), proto.R(0))
	b.Comment("call fib function")
	b.Add(proto.JMP, proto.Imm, "fib(n)")
	b.Label("END")
	b.Add(proto.PrintLn, proto.Int, proto.R(1))
	b.Exit()
	b.Label("fib(n)")
	b.Comment("jump to next if condition if n != 0")
	b.Add(proto.JMPNEQ, proto.ImmI, "if n == 1", int64(0), proto.R(0))
	b.Comment("set return value to 0")
	b.Add(proto.Store, proto.None, uint64(0), proto.R(1))
	b.Comment("jump to return address")
	b.Add(proto.JMP, proto.SP)
	b.Label("if n == 1")
	b.Comment("jump to recursive calls if n != 1")
	b.Add(proto.JMPNEQ, proto.ImmI, "fib(n - 1)", int64(1), proto.R(0))
	b.Comment("set return value to 1")
	b.Add(proto.Store, proto.None, uint64(1), proto.R(1))
	b.Comment("jump to return address")
	b.Add(proto.JMP, proto.SP)
	b.Label("fib(n - 1)")
	b.Comment("push n to stack to recover after recursive call")
	b.Add(proto.Store, proto.SPR, proto.R(0))
	b.Comment("n - 1 => R0")
	b.Add(proto.Sub, proto.IImm, proto.R(0), int64(1), proto.R(0))
	b.Comment("push return address to stack")
	b.Add(proto.Store, proto.SP, "fib(n - 2)")
	b.Comment("goto recursive function call")
	b.Add(proto.JMP, proto.Imm, "fib(n)") // return value from call should be in R1
	b.Label("fib(n - 2)")
	b.Comment("pop n off of stack into R0")
	b.Add(proto.Load, proto.SP, nil, proto.R(0))
	b.Comment("push previous result back R1 onto the stack")
	b.Add(proto.Store, proto.SPR, proto.R(1))
	b.Comment("n - 2 => R0")
	b.Add(proto.Sub, proto.IImm, proto.R(0), int64(2), proto.R(0))
	b.Comment("push return address to stack")
	b.Add(proto.Store, proto.SP, "fib + fib")
	b.Comment("goto recursive function call")
	b.Add(proto.JMP, proto.Imm, "fib(n)")
	b.Label("fib + fib")
	b.Comment("pop fib(n-1) result into R2 => fib(n-2) result is in R1")
	b.Add(proto.Load, proto.SP, nil, proto.R(2))
	b.Comment("add R1 R2 into R1")
	b.Add(proto.Add, proto.Int, proto.R(1), proto.R(2), proto.R(1))
	b.Comment("jump to return address")
	b.Add(proto.JMP, proto.SP)
	//

	// print out bc representation
	fmt.Println(b)

	// generate ByteCode
	p, err := b.BC()
	if err != nil {
		log.Fatalln(err)
	}
	// create vm
	vm := proto.NewVm(os.Stdout)
	// execute ByteCode
	if err := vm.Exec(p); err != nil {
		log.Fatalln(err)
	}
}
