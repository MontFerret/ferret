package core

import "context"

type Expression interface {
	Exec(ctx context.Context, scope *Scope) (Value, error)
}
