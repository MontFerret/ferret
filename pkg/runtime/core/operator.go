package core

import "context"

type OperatorExpression interface {
	Expression
	Eval(ctx context.Context, left, right Value) (Value, error)
}
