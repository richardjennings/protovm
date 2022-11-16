package main

import (
	"fmt"
	"github.com/richardjennings/protovm/asm"
	"github.com/richardjennings/protovm/vm"
	"log"
	"os"
)

func main() {
	b := asm.NewBuilder()
	b.Comment("add return position to stack")
	b.Add(vm.Store, vm.SP, "END")
	b.Comment("set function arg n as R0")
	b.Add(vm.Store, vm.None, uint64(35), asm.R(0))
	b.Comment("call fib function")
	b.Add(vm.JMP, vm.Imm, "fib(n)")
	b.Label("END")
	b.Add(vm.PrintLn, vm.Int, asm.R(1))
	b.Exit()
	b.Label("fib(n)")
	b.Comment("jump to next if condition if n != 0")
	b.Add(vm.JMPNEQ, vm.ImmI, "if n == 1", int64(0), asm.R(0))
	b.Comment("set return value to 0")
	b.Add(vm.Store, vm.None, uint64(0), asm.R(1))
	b.Comment("jump to return address")
	b.Add(vm.JMP, vm.SP)
	b.Label("if n == 1")
	b.Comment("jump to recursive calls if n != 1")
	b.Add(vm.JMPNEQ, vm.ImmI, "fib(n - 1)", int64(1), asm.R(0))
	b.Comment("set return value to 1")
	b.Add(vm.Store, vm.None, uint64(1), asm.R(1))
	b.Comment("jump to return address")
	b.Add(vm.JMP, vm.SP)
	b.Label("fib(n - 1)")
	b.Comment("push n to stack to recover after recursive call")
	b.Add(vm.Store, vm.SPR, asm.R(0))
	b.Comment("n - 1 => R0")
	b.Add(vm.Sub, vm.IImm, asm.R(0), int64(1), asm.R(0))
	b.Comment("push return address to stack")
	b.Add(vm.Store, vm.SP, "fib(n - 2)")
	b.Comment("goto recursive function call")
	b.Add(vm.JMP, vm.Imm, "fib(n)") // return value from call should be in R1
	b.Label("fib(n - 2)")
	b.Comment("pop n off of stack into R0")
	b.Add(vm.Load, vm.SP, asm.R(0))
	b.Comment("push previous result back R1 onto the stack")
	b.Add(vm.Store, vm.SPR, asm.R(1))
	b.Comment("n - 2 => R0")
	b.Add(vm.Sub, vm.IImm, asm.R(0), int64(2), asm.R(0))
	b.Comment("push return address to stack")
	b.Add(vm.Store, vm.SP, "fib + fib")
	b.Comment("goto recursive function call")
	b.Add(vm.JMP, vm.Imm, "fib(n)")
	b.Label("fib + fib")
	b.Comment("pop fib(n-1) result into R2 => fib(n-2) result is in R1")
	b.Add(vm.Load, vm.SP, asm.R(2))
	b.Comment("add R1 R2 into R1")
	b.Add(vm.Add, vm.Int, asm.R(1), asm.R(2), asm.R(1))
	b.Comment("jump to return address")
	b.Add(vm.JMP, vm.SP)
	//

	// print out ASM representation
	fmt.Println(b)

	// generate ByteCode
	bc, err := b.BC()
	if err != nil {
		log.Fatalln(err)
	}
	// create vm
	vm := vm.NewVm(os.Stdout)
	// execute ByteCode
	if err := vm.Exec(bc); err != nil {
		log.Fatalln(err)
	}
}
