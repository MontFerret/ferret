package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const (
	DefaultValueVar = "value"
	DefaultKeyVar   = "key"
)

type ArrayIterator struct {
	valVar string
	keyVar string
	values *values.Array
	pos    int
}

func NewArrayIterator(
	valVar,
	keyVar string,
	input *values.Array,
) Iterator {
	return &ArrayIterator{valVar, keyVar, input, 0}
}

func NewDefaultArrayIterator(
	input *values.Array,
) Iterator {
	return &ArrayIterator{DefaultValueVar, DefaultKeyVar, input, 0}
}

func (iterator *ArrayIterator) HasNext() bool {
	return int(iterator.values.Length()) > iterator.pos
}

func (iterator *ArrayIterator) Next() (DataSet, error) {
	if int(iterator.values.Length()) > iterator.pos {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.Get(idx)
		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: idx,
		}, nil
	}

	return nil, ErrExhausted
}
