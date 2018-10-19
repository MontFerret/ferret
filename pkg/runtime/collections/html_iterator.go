package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLNodeIterator struct {
	values values.HTMLNode
	pos    int
}

func NewHTMLNodeIterator(input values.HTMLNode) *HTMLNodeIterator {
	return &HTMLNodeIterator{input, 0}
}

func (iterator *HTMLNodeIterator) HasNext() bool {
	return iterator.values.Length() > values.NewInt(iterator.pos)
}

func (iterator *HTMLNodeIterator) Next() (ResultSet, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := iterator.pos
		val := iterator.values.GetChildNode(values.NewInt(idx))

		iterator.pos++

		return ResultSet{val, values.NewInt(idx)}, nil
	}

	return nil, ErrExhausted
}
