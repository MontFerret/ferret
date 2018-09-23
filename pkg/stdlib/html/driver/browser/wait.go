package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
	"time"
)

type WaitTask struct {
	client    *cdp.Client
	predicate string
	timeout   time.Duration
}

func NewWaitTask(
	client *cdp.Client,
	predicate string,
	timeout time.Duration,
) *WaitTask {
	return &WaitTask{
		client,
		fmt.Sprintf("((function () {%s})())", predicate),
		timeout,
	}
}

func (task *WaitTask) Run() (core.Value, error) {
	var result core.Value = values.None
	var err error
	var done bool
	timer := time.NewTimer(task.timeout)

	for !done {
		select {
		case <-timer.C:
			err = core.ErrTimeout
			done = true
		default:
			out, e := task.exec()

			if e != nil {
				done = true
				timer.Stop()
				err = e

				break
			}

			if out != values.None {
				timer.Stop()

				result = out
				done = true
				break
			}
		}
	}

	return result, err
}

func (task *WaitTask) exec() (core.Value, error) {
	args := runtime.NewEvaluateArgs(task.predicate).SetReturnByValue(true)
	out, err := task.client.Runtime.Evaluate(context.Background(), args)

	if err != nil {
		return values.None, err
	}

	if out.ExceptionDetails != nil {
		ex := out.ExceptionDetails
		return values.None, core.Error(
			core.ErrUnexpected,
			fmt.Sprintf("%s %s", ex.Text, *ex.Exception.Description),
		)
	}

	if out.Result.Type != "undefined" {
		var o interface{}

		err := json.Unmarshal(out.Result.Value, &o)

		if err != nil {
			return values.None, core.Error(core.ErrUnexpected, err.Error())
		}

		return values.Parse(o), nil
	}

	return values.None, nil
}
