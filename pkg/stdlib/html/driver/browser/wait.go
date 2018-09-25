package browser

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"time"
)

type WaitTask struct {
	client    *cdp.Client
	predicate string
	timeout   time.Duration
	polling   time.Duration
}

const DefaultPolling = time.Millisecond * time.Duration(200)

func NewWaitTask(
	client *cdp.Client,
	predicate string,
	timeout time.Duration,
	polling time.Duration,
) *WaitTask {
	return &WaitTask{
		client,
		predicate,
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
			out, err := task.eval()

			// JS expression failed
			// terminating
			if err != nil {
				timer.Stop()

				return values.None, err
			}

			// JS output is not empty
			// terminating
			if out != values.None {
				timer.Stop()

				return out, nil
			}

			// Nothing yet, let's wait before the next try
			time.Sleep(task.polling)
		}
	}

	// TODO: Do we need this code?
	return values.None, core.ErrTimeout
}

func (task *WaitTask) eval() (core.Value, error) {
	return Eval(
		task.client,
		task.predicate,
		true,
		false,
	)
}
