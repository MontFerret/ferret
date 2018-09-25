package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
)

func PrepareEval(exp string) string {
	return fmt.Sprintf("((function () {%s})())", exp)
}

func Eval(client *cdp.Client, exp string, ret bool, async bool) (core.Value, error) {
	args := runtime.
		NewEvaluateArgs(PrepareEval(exp)).
		SetReturnByValue(ret).
		SetAwaitPromise(async)

	out, err := client.Runtime.Evaluate(context.Background(), args)

	if err != nil {
		return values.None, err
	}

	if out.ExceptionDetails != nil {
		ex := out.ExceptionDetails

		return values.None, core.Error(
			core.ErrUnexpected,
			fmt.Sprintf("%s: %s", ex.Text, *ex.Exception.Description),
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
