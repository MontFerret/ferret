package literals

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type ArrayLiteral struct {
	elements []core.Expression
}

func NewArrayLiteral(size int) *ArrayLiteral {
	return &ArrayLiteral{make([]core.Expression, 0, size)}
}

func NewArrayLiteralWith(elements []core.Expression) *ArrayLiteral {
	return &ArrayLiteral{elements}
}

func (l *ArrayLiteral) Push(expression core.Expression) {
	l.elements = append(l.elements, expression)
}

func (l *ArrayLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	arr := values.NewArray(len(l.elements))

	for _, el := range l.elements {
		val, err := el.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		arr.Push(val)
	}

	return arr, nil
}
