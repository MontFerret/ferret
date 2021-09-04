package eval

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/rs/zerolog"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const EmptyExecutionContextID = runtime.ExecutionContextID(-1)

type Runtime struct {
	logger    zerolog.Logger
	client    *cdp.Client
	frame     page.Frame
	contextID runtime.ExecutionContextID
}

func New(ctx context.Context, logger zerolog.Logger, client *cdp.Client, frameID page.FrameID) (*Runtime, error) {
	world, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(frameID))

	if err != nil {
		return nil, err
	}

	return Create(logger, client, world.ExecutionContextID), nil
}

func Create(logger zerolog.Logger, client *cdp.Client, contextID runtime.ExecutionContextID) *Runtime {
	rt := new(Runtime)
	rt.logger = logging.WithName(logger.With(), "js-eval").Logger()
	rt.client = client
	rt.contextID = contextID

	return rt
}

func (rt *Runtime) ContextID() runtime.ExecutionContextID {
	return rt.contextID
}

func (rt *Runtime) Eval(ctx context.Context, fn *Function) error {
	_, err := rt.call(ctx, fn)

	return err
}

func (rt *Runtime) EvalValue(ctx context.Context, fn *Function) (core.Value, error) {
	out, err := rt.call(ctx, fn.returnValue())

	if err != nil {
		return values.None, err
	}

	return CastToValue(out)
}

func (rt *Runtime) EvalRef(ctx context.Context, fn *Function) (runtime.RemoteObject, error) {
	out, err := rt.call(ctx, fn.returnRef())

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	return CastToReference(out)
}

func (rt *Runtime) ReadProperty(
	ctx context.Context,
	objectID runtime.RemoteObjectID,
	propName string,
) (core.Value, error) {
	res, err := rt.client.Runtime.GetProperties(
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
			if prop.Value != nil {
				val, err := Unmarshal(*prop.Value)

				if err != nil {
					return values.None, err
				}

				arr.Push(val)
			}
		}

		return arr, nil
	}

	for _, prop := range res.Result {
		if prop.Name == propName {
			if prop.Value != nil {
				return Unmarshal(*prop.Value)
			}

			return values.None, nil
		}
	}

	return values.None, nil
}

func (rt *Runtime) call(ctx context.Context, fn *Function) (interface{}, error) {
	expression := fn.String()
	rt.logger.Trace().Str("expression", expression).Msg("executing an expression...")
	repl, err := rt.client.Runtime.CallFunctionOn(ctx, fn.build(rt.contextID))

	if err != nil {
		rt.logger.Trace().Err(err).Str("expression", expression).Msg("failed to execute an expression")

		return nil, errors.Wrap(err, "runtime call")
	}

	if err := parseRuntimeException(repl.ExceptionDetails); err != nil {
		rt.logger.Trace().Err(err).Str("expression", expression).Msg("expression execution has failed")

		return nil, err
	}

	var className string

	if repl.Result.ClassName != nil {
		className = *repl.Result.ClassName
	}

	var subtype string

	if repl.Result.Subtype != nil {
		subtype = *repl.Result.Subtype
	}

	rt.logger.Trace().
		Str("expression", expression).
		Str("return-type", repl.Result.Type).
		Str("return-sub-type", subtype).
		Str("return-class-name", className).
		Str("return-value", string(repl.Result.Value)).
		Msg("succeeded to executed an expression")

	switch fn.returnType {
	case ReturnValue:
		return Unmarshal(repl.Result)
	case ReturnRef:
		return repl.Result, nil
	default:
		return nil, nil
	}
}
