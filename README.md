# ProtoVM

## About

Prototype VM for Elodie.

```/cmd/main.go``` contains a handcrafted recursive fib 35 implementation
used to ballpark VM performance.

## Test

```$ go test -v ./...```

## Run
``` 
$ go build cmd/main.go
$ time ./main
9227465

real    0m1.366s
user    0m1.359s
sys     0m0.005s
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


