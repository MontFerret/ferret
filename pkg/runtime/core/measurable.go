package core

import "context"

type Measurable interface {
	Length(ctx context.Context) (Int, error)
}
