package events

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"time"
)

type (
	Function func() (core.Value, error)
	WaitTask struct {
		fun     Function
		timeout time.Duration
		polling time.Duration
	}
)

const DefaultPolling = time.Millisecond * time.Duration(200)

func NewWaitTask(
	fun Function,
	timeout time.Duration,
	polling time.Duration,
) *WaitTask {
	return &WaitTask{
		fun,
		timeout,
		polling,
	}
}

func (task *WaitTask) Run() (core.Value, error) {
	timer := time.NewTimer(task.timeout)

	for {
		select {
		case <-timer.C:
			return values.None, core.ErrTimeout
		default:
			out, err := task.fun()

			// expression failed
			// terminating
			if err != nil {
				timer.Stop()

				return values.None, err
			}

			// output is not empty
			// terminating
			if out != values.None {
				timer.Stop()

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
	timeout time.Duration,
	polling time.Duration,
) *WaitTask {
	return NewWaitTask(
		func() (core.Value, error) {
			return eval.Eval(
				client,
				predicate,
				true,
				false,
			)
		},
		timeout,
		polling,
	)
}
