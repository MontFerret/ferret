package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
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
	values values.HTMLNode,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &HTMLNodeIterator{valVar, keyVar, values, 0}, nil
}

func (iterator *HTMLNodeIterator) Next(_ context.Context, _ *core.Scope) (DataSet, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.GetChildNode(idx)

		iterator.pos++

		return DataSet{
			iterator.valVar: val,
			iterator.keyVar: idx,
		}, nil
	}

	return nil, nil
}
