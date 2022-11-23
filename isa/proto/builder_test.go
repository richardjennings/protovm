package proto

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
	bc, err := b.BC()
	if err != nil {
		t.Error(err)
	}
	expected := 3
	if bc[1].X != 3 {
		t.Errorf("expected jump offset %d got %d", expected, bc[1].X)
	}

	buf := bytes.Buffer{}
	vm := NewVm(&buf)
	if err := vm.Exec(bc); err != nil {
		t.Error(err)
	}
}
