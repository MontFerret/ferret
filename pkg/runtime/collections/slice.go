package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type SliceIterator struct {
	values []core.Value
	pos    int
}

func NewSliceIterator(input []core.Value) Iterator {
	return &SliceIterator{input, 0}
}

func (iterator *SliceIterator) HasNext() bool {
	return len(iterator.values) > iterator.pos
}

func (iterator *SliceIterator) Next() (core.Value, core.Value, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return val, values.NewInt(idx), nil
	}

	return values.None, values.None, ErrExhausted
}
