# Proto

# VM

## About

Prototype VM for use as a backend to the Elodie programming language.

```example/fib35recursive/main.go``` contains a handcrafted recursive fib 35 implementation
used to ballpark VM performance.

Having tried a couple of different methods for op selection / handling, I found a manually 
implemented search tree with if / else to be the fastest.    
Go 1.x does not provide any mechanism to jump dynamically, so techniques such as Computed Goto cannot be applied.

## Benchmarks
Benchmarks currently include a Fib recursive implementation using a number of n values.    
The results on my system look like:
```
goos: darwin
goarch: amd64
pkg: github.com/richardjennings/protovm
BenchmarkFibRecursive_0-4        4530684               261 ns/op
BenchmarkFibRecursive_1-4        4304808               279 ns/op
BenchmarkFibRecursive_5-4        1417191               842 ns/op
BenchmarkFibRecursive_20-4          1351            859398 ns/op
BenchmarkFibRecursive_30-4            10         103942165 ns/op
BenchmarkFibRecursive_35-4             1        1154557894 ns/op
```

## Run
``` 
$ go build example/fib35recursive/main.go
$ time ./main
9227465

real    0m1.156s
user    0m1.153s
sys     0m0.004s
```

In comparison, PHP 7.3.11 (cli) with cli op-caching enabled achieves:

```
time php fib.php 
9227465

real    0m1.410s
user    0m1.374s
sys     0m0.018s
```

The equivalent Go implementation completes in around 0.07 seconds on my system. Hopefully an x86_64 backend for Elodie
will eventually be able to compete with Go, C and Rust performance. One day.


## Experiments

### Op Selection and Dispatch
Op selection via an array of functions.   
Something like:
```
ops := [50]func(){}{
    func(){}//NoOp
    func(){//Add
        ...
    },
    ...
}
for {
    i := insts[vm.pc]
    ops[i[0]]()
}
```
This worked out at around 2.6 seconds.

Using select case in a for loop.    
Something like:
```
for {
    i := insts[vm.pc]
    select ops[i[0]] {
    case NoOp:
        ...
    case Add:
        ...
    }
    vm.pc++
}
```
This worked out at around 1.35 seconds.


# ASM

## About

The asm package is intended to provide assembly parsing and programmatic construction of Proto VM ByteCode.   

A Builder is included in the asm package to construct the bytecode format expected by the VM.

```
	b := asm.NewBuilder()
	// add a comment for the proceeding instruction
	b.Comment("print immediate int 1")
	// adds a print immmediate int 1 instruction
	b.Add(p.PrintLn, p.ImmI, int64(1))
	// end program
	b.Add(p.Exit)
	
	// print out asm
	fmt.Println(b)
```

