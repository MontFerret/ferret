package events

import (
	"context"
	"time"

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
		if ctx.Err() != nil {
			return values.None, ctx.Err()
		}

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
		<-time.After(task.polling)
	}
}

func NewEvalWaitTask(
	ec *eval.Runtime,
	fn *eval.Function,
	polling time.Duration,
) *WaitTask {
	return NewWaitTask(
		func(ctx context.Context) (core.Value, error) {
			return ec.EvalValue(ctx, fn)
		},
		polling,
	)
}
