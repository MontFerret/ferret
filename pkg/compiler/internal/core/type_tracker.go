package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type TypeTracker struct {
	regs map[bytecode.Operand]ValueType
}

func NewTypeTracker() *TypeTracker {
	return &TypeTracker{
		regs: make(map[bytecode.Operand]ValueType),
	}
}

func (t *TypeTracker) Set(op bytecode.Operand, typ ValueType) {
	if t == nil || !op.IsRegister() {
		return
	}

	t.regs[op] = typ
}

func (t *TypeTracker) Get(op bytecode.Operand) (ValueType, bool) {
	if t == nil || !op.IsRegister() {
		return TypeUnknown, false
	}

	typ, ok := t.regs[op]

	return typ, ok
}

func (t *TypeTracker) Resolve(op bytecode.Operand) ValueType {
	typ, ok := t.Get(op)

	if !ok {
		return TypeUnknown
	}

	return typ
}
