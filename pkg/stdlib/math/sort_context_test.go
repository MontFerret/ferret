package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Cancel during comparator polling without depending on scheduler timing.
	sortCancelContext struct {
		context.Context
		cancel    context.CancelFunc
		remaining int
	}

	// Fail inside runtime comparison, optionally canceling the same operation.
	sortErrorValue struct {
		runtime.Value
		err    error
		cancel context.CancelFunc
	}
)

func (ctx *sortCancelContext) Err() error {
	ctx.remaining--
	if ctx.remaining == 0 {
		ctx.cancel()
	}

	return ctx.Context.Err()
}

func (value *sortErrorValue) Compare(context.Context, runtime.Value) (runtime.Ordering, error) {
	if value.cancel != nil {
		value.cancel()
	}

	return runtime.Equal, value.err
}
