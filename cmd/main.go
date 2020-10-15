package main

import (
	p "github.com/richardjennings/protovm"
	"os"
)
// < me
// > them
func main() {
	b := p.NewBuilder()
	// add return position to stack
	b.Add(p.Store, p.SP, "END")
	// add function arg $n to the stack
	b.Add(p.Store, p.SP, uint64(35))
	//
	b.Label("fib($n)")
	// pop fib arg from stack into R1
	b.Add(p.Load, p.SP, p.R(1))
	// jump to next if condition if $n != 0
	b.Add(p.JMPNEQ, p.ImmI, "if $n == 1", int64(0), p.R(1))
	// pop return address into R2
	b.Add(p.Load, p.SP, p.R(2))
	// push 0 to stack and return
	b.Add(p.Store, p.SP, int64(0))
	// goto R2 caller address
	b.Add(p.JMP, p.None,  p.R(2))
	//
	b.Label("if $n == 1")
	// jump to recursive calls if $n != 1
	b.Add(p.JMPNEQ, p.ImmI, "fib($n - 1)", int64(1), p.R(1))
	// pop return address into R2
	b.Add(p.Load, p.SP, p.R(2))
	// push 0 to stack and return
	b.Add(p.Store, p.SP, int64(1)) //10
	// goto R2 caller address
	b.Add(p.JMP, p.None, p.R(2))
	//
	b.Label("fib($n - 1)")
	// $n - 1 => R2
	b.Add(p.Sub, p.IImm, p.R(1), int64(1), p.R(2))
	// push n to stack to recover after recursive call
	b.Add(p.Store, p.SPR, p.R(1))
	// push return address to stack
	b.Add(p.Store, p.SP, "fib($n - 2)")
	// push n - 1 to stack as argument to recursive call
	b.Add(p.Store, p.SPR, p.R(2))
	// goto recursive function call
	b.Add(p.JMP, p.Imm, "fib($n)")
	//
	b.Label("fib($n - 2)")
	// pop result into r(2)
	b.Add(p.Load, p.SP, p.R(2))
	// pop $n off of stack into R1
	b.Add(p.Load, p.SP, p.R(1))
	// push previous result back R2 onto the stack
	b.Add(p.Store, p.SPR, p.R(2))
	// $n - 2 => R2
	b.Add(p.Sub, p.IImm, p.R(1), int64(2), p.R(2))
	// push return address to stack
	b.Add(p.Store, p.SP, "fib($n - 1) + fib($n - 2)")
	// push n - 2 to stack as argument to recursive call
	b.Add(p.Store, p.SPR, p.R(2))
	// goto recursive function call
	b.Add(p.JMP, p.Imm, "fib($n)")
	//
	b.Label("fib($n - 1) + fib($n - 2)")
	// pop fib($n-2) result into R2
	b.Add(p.Load, p.SP, p.R(2))
	// pop fib($n-1) result into R1
	b.Add(p.Load, p.SP, p.R(1))
	// add R1 R2 into R1
	b.Add(p.Add, p.Int, p.R(1), p.R(2), p.R(4))
	// pop return address into R2
	b.Add(p.Load, p.SP, p.R(2))
	// push result R1 onto stack
	b.Add(p.Store, p.SPR, p.R(4))
	// jump to return address
	b.Add(p.JMP, p.None, p.R(2))
	//
	b.Label("END")
	b.Add(p.Load, p.SP, p.R(0))
	b.Add(p.PrintLn, p.Int, p.R(0))
	b.Exit()

	//fmt.Println(b)

	bc := b.BC()

	w := os.Stdout
	vm := p.NewVm(w, 10, 100)
	err := vm.Exec(bc)
	if err != nil {
		panic(err)
	}
}
