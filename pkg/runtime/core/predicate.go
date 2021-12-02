package core

import "context"

type Predicate interface {
	Expression
	Eval(ctx context.Context, left, right Value) (Value, error)
}
