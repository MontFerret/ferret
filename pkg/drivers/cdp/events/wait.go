package events

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Function func(ctx context.Context) (core.Value, error)

	WaitTask struct {
		fun     Function
		polling time.Duration
	}
)

const DefaultPolling = time.Millisecond * time.Duration(200)

func NewWaitTask(
	fun Function,
	polling time.Duration,
) *WaitTask {
	return &WaitTask{
		fun,
		polling,
	}
}

func (task *WaitTask) Run(ctx context.Context) (core.Value, error) {
	for {
		select {
		case <-ctx.Done():
			return values.None, core.ErrTimeout
		default:
			out, err := task.fun(ctx)

			// expression failed
			// terminating
			if err != nil {
				return values.None, err
			}

			// output is not empty
			// terminating
			if out != values.None {
				return out, nil
			}

			// Nothing yet, let's wait before the next try
			time.Sleep(task.polling)
		}
	}
}

func NewEvalWaitTask(
	ec *eval.ExecutionContext,
	predicate string,
	polling time.Duration,
) *WaitTask {
	return NewWaitTask(
		func(ctx context.Context) (core.Value, error) {
			return ec.EvalWithReturnValue(
				ctx,
				predicate,
			)
		},
		polling,
	)
}

func NewValueWaitTask(
	when drivers.WaitEvent,
	value core.Value,
	getter Function,
	polling time.Duration,
) *WaitTask {
	return &WaitTask{
		func(ctx context.Context) (core.Value, error) {
			current, err := getter(ctx)

			if err != nil {
				return values.None, err
			}

			if when == drivers.WaitEventPresence {
				// Values appeared, exit
				if current.Compare(value) == 0 {
					// The value does not really matter if it's not None
					// None indicates that operation needs to be repeated
					return values.True, nil
				}
			} else {
				// Value disappeared, exit
				if current.Compare(value) != 0 {
					// The value does not really matter if it's not None
					// None indicates that operation needs to be repeated
					return values.True, nil
				}
			}

			return values.None, nil
		},
		polling,
	}
}
