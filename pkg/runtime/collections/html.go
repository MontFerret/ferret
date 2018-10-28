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

func (iterator *HTMLNodeIterator) Next(_ context.Context, scope *core.Scope) (*core.Scope, error) {
	if iterator.values.Length() > values.NewInt(iterator.pos) {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.GetChildNode(idx)

		iterator.pos++

		nextScope := scope.Fork()

		if err := nextScope.SetVariable(iterator.valVar, val); err != nil {
			return nil, err
		}

		if iterator.keyVar != "" {
			if err := nextScope.SetVariable(iterator.keyVar, idx); err != nil {
				return nil, err
			}
		}

		return nextScope, nil
	}

	return nil, nil
}
