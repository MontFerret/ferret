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

func (iterator *SliceIterator) Next(_ context.Context, scope *core.Scope) (*core.Scope, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]

		iterator.pos++

		cs := scope.Fork()

		if err := cs.SetVariable(iterator.valVar, val); err != nil {
			return nil, err
		}

		if iterator.keyVar != "" {
			if err := cs.SetVariable(iterator.keyVar, values.NewInt(idx)); err != nil {
				return nil, err
			}
		}

		return cs, nil
	}

	return nil, nil
}
