package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"time"
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
	client *cdp.Client,
	predicate string,
	polling time.Duration,
) *WaitTask {
	return NewWaitTask(
		func(ctx context.Context) (core.Value, error) {
			return eval.Eval(
				ctx,
				client,
				predicate,
				true,
				false,
			)
		},
		polling,
	)
}
