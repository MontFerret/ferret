package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MapIterator struct {
	values map[string]core.Value
	keys   []string
	pos    int
}

func NewMapIterator(input map[string]core.Value) *MapIterator {
	return &MapIterator{input, nil, 0}
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

func (iterator *MapIterator) Next() (core.Value, core.Value, error) {
	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val := iterator.values[key]
		iterator.pos++

		return val, values.NewString(key), nil
	}

	return values.None, values.None, ErrExhausted
}
