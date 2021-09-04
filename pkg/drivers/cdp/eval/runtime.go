package eval

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const EmptyExecutionContextID = runtime.ExecutionContextID(-1)

type Runtime struct {
	client    *cdp.Client
	frame     page.Frame
	contextID runtime.ExecutionContextID
}

func New(ctx context.Context, client *cdp.Client, frameID page.FrameID) (*Runtime, error) {
	world, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(frameID))

	if err != nil {
		return nil, err
	}

	return Create(client, world.ExecutionContextID), nil
}

func Create(client *cdp.Client, contextID runtime.ExecutionContextID) *Runtime {
	ec := new(Runtime)
	ec.client = client
	ec.contextID = contextID

	return ec
}

func (ex *Runtime) ContextID() runtime.ExecutionContextID {
	return ex.contextID
}

func (ex *Runtime) Eval(ctx context.Context, fn *Function) error {
	_, err := ex.call(ctx, fn)

	return err
}

func (ex *Runtime) EvalValue(ctx context.Context, fn *Function) (core.Value, error) {
	out, err := ex.call(ctx, fn.returnValue())

	if err != nil {
		return values.None, err
	}

	return CastToValue(out)
}

func (ex *Runtime) EvalRef(ctx context.Context, fn *Function) (runtime.RemoteObject, error) {
	out, err := ex.call(ctx, fn.returnRef())

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	return CastToReference(out)
}

func (ex *Runtime) ReadProperty(
	ctx context.Context,
	objectID runtime.RemoteObjectID,
	propName string,
) (core.Value, error) {
	res, err := ex.client.Runtime.GetProperties(
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

func (ex *Runtime) call(ctx context.Context, fn *Function) (interface{}, error) {
	repl, err := ex.client.Runtime.CallFunctionOn(ctx, fn.build(ex.contextID))

	if err != nil {
		return nil, errors.Wrap(err, "runtime call")
	}

	if err := parseRuntimeException(repl.ExceptionDetails); err != nil {
		return nil, err
	}

	switch fn.returnType {
	case ReturnValue:
		out := repl.Result

		if out.Type != "undefined" && out.Type != "null" {
			return values.Unmarshal(out.Value)
		}

		return Unmarshal(&out)
	case ReturnRef:
		return repl.Result, nil
	default:
		return nil, nil
	}
}
