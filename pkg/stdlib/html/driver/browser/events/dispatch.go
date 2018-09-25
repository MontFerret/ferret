package events

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser/eval"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
)

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

	evt, err := client.Runtime.Evaluate(ctx, runtime.NewEvaluateArgs(eval.PrepareEval(fmt.Sprintf(`
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
