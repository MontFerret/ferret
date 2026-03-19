package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type counterIterator struct {
	value runtime.Int
	done  bool
}

func (it *counterIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it.done {
		return runtime.None, runtime.None, io.EOF
	}

	it.done = true

	return it.value, runtime.ZeroInt, nil
}
