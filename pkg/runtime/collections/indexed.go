package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const (
	DefaultValueVar = "value"
	DefaultKeyVar   = "key"
)

type IndexedIterator struct {
	valVar string
	keyVar string
	values IndexedCollection
	pos    int
}

func NewIndexedIterator(
	valVar,
	keyVar string,
	values IndexedCollection,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &IndexedIterator{valVar, keyVar, values, 0}, nil
}

func NewDefaultIndexedIterator(
	values IndexedCollection,
) (Iterator, error) {
	return NewIndexedIterator(DefaultValueVar, DefaultKeyVar, values)
}

func (iterator *IndexedIterator) Next(_ context.Context, _ *core.Scope) (DataSet, error) {
	if int(iterator.values.Length()) > iterator.pos {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.Get(idx)
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: idx,
		}, nil
	}

	return nil, nil
}
