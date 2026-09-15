package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// A fallible comparison can cancel the caller independently of its result.
type sortErrorValue struct {
	runtime.Value
	err    error
	cancel context.CancelFunc
}

func (value *sortErrorValue) Compare(context.Context, runtime.Value) (runtime.Ordering, error) {
	if value.cancel != nil {
		value.cancel()
	}

	return runtime.Equal, value.err
}
