package core

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	RegisterSequence []vm.Operand

	RegisterAllocator struct {
		next vm.Operand
	}
)

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		next: vm.NoopOperand + 1,
	}
}

func (ra *RegisterAllocator) Allocate() vm.Operand {
	reg := ra.next
	ra.next++

	return reg
}

func (ra *RegisterAllocator) AllocateSequence(count int) RegisterSequence {
	seq := make(RegisterSequence, count)
	start := ra.next

	for i := 0; i < count; i++ {
		seq[i] = start + vm.Operand(i)
	}

	ra.next += vm.Operand(count)
	return seq
}

func (ra *RegisterAllocator) Size() int {
	return int(ra.next)
}
