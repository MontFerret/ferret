package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ArrayLiteral struct {
	elements []core.Expression
}

func NewArrayLiteral(size int) *ArrayLiteral {
	return &ArrayLiteral{make([]core.Expression, 0, size)}
}

func NewArrayLiteralWith(elements ...core.Expression) *ArrayLiteral {
	return &ArrayLiteral{elements}
}

func (l *ArrayLiteral) Push(expression core.Expression) {
	l.elements = append(l.elements, expression)
}

func (l *ArrayLiteral) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	arr, err := l.doExec(ctx, scope)

	if err != nil {
		return nil, err
	}

	return collections.NewArrayIterator(arr), nil
}

func (l *ArrayLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	arr, err := l.doExec(ctx, scope)

	if err != nil {
		return values.None, err
	}

	return arr, nil
}

func (l *ArrayLiteral) doExec(ctx context.Context, scope *core.Scope) (*values.Array, error) {
	arr := values.NewArray(len(l.elements))

	for _, el := range l.elements {
		val, err := el.Exec(ctx, scope)

		if err != nil {
			return nil, err
		}

		arr.Push(val)
	}

	return arr, nil
}
