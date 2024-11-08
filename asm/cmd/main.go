package main

import (
	"bytes"
	"github.com/richardjennings/protovm/asm"
	"github.com/richardjennings/protovm/isa/proto"
	"log"
	"os"
)

func main() {
	a := asm.NewAssembler()
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	bc, err := a.Assemble(bytes.NewBuffer(b))
	if err != nil {
		log.Fatalln(err)
	}
	vm := proto.NewVm(os.Stdout, 100, bc)
	if err := vm.Exec(bc); err != nil {
		log.Fatalln(err)
	}
}
