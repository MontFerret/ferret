package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Transformer interface {
	runtime.Value
	runtime.Iterable

	Add(ctx context.Context, key, value runtime.Value) error
}
