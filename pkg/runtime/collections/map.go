package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MapIterator struct {
	valVar string
	keyVar string
	values map[string]core.Value
	keys   []string
	pos    int
}

func NewMapIterator(
	valVar,
	keyVar string,
	input map[string]core.Value,
) Iterator {
	return &MapIterator{valVar, keyVar, input, nil, 0}
}

func (iterator *MapIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		keys := make([]string, len(iterator.values))

		i := 0
		for k := range iterator.values {
			keys[i] = k
			i++
		}

		iterator.keys = keys
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *MapIterator) Next() (DataSet, error) {
	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val := iterator.values[key]
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: values.NewString(key),
		}, nil
	}

	return nil, ErrExhausted
}
