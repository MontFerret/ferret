package core

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// ConstantPool stores and deduplicates constants
	ConstantPool struct {
		index  map[uint64]constantBucket
		values []runtime.Value
	}

	constantBucket struct {
		rest  []int
		first int
	}
)

func NewConstantPool() *ConstantPool {
	return &ConstantPool{
		values: make([]runtime.Value, 0),
		index:  make(map[uint64]constantBucket),
	}
}

func (cp *ConstantPool) Add(val runtime.Value) bytecode.Operand {
	var hash uint64
	isNone := val == runtime.None

	if runtime.IsScalar(val) {
		hash = val.Hash()
	}

	if hash > 0 || isNone {
		if bucket, ok := cp.index[hash]; ok {
			if cp.valueEqualsAt(bucket.first, val) {
				return bytecode.NewConstant(bucket.first)
			}

			for _, idx := range bucket.rest {
				if cp.valueEqualsAt(idx, val) {
					return bytecode.NewConstant(idx)
				}
			}
		}
	}

	cp.values = append(cp.values, val)
	idx := len(cp.values) - 1

	if hash > 0 || isNone {
		bucket, exists := cp.index[hash]
		if !exists {
			cp.index[hash] = constantBucket{first: idx}
		} else {
			bucket.rest = append(bucket.rest, idx)
			cp.index[hash] = bucket
		}
	}

	return bytecode.NewConstant(idx)
}

func (cp *ConstantPool) valueEqualsAt(index int, right runtime.Value) bool {
	left := cp.values[index]
	if left == runtime.None || right == runtime.None {
		return left == runtime.None && right == runtime.None
	}

	if runtime.TypeOf(left) != runtime.TypeOf(right) {
		return false
	}

	equal, err := runtime.EqualValues(context.Background(), left, right)
	return err == nil && bool(equal)
}

func (cp *ConstantPool) Get(addr bytecode.Operand) runtime.Value {
	if !addr.IsConstant() {
		PanicInvariantf("invalid operand used in the constant pool: %s", addr)
	}

	idx := addr.Constant()

	if idx < 0 || idx >= len(cp.values) {
		PanicInvariantf("constant operand out of range: %s", addr)
	}

	return cp.values[idx]
}

func (cp *ConstantPool) All() []runtime.Value {
	return cp.values
}
