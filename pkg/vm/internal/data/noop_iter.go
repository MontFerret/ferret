package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type noopIter struct{}

func (n noopIter) Next(_ context.Context) (value runtime.Value, key runtime.Value, err error) {
	return runtime.None, runtime.None, io.EOF
}
