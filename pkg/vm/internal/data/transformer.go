package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Transformer interface {
	runtime.Value
	runtime.Iterable
	runtime.Keyed
	runtime.Measurable
	io.Closer

	Add(ctx context.Context, key, value runtime.Value) error
}
