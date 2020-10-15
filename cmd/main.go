package main

import (
	"fmt"
	p "github.com/richardjennings/protovm"
	"os"
)

func main() {
	b := p.NewBuilder()
	b.Comment("add return position to stack")
	b.Add(p.Store, p.SP, "END")
	b.Comment("set function arg $n as R0")
	b.Add(p.Store, p.None, uint64(35), p.R(0))
	b.Comment("call fib function")
	b.Add(p.JMP, p.Imm,"fib($n)")
	b.Label("END")
	b.Add(p.PrintLn, p.Int, p.R(1))
	b.Exit()
	b.Label("fib($n)")
	b.Comment("jump to next if condition if $n != 0")
	b.Add(p.JMPNEQ, p.ImmI, "if $n == 1", int64(0), p.R(0))
	b.Comment("set return value to 0")
	b.Add(p.Store, p.None, uint64(0), p.R(1))
	b.Comment("jump to return address")
	b.Add(p.JMP, p.SP)
	//
	b.Label("if $n == 1")
	b.Comment("jump to recursive calls if $n != 1")
	b.Add(p.JMPNEQ, p.ImmI, "fib($n - 1)", int64(1), p.R(0))
	b.Comment("set return value to 1")
	b.Add(p.Store, p.None, uint64(1), p.R(1))
	b.Comment("jump to return address")
	b.Add(p.JMP, p.SP)
	//
	b.Label("fib($n - 1)")
	b.Comment("push n to stack to recover after recursive call")
	b.Add(p.Store, p.SPR, p.R(0))
	b.Comment("$n - 1 => R0")
	b.Add(p.Sub, p.IImm, p.R(0), int64(1), p.R(0))
	b.Comment("push return address to stack")
	b.Add(p.Store, p.SP, "fib($n - 2)")
	b.Comment("goto recursive function call")
	b.Add(p.JMP, p.Imm, "fib($n)") // return value from call should be in R1
	//
	b.Label("fib($n - 2)")
	b.Comment("pop $n off of stack into R0")
	b.Add(p.Load, p.SP, p.R(0))
	b.Comment("push previous result back R1 onto the stack")
	b.Add(p.Store, p.SPR, p.R(1))
	b.Comment("$n - 2 => R0")
	b.Add(p.Sub, p.IImm, p.R(0), int64(2), p.R(0))
	b.Comment("push return address to stack")
	b.Add(p.Store, p.SP, "fib($n - 1) + fib($n - 2)")
	b.Comment("goto recursive function call")
	b.Add(p.JMP, p.Imm, "fib($n)")
	//
	b.Label("fib($n - 1) + fib($n - 2)")
	b.Comment("pop fib($n-1) result into R2 => fib($n-2) result is in R1")
	b.Add(p.Load, p.SP, p.R(2))
	b.Comment("add R1 R2 into R1")
	b.Add(p.Add, p.Int, p.R(1), p.R(2), p.R(1))
	b.Comment("jump to return address")
	b.Add(p.JMP, p.SP)
	//


	fmt.Println(b)

	bc := b.BC()

	w := os.Stdout
	vm := p.NewVm(w)
	err := vm.Exec(bc)
	if err != nil {
		panic(err)
	}
}
