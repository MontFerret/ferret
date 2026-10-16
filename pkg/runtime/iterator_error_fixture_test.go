package runtime_test

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	failedIterator  struct{ err error }
	exhaustionError struct{}
)

func (it failedIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.None, runtime.None, it.err
}

func (exhaustionError) Error() string        { return "exhausted" }
func (exhaustionError) Is(target error) bool { return target == io.EOF }
