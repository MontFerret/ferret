package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	RegisterSequence []bytecode.Operand

	RegisterAllocator struct {
		next bytecode.Operand
	}
)

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		next: bytecode.NoopOperand + 1,
	}
}

func (ra *RegisterAllocator) Allocate() bytecode.Operand {
	reg := ra.next
	ra.next++

	return reg
}

func (ra *RegisterAllocator) AllocateSequence(count int) RegisterSequence {
	seq := make(RegisterSequence, count)
	start := ra.next

	for i := 0; i < count; i++ {
		seq[i] = start + bytecode.Operand(i)
	}

	ra.next += bytecode.Operand(count)
	return seq
}

func (ra *RegisterAllocator) Size() int {
	return int(ra.next)
}
