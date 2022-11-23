# ASM

## About
Assembler

## Format


Example:
```
section	.text
global _start

# fib(n int) int 
fib:
    jne         $0, %r0, .LB0_1 
    store       $0, %r1
    jmp         %sp
.LB0_1:
    jne         $1, %r0, .LB0_2
    store       $1, %r1
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
```