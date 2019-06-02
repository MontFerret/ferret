package eval

import (
	"context"
	"fmt"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const EmptyExecutionContextID = runtime.ExecutionContextID(-1)

type ExecutionContext struct {
	client    *cdp.Client
	frame     page.Frame
	contextID runtime.ExecutionContextID
}

func NewExecutionContext(client *cdp.Client, frame page.Frame, contextID runtime.ExecutionContextID) *ExecutionContext {
	ec := new(ExecutionContext)
	ec.client = client
	ec.frame = frame
	ec.contextID = contextID

	return ec
}

func (ec *ExecutionContext) ID() runtime.ExecutionContextID {
	return ec.contextID
}

func (ec *ExecutionContext) Eval(ctx context.Context, exp string) error {
	_, err := ec.eval(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)),
	)

	return err
}

func (ec *ExecutionContext) EvalWithReturn(ctx context.Context, exp string) (core.Value, error) {
	return ec.eval(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)).
			SetReturnByValue(true),
	)
}

func (ec *ExecutionContext) EvalAsync(ctx context.Context, exp string) (core.Value, error) {
	return ec.eval(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)).
			SetReturnByValue(true).
			SetAwaitPromise(true),
	)
}

func (ec *ExecutionContext) eval(ctx context.Context, args *runtime.EvaluateArgs) (core.Value, error) {
	if ec.contextID != EmptyExecutionContextID {
		args.SetContextID(ec.contextID)
	}

	out, err := ec.client.Runtime.Evaluate(ctx, args)

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
		return values.Unmarshal(out.Result.Value)
	}

	return Unmarshal(&out.Result)
}

func (ec *ExecutionContext) CallMethod(
	ctx context.Context,
	objectID runtime.RemoteObjectID,
	methodName string,
	args []runtime.CallArgument,
) (*runtime.RemoteObject, error) {
	callArgs := runtime.NewCallFunctionOnArgs(methodName).
		SetObjectID(objectID).
		SetArguments(args)

	if ec.contextID != EmptyExecutionContextID {
		callArgs.SetExecutionContextID(ec.contextID)
	}

	found, err := ec.client.Runtime.CallFunctionOn(
		ctx,
		callArgs,
	)

	if err != nil {
		return nil, err
	}

	if found.ExceptionDetails != nil {
		return nil, found.ExceptionDetails
	}

	if found.Result.ObjectID == nil {
		return nil, nil
	}

	return &found.Result, nil
}

func (ec *ExecutionContext) ReadProperty(
	ctx context.Context,
	objectID runtime.RemoteObjectID,
	propName string,
) (core.Value, error) {
	res, err := ec.client.Runtime.GetProperties(
		ctx,
		runtime.NewGetPropertiesArgs(objectID),
	)

	if err != nil {
		return values.None, err
	}

	if res.ExceptionDetails != nil {
		return values.None, res.ExceptionDetails
	}

	// all props
	if propName == "" {
		arr := values.NewArray(len(res.Result))

		for _, prop := range res.Result {
			val, err := Unmarshal(prop.Value)

			if err != nil {
				return values.None, err
			}

			arr.Push(val)
		}

		return arr, nil
	}

	for _, prop := range res.Result {
		if prop.Name == propName {
			return Unmarshal(prop.Value)
		}
	}

	return values.None, nil
}

func (ec *ExecutionContext) DispatchEvent(
	ctx context.Context,
	objectID runtime.RemoteObjectID,
	eventName string,
) (values.Boolean, error) {
	args := runtime.NewEvaluateArgs(PrepareEval(fmt.Sprintf(`
		return new window.MouseEvent('%s', { bubbles: true, cancelable: true })
	`, eventName)))

	if ec.contextID != EmptyExecutionContextID {
		args.SetContextID(ec.contextID)
	}

	evt, err := ec.client.Runtime.Evaluate(ctx, args)

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
	defer ec.client.Runtime.ReleaseObject(ctx, runtime.NewReleaseObjectArgs(*evtID))

	_, err = ec.CallMethod(
		ctx,
		objectID,
		"dispatchEvent",
		[]runtime.CallArgument{
			{
				ObjectID: evt.Result.ObjectID,
			},
		},
	)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}
