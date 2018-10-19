package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ObjectIterator struct {
	values *values.Object
	keys   []string
	pos    int
}

func NewObjectIterator(input *values.Object) *ObjectIterator {
	return &ObjectIterator{input, nil, 0}
}

func (iterator *ObjectIterator) HasNext() bool {
	// lazy initialization
	if iterator.keys == nil {
		iterator.keys = iterator.values.Keys()
	}

	return len(iterator.keys) > iterator.pos
}

func (iterator *ObjectIterator) Next() (ResultSet, error) {
	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val, _ := iterator.values.Get(values.NewString(key))
		iterator.pos++

		return ResultSet{val, values.NewString(key)}, nil
	}

	return nil, ErrExhausted
}
