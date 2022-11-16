# ProtoVM

## About

This is a register based virtual machine, instruction set and byte-code builder I wrote whilst
playing around with writing a language in GO some time ago. The intention was to calculate fib(35) recursively
quicker than other interpreted languages. This example does not include the parser / frontend and
mostly serves for my own historic reference. The ISA is not complete. I updated the dispatch loop to use a switch 
statement as Go 1.19 improved `switch` performance:
> The compiler now uses a jump table to implement large integer and string switch statements. Performance improvements 
> for the switch statement vary but can be on the order of 20% faster. (GOARCH=amd64 and GOARCH=arm64 only)

and tidied up a number of things to make this nicer to pick up and play with if and when wanted.

```
go build example/fib35recursive/main.go
time ./main
         0 Store    SP       END                        //add return position to stack
         1 Store             35       r0                //set function arg n as R0
         2 JMP      Imm      fib(n)                     //call fib function
END:
         3 PrintLn  Int      r1                         //
         4 Exit                                         //
fib(n):
         5 JMPNEQ   ImmI     if n == 1 0        r0       //jump to next if condition if n != 0
         6 Store             0        r1                //set return value to 0
         7 JMP      SP                                  //jump to return address
if n == 1:
         8 JMPNEQ   ImmI     fib(n - 1) 1        r0       //jump to recursive calls if n != 1
         9 Store             1        r1                //set return value to 1
        10 JMP      SP                                  //jump to return address
fib(n - 1):
        11 Store    SPR      r0                         //push n to stack to recover after recursive call
        12 Sub      IImm     r0       1        r0       //n - 1 => R0
        13 Store    SP       fib(n - 2)                   //push return address to stack
        14 JMP      Imm      fib(n)                     //goto recursive function call
fib(n - 2):
        15 Load     SP       r0                         //pop n off of stack into R0
        16 Store    SPR      r1                         //push previous result back R1 onto the stack
        17 Sub      IImm     r0       2        r0       //n - 2 => R0
        18 Store    SP       fib + fib                   //push return address to stack
        19 JMP      Imm      fib(n)                     //goto recursive function call
fib + fib:
        20 Load     SP       r2                         //pop fib(n-1) result into R2 => fib(n-2) result is in R1
        21 Add      Int      r1       r2       r1       //add R1 R2 into R1
        22 JMP      SP                                  //jump to return address

9227465
./main  0.71s user 0.01s system 99% cpu 0.716 total
```

I did have an assembler written for this which I would like to include at some point.

I moved on from this project to implementing some of the RISC-V ISA which I did not finish. Something like this was intended
as one of the low level intermediate representations for the compiler backend, similar to how Go leverages its own
intermediate language.

## ISA

The ISA includes the following OPs:
```
And      - [None|Bool|ImmB],                     R(Z) = X && Y
Or       - [None|Bool|ImmB],                     R(Z) = X || Y
Not      - [],                                   R(Z) = ! R(X)
Add      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = X + Y
Sub      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = X +- Y
Mul      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = X * Y
Quo      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = X / Y
Pow      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = POW(X, Y)
Rem      - [None|Int|ImmI|IImm|Float|ImmF|FImm], R(Z) = X % Y
Eq       - [None|Int|ImmI|Float|ImmF|Bool|ImmB], R(Z) = X == Y
NEq      - [None|Int|ImmI|Float|ImmF|Bool|ImmB], R(Z) = X != Y
LT       - [None|Int|ImmI|Float|ImmF],           R(Z) = X < Y
LTE      - [None|Int|ImmI|Float|ImmF],           R(Z) = X <= Y      
GT       - [None|Int|ImmI|Float|ImmF],           R(Z) = X > Y
GTE      - [None|Int|ImmI|Float|ImmF],           R(Z) = X >= Y  
Print    - [None|Int|ImmI|Float|ImmF|Bool|ImmB], PRINT( X )
PrintLn  - [None|Int|ImmI|Float|ImmF|Bool|ImmB], PRINTLN( X )
Load     - [SP],                                 R(X) = LOAD(SP)
Store    - [None|SP|SPR],                        R(X) = Y, MEM[SP] = X; SP++; MEM[SP] = R(X)
JMP      - [None|Imm|SP],                        PC = R(X); PC = X; PC = MEM[SP]
JMPEQ    - [None|ImmB|ImmI],                     IF Y == Z THEN PC = X; 
JMPNEQ   - [None|ImmB|ImmI],                     IF Y != Z THEN PC = X; 
Exit     - [],                                   EXIT                  
```

An Instruction is defined as:
```
type Inst struct {
	O Opcode
	F Funct
	X uint64
	Y uint64
	Z uint64
}
```

Because `clock cycles` are expensive via the dispatch loop, the instruction is large to prevent needing multiple
instructions for 64 bit values.


The instruction format is:

`|-   OP   -|-  FUNCT  -|-    X    -|-    Y    -|-    Z    -|`

Where `X` and `Y` can be a register or immediate value. `Z` specifies the result register.

FUNCT values may be:

```
Imm   - An Immediate (X)
Int   - Treat Operands as a 64 bit Int
ImmI  - Treat X as an Immediate Int and Y as an Int in a register
IImm  - Treat X as an Int in a register and Y as an Immediate Int
Float - Treat Operands as a 64 bit Float
ImmF  - Treat X as an Immediate Float and Y as an Float in a register
FImm  - Treat X as an Float in a register and Y as an Immediate Float
Bool  - Treat Operands as a Bool
ImmB  - Treat X as a Bool and Y as a Bool in a register
SP    - Opertation with reference to Stack Pointer value
SPR   - Stack Pointer Operation with reference to a Register
```

There is no Heap, only a Stack of size `[100][8]byte`. The stack was arranged like this to make it easy to align types
in the VM. For example a Load from the Stack Pointer into a register:
```
vm.r[*(*uint64)(unsafe.Pointer(&i[2]))] = vm.s[vm.sp]
```
and then setting a value from a register into the stack:

```
vm.s[vm.sp] = i[2]
```

The VM has only a Program Counter and a Stack Pointer.


## Calling Convention

Any register values required to be restored after a function call are added to the stack.

The return address is pushed to the stack. 
Using ASM builder this is expressed as a Label, e.g. `b.Add(vm.Store, vm.SP, "END")`

Arguments to the function are added to Registers from 0 onwards.

It is the Callees responsibility to jump back to the callers address pointed to by the Stack Pointer.

The Caller should then restore register values from the stack as needed.


