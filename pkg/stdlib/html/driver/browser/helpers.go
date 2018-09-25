package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"golang.org/x/sync/errgroup"
)

func PointerInt(input int) *int {
	return &input
}

type BatchFunc = func() error

func RunBatch(funcs ...BatchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}

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

func DispatchEvent(
	ctx context.Context,
	client *cdp.Client,
	id dom.NodeID,
	eventName string,
) (values.Boolean, error) {
	// get a ref to remote object representing the node
	obj, err := client.DOM.ResolveNode(
		ctx,
		dom.NewResolveNodeArgs().
			SetNodeID(id),
	)

	if err != nil {
		return values.False, err
	}

	if obj.Object.ObjectID == nil {
		return values.False, nil
	}

	evt, err := client.Runtime.Evaluate(ctx, runtime.NewEvaluateArgs(PrepareEval(fmt.Sprintf(`
		return new window.MouseEvent('%s', { bubbles: true })
	`, eventName))))

	if err != nil {
		return values.False, nil
	}

	if evt.ExceptionDetails != nil {
		return values.False, evt.ExceptionDetails
	}

	if evt.Result.ObjectID == nil {
		return values.False, nil
	}

	evtId := evt.Result.ObjectID

	// release the event object
	defer client.Runtime.ReleaseObject(ctx, runtime.NewReleaseObjectArgs(*evtId))

	res, err := client.Runtime.CallFunctionOn(
		ctx,
		runtime.NewCallFunctionOnArgs("dispatchEvent").
			SetObjectID(*obj.Object.ObjectID).
			SetArguments([]runtime.CallArgument{
				{
					ObjectID: evt.Result.ObjectID,
				},
			}),
	)

	if err != nil {
		return values.False, err
	}

	if res.ExceptionDetails != nil {
		return values.False, res.ExceptionDetails
	}

	return values.True, nil
}
