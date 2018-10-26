package collections

import (
	"context"
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
	values []core.Value,
) (Iterator, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &SliceIterator{valVar, keyVar, values, 0}, nil
}

func NewDefaultSliceIterator(input []core.Value) (Iterator, error) {
	return NewSliceIterator(DefaultValueVar, DefaultKeyVar, input)
}

func (iterator *SliceIterator) Next(_ context.Context, _ *core.Scope) (DataSet, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: values.NewInt(idx),
		}, nil
	}

	return nil, nil
}
