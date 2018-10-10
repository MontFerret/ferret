package events

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/html/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
)

func DispatchEvent(
	ctx context.Context,
	client *cdp.Client,
	objectId runtime.RemoteObjectID,
	eventName string,
) (values.Boolean, error) {
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

	evtID := evt.Result.ObjectID

	// release the event object
	defer client.Runtime.ReleaseObject(ctx, runtime.NewReleaseObjectArgs(*evtID))

	res, err := client.Runtime.CallFunctionOn(
		ctx,
		runtime.NewCallFunctionOnArgs("dispatchEvent").
			SetObjectID(objectId).
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
