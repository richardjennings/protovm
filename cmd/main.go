package main

import (
	p "github.com/richardjennings/protovm"
	"os"
)
// < me
// > them
func main() {
	w := os.Stdout
	vm := p.NewVm(w, 20, 1000)
	bc := p.ByteCode{
		{p.Op(p.Store), p.F(p.SP), p.Uint64(31)},// <add return position to stack @todo jump address
		{p.Op(p.Store), p.F(p.SP), p.Int64(35)},// >add 35 to stack
// fn fib ($n int) (int)
		{p.Op(p.Load), p.F(p.SP), p.R(1)},// pop fib arg from stack into R1
// $n == 0?
		{p.Op(p.Eq), p.F(p.ImmediateInt), p.Int64(0), p.R(1), p.R(2)},// $n == 0 => R2
		{p.Op(p.JMPEQ), p.F(p.ImmediateBool), p.Uint64(8), p.Boolean(false), p.R(2)}, //@todo update jump address (if $n == 1)
		{p.Op(p.Load), p.F(p.SP), p.R(2)},// pop return address into R2
		{p.Op(p.Store), p.F(p.SP), p.Int64(0)},//>is 0 so push to stack and return
		{p.Op(p.JMP), p.F(p.None), p.R(2)},// goto R2 caller address
// $n == 1?
		{p.Op(p.Eq), p.F(p.ImmediateInt), p.Int64(1), p.R(1), p.R(2)},/// not 0 so try 1
		{p.Op(p.JMPEQ), p.F(p.ImmediateBool), p.Uint64(13), p.Boolean(false), p.R(2)}, //@todo update jump address (fib($n -1 ...
		{p.Op(p.Load), p.F(p.SP), p.R(2)},// pop return address into R2
		{p.Op(p.Store), p.F(p.SP), p.Int64(1)},// push result 1 to stack
		{p.Op(p.JMP), p.F(p.None), p.R(2)},// goto R2 caller address
// fib($n - 1)
		{p.Op(p.Sub), p.F(p.IntImmediate), p.R(1), p.Int64(1), p.R(2)},// n - 1 R2
		{p.Op(p.Store), p.F(p.SPR), p.R(1)},// <push n to stack to recover after recursive call
		{p.Op(p.Store), p.F(p.SP), p.Uint64(18)}, // >push return address to stack @todo jump address
		{p.Op(p.Store), p.F(p.SPR), p.R(2)},// >push n - 1 to stack as argument to recursive call
		{p.Op(p.JMP), p.F(p.Immediate), p.Uint64(2)},// goto recursive function call @todo update jmp address
		{p.Op(p.Load), p.F(p.SP), p.R(2)},// pop result into R(2)
// fib ($n - 2)
		{p.Op(p.Load), p.F(p.SP), p.R(1)},// pop $n off of stack into R1
		{p.Op(p.Store), p.F(p.SPR), p.R(2)},// push previous result back R2 onto the stack
		{p.Op(p.Sub), p.F(p.IntImmediate), p.R(1), p.Int64(2), p.R(2)},// n - 2 R2
		{p.Op(p.Store), p.F(p.SP), p.Uint64(25)}, // >push return address to stack @todo jump address
		{p.Op(p.Store), p.F(p.SPR), p.R(2)},// >push n - 2 to stack as argument to recursive call
		{p.Op(p.JMP), p.F(p.Immediate), p.Uint64(2)},// goto recursive function call @todo update jmp address
// fib($n-1) + fib($n-2)
		{p.Op(p.Load), p.F(p.SP), p.R(2)},// pop fib($n-2) result into R2
		{p.Op(p.Load), p.F(p.SP), p.R(1)},// pop fib($n-1) result into R1
		{p.Op(p.Add), p.F(p.Int), p.R(1), p.R(2), p.R(4)},// add R1 R2 into R1
		{p.Op(p.Load), p.F(p.SP), p.R(2)},// pop return address into R2
		{p.Op(p.Store), p.F(p.SPR), p.R(4)},// push result R1 onto stack
		{p.Op(p.JMP), p.F(p.None), p.R(2)},// jump to return address
// end of fib
		{p.Op(p.Load), p.F(p.SP), p.R(0)},// pop result from stack into r0
		{p.Op(p.PrintLn), p.F(p.Int), p.R(0)},// print result r1
		{p.Op(p.Exit)},// done
	}
	err := vm.Exec(bc)
	if err != nil {
		panic(err)
	}
}
