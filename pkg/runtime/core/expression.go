package core

import "context"

type (
	Expression interface {
		Exec(ctx context.Context, scope *Scope) (Value, error)
	}

	ExpressionFn struct {
		fn func(ctx context.Context, scope *Scope) (Value, error)
	}
)

func AsExpression(fn func(ctx context.Context, scope *Scope) (Value, error)) Expression {
	return &ExpressionFn{fn}
}

func (f *ExpressionFn) Exec(ctx context.Context, scope *Scope) (Value, error) {
	return f.fn(ctx, scope)
}
