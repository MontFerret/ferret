package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type SliceIterator struct {
	valVar string
	keyVar string
	values []core.Value
	pos    int
}

func NewSliceIterator(
	valVar,
	keyVar string,
	input []core.Value,
) Iterator {
	return &SliceIterator{valVar, keyVar, input, 0}
}

func (iterator *SliceIterator) HasNext() bool {
	return len(iterator.values) > iterator.pos
}

func (iterator *SliceIterator) Next() (DataSet, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: values.NewInt(idx),
		}, nil
	}

	return nil, ErrExhausted
}
