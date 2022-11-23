package main

import (
	"bytes"
	asm2 "github.com/richardjennings/protovm/asm"
	"github.com/richardjennings/protovm/vm"
	"log"
	"os"
)

func main() {
	asm := asm2.NewAssembler()
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	bc, err := asm.Assemble(bytes.NewBuffer(b))
	if err != nil {
		log.Fatalln(err)
	}
	proto := vm.NewVm(os.Stdout)
	if err := proto.Exec(bc); err != nil {
		log.Fatalln(err)
	}
}
