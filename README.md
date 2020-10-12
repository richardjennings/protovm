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

real    0m2.790s
user    0m2.786s
sys     0m0.006s
```

