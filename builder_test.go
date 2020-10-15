package protovm

import (
	"bytes"
	"testing"
)

func TestBuilder(t *testing.T) {
	b := NewBuilder()
	b.Add(Print, ImmI, int64(1))
	b.Add(JMP, Imm, "end")
	b.Add(Print, ImmI, int64(2))
	b.Label("end")
	b.Add(Print, ImmI, int64(3))
	b.Exit()
	bc := b.BC()

	expected := Uint64(3)
	if bc[1][2] != Uint64(3) {
		t.Errorf("expected jump offset %s got %s", expected, bc[1][2])
	}

	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	err := vm.Exec(bc)
	if err != nil {
		t.Error(err)
	}
}
