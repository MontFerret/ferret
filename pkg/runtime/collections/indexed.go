package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type IndexedIterator struct {
	valVar values.String
	keyVar values.String
	values IndexedCollection
	pos    int
}

func NewIndexedIterator(
	valVar,
	keyVar values.String,
	input IndexedCollection,
) Iterator {
	return &IndexedIterator{valVar, keyVar, input, 0}
}

func NewDefaultIndexedIterator(
	input IndexedCollection,
) Iterator {
	return &IndexedIterator{DefaultValueVar, DefaultKeyVar, input, 0}
}

func (iterator *IndexedIterator) HasNext() bool {
	return int(iterator.values.Length()) > iterator.pos
}

func (iterator *IndexedIterator) Next() (DataSet, error) {
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
