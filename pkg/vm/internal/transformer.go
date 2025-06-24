package internal

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Transformer interface {
	runtime.Value
	runtime.Iterable
	runtime.Keyed
	io.Closer

	Add(ctx context.Context, key, value runtime.Value) error
}
