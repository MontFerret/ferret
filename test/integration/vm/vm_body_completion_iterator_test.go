package vm_test

import (
	"context"
	"io"
	"sync/atomic"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type bodyCompletionIterator struct {
	closeCount atomic.Int32
	index      int
}

func (i *bodyCompletionIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	if i.index >= 2 {
		return runtime.None, runtime.None, io.EOF
	}

	i.index++

	return runtime.NewInt(i.index), runtime.NewInt(i.index - 1), nil
}

func (i *bodyCompletionIterator) Close() error {
	i.closeCount.Add(1)

	return nil
}
