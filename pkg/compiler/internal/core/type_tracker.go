package core

import "github.com/MontFerret/ferret/pkg/vm"

type TypeTracker struct {
	regs map[vm.Operand]ValueType
}

func NewTypeTracker() *TypeTracker {
	return &TypeTracker{
		regs: make(map[vm.Operand]ValueType),
	}
}

func (t *TypeTracker) Set(op vm.Operand, typ ValueType) {
	if t == nil || !op.IsRegister() {
		return
	}

	t.regs[op] = typ
}

func (t *TypeTracker) Get(op vm.Operand) (ValueType, bool) {
	if t == nil || !op.IsRegister() {
		return TypeUnknown, false
	}

	typ, ok := t.regs[op]

	return typ, ok
}

func (t *TypeTracker) Resolve(op vm.Operand) ValueType {
	typ, ok := t.Get(op)

	if !ok {
		return TypeUnknown
	}

	return typ
}
