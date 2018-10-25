package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLNodeIterator struct {
	valVar string
	keyVar string
	values values.HTMLNode
	pos    int
}

func NewHTMLNodeIterator(
	valVar,
	keyVar string,
	input values.HTMLNode,
) Iterator {
	return &HTMLNodeIterator{valVar, keyVar, input, 0}
}

func (iterator *HTMLNodeIterator) HasNext() bool {
	return iterator.values.Length() > values.NewInt(iterator.pos)
}

func (iterator *HTMLNodeIterator) Next() (DataSet, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.GetChildNode(idx)

		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: idx,
		}, nil
	}

	return nil, ErrExhausted
}
