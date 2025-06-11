package core

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// ConstantPool stores and deduplicates constants
type ConstantPool struct {
	values []runtime.Value
	index  map[uint64]int
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		values: make([]runtime.Value, 0),
		index:  make(map[uint64]int),
	}
}

func (cp *ConstantPool) Add(val runtime.Value) vm.Operand {
	var hash uint64
	isNone := val == runtime.None

	if runtime.IsScalar(val) {
		hash = val.Hash()
	}

	if hash > 0 || isNone {
		if idx, ok := cp.index[hash]; ok {
			return vm.NewConstantOperand(idx)
		}
	}

	cp.values = append(cp.values, val)
	idx := len(cp.values) - 1

	if hash > 0 || isNone {
		cp.index[hash] = idx
	}

	return vm.NewConstantOperand(idx)
}

func (cp *ConstantPool) Get(addr vm.Operand) runtime.Value {
	if !addr.IsConstant() {
		panic(runtime.Error(ErrInvalidOperandType, strconv.Itoa(int(addr))))
	}

	idx := addr.Constant()

	if idx < 0 || idx >= len(cp.values) {
		panic(runtime.Error(ErrConstantNotFound, strconv.Itoa(idx)))
	}

	return cp.values[idx]
}

func (cp *ConstantPool) All() []runtime.Value {
	return cp.values
}
