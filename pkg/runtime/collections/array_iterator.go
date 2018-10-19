package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ArrayIterator struct {
	values *values.Array
	pos    int
}

func NewArrayIterator(input *values.Array) *ArrayIterator {
	return &ArrayIterator{input, 0}
}

func (iterator *ArrayIterator) HasNext() bool {
	return int(iterator.values.Length()) > iterator.pos
}

func (iterator *ArrayIterator) Next() (ResultSet, error) {
	if int(iterator.values.Length()) > iterator.pos {
		idx := iterator.pos
		val := iterator.values.Get(values.NewInt(idx))
		iterator.pos++

		return ResultSet{val, values.NewInt(idx)}, nil
	}

	return nil, ErrExhausted
}
