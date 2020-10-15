# ProtoVM

## About

Prototype VM for Elodie.

```/cmd/main.go``` contains a handcrafted recursive fib 35 implementation
used to ballpark VM performance.

Having tried a couple of different methods for op selection / handling, I found a manually 
implemented search tree with if / else to be the fastest. Go does not provide any mechanism to jump dynamically,
so techniques such as Computed Goto can not be applied.

## Test

```$ go test -v ./...```

## Run
``` 
$ go build cmd/main.go
$ time ./main
9227465

real    0m1.197s
user    0m1.193s
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

## Conclusion

I set out to create a VM that Elodie can target which can run recursive fibonacci 35 faster than PHP 7.3.
This has been achieved but potentially only due to the limited number of Opcodes that are currently implemented in ProtoVM.
I expect the select statement in Go is a limiting factor and will result in a slow down as the VM is further built out. 


## Experiments

Using a select case in a for loop.    
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