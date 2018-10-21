package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ObjectIterator struct {
	valVar string
	keyVar string
	values *values.Object
	keys   []string
	pos    int
}

func NewObjectIterator(
	valVar,
	keyVar string,
	input *values.Object,
) Iterator {
	return &ObjectIterator{valVar, keyVar, input, nil, 0}
}

func NewDefaultObjectIterator(input *values.Object) Iterator {
	return NewObjectIterator(DefaultValueVar, DefaultKeyVar, input)
}

func (iterator *ObjectIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		iterator.keys = iterator.values.Keys()
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *ObjectIterator) Next() (DataSet, error) {
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
