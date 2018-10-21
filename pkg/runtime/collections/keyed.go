package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type KeyedIterator struct {
	valVar values.String
	keyVar values.String
	values KeyedCollection
	keys   []string
	pos    int
}

func NewKeyedIterator(
	valVar,
	keyVar values.String,
	input KeyedCollection,
) Iterator {
	return &KeyedIterator{valVar, keyVar, input, nil, 0}
}

func NewDefaultKeyedIterator(input KeyedCollection) Iterator {
	return NewKeyedIterator(DefaultValueVar, DefaultKeyVar, input)
}

func (iterator *KeyedIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		iterator.keys = iterator.values.Keys()
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *KeyedIterator) Next() (DataSet, error) {
	if len(iterator.keys) > iterator.pos {
		key := values.NewString(iterator.keys[iterator.pos])
		val, _ := iterator.values.Get(key)
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: key,
		}, nil
	}

	return nil, ErrExhausted
}
