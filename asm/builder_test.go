package asm

import (
	"bytes"
	"testing"
	vm "github.com/richardjennings/proto/vm"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder()
	b.Add(vm.Print, vm.ImmI, int64(1))
	b.Add(vm.JMP, vm.Imm, "end")
	b.Add(vm.Print, vm.ImmI, int64(2))
	b.Label("end")
	b.Add(vm.Print, vm.ImmI, int64(3))
	b.Exit()
	bc := b.BC()

	expected := vm.Uint64(3)
	if bc[1][2] != vm.Uint64(3) {
		t.Errorf("expected jump offset %s got %s", expected, bc[1][2])
	}

	buf := bytes.Buffer{}
	vm := vm.NewVm(&buf)
	err := vm.Exec(bc)
	if err != nil {
		t.Error(err)
	}
}
