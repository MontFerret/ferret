package core

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ConstantPool stores and deduplicates constants
type ConstantPool struct {
	index  map[uint64]int
	values []runtime.Value
}

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		values: make([]runtime.Value, 0),
		index:  make(map[uint64]int),
	}
}

func (cp *ConstantPool) Add(val runtime.Value) bytecode.Operand {
	var hash uint64
	isNone := val == runtime.None

	if runtime.IsScalar(val) {
		hash = val.Hash()
	}

	if hash > 0 || isNone {
		if idx, ok := cp.index[hash]; ok {
			return bytecode.NewConstant(idx)
		}
	}

	cp.values = append(cp.values, val)
	idx := len(cp.values) - 1

	if hash > 0 || isNone {
		cp.index[hash] = idx
	}

	return bytecode.NewConstant(idx)
}

func (cp *ConstantPool) Get(addr bytecode.Operand) runtime.Value {
	if !addr.IsConstant() {
		panic(fmt.Errorf("invalid operand type used in the constant pool: %s", addr))
	}

	idx := addr.Constant()

	if idx < 0 || idx >= len(cp.values) {
		panic(fmt.Errorf("invalid operand type used in the constant pool: %s", addr))
	}

	return cp.values[idx]
}

func (cp *ConstantPool) All() []runtime.Value {
	return cp.values
}
