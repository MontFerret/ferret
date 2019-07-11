package eval

import (
	"context"
	"fmt"
	"github.com/mafredri/cdp/protocol/dom"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"

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
	_, err := ec.evalWithValueInternal(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)),
	)

	return err
}

func (ec *ExecutionContext) EvalWithValue(ctx context.Context, exp string) (core.Value, error) {
	return ec.evalWithValueInternal(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)).
			SetReturnByValue(true),
	)
}

func (ec *ExecutionContext) EvalAsync(ctx context.Context, exp string) (core.Value, error) {
	return ec.evalWithValueInternal(
		ctx,
		runtime.
			NewEvaluateArgs(PrepareEval(exp)).
			SetReturnByValue(true).
			SetAwaitPromise(true),
	)
}

func (ec *ExecutionContext) ResolveRemoteObject(ctx context.Context, exp string) (runtime.RemoteObject, error) {
	res, err := ec.evalInternal(ctx, runtime.NewEvaluateArgs(PrepareEval(exp)))

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	if res.ObjectID == nil {
		return runtime.RemoteObject{}, errors.Wrap(core.ErrUnexpected, "unable to resolve remote object")
	}

	return res, nil
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

func (ec *ExecutionContext) ReadPropertyByNodeID(
	ctx context.Context,
	nodeID dom.NodeID,
	propName string,
) (core.Value, error) {
	obj, err := ec.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(nodeID))

	if err != nil {
		return values.None, err
	}

	if obj.Object.ObjectID == nil {
		return values.None, nil
	}

	return ec.ReadProperty(ctx, *obj.Object.ObjectID, propName)
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

func (ec *ExecutionContext) CallFunction(ctx context.Context, declaration string, args ...runtime.CallArgument) (runtime.RemoteObject, error) {
	cfArgs := runtime.NewCallFunctionOnArgs(declaration).SetArguments(args)

	if ec.contextID != EmptyExecutionContextID {
		cfArgs.SetExecutionContextID(ec.contextID)
	}

	repl, err := ec.client.Runtime.CallFunctionOn(ctx, cfArgs)

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	if repl.ExceptionDetails != nil {
		exception := *repl.ExceptionDetails

		return runtime.RemoteObject{}, errors.New(exception.Error())
	}

	return repl.Result, nil
}

func (ec *ExecutionContext) evalWithValueInternal(ctx context.Context, args *runtime.EvaluateArgs) (core.Value, error) {
	obj, err := ec.evalInternal(ctx, args)

	if err != nil {
		return values.None, err
	}

	if obj.Type != "undefined" && obj.Type != "null" {
		return values.Unmarshal(obj.Value)
	}

	return Unmarshal(&obj)
}

func (ec *ExecutionContext) evalInternal(ctx context.Context, args *runtime.EvaluateArgs) (runtime.RemoteObject, error) {
	if ec.contextID != EmptyExecutionContextID {
		args.SetContextID(ec.contextID)
	}

	out, err := ec.client.Runtime.Evaluate(ctx, args)

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	if out.ExceptionDetails != nil {
		ex := out.ExceptionDetails

		return runtime.RemoteObject{}, core.Error(
			core.ErrUnexpected,
			fmt.Sprintf("%s: %s", ex.Text, *ex.Exception.Description),
		)
	}

	return out.Result, nil
}
