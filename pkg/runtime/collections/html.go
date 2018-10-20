package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLNodeIterator struct {
	values values.HTMLNode
	pos    int
}

func NewHTMLNodeIterator(input values.HTMLNode) Iterator {
	return &HTMLNodeIterator{input, 0}
}

func (iterator *HTMLNodeIterator) HasNext() bool {
	return iterator.values.Length() > values.NewInt(iterator.pos)
}

func (iterator *HTMLNodeIterator) Next() (core.Value, core.Value, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := iterator.pos
		val := iterator.values.GetChildNode(values.NewInt(idx))

		iterator.pos++

		return val, values.NewInt(idx), nil
	}

	return values.None, values.None, ErrExhausted
}
