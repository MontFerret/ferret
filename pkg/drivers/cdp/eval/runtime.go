package eval

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const EmptyExecutionContextID = runtime.ExecutionContextID(-1)

type Runtime struct {
	logger    zerolog.Logger
	client    *cdp.Client
	frame     page.Frame
	contextID runtime.ExecutionContextID
	resolver  *Resolver
}

func Create(
	ctx context.Context,
	logger zerolog.Logger,
	client *cdp.Client,
	frameID page.FrameID,
) (*Runtime, error) {
	world, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(frameID))

	if err != nil {
		return nil, err
	}

	return New(logger, client, world.ExecutionContextID), nil
}

func New(
	logger zerolog.Logger,
	client *cdp.Client,
	contextID runtime.ExecutionContextID,
) *Runtime {
	rt := new(Runtime)
	rt.logger = logging.WithName(logger.With(), "js-eval").Logger()
	rt.client = client
	rt.contextID = contextID
	rt.resolver = NewResolver(client.Runtime)

	return rt
}

func (rt *Runtime) SetLoader(loader ValueLoader) *Runtime {
	rt.resolver.SetLoader(loader)

	return rt
}

func (rt *Runtime) ContextID() runtime.ExecutionContextID {
	return rt.contextID
}

func (rt *Runtime) Eval(ctx context.Context, fn *Function) error {
	_, err := rt.call(ctx, fn)

	return err
}

func (rt *Runtime) EvalRef(ctx context.Context, fn *Function) (runtime.RemoteObject, error) {
	out, err := rt.call(ctx, fn.returnRef())

	if err != nil {
		return runtime.RemoteObject{}, err
	}

	return out, nil
}

func (rt *Runtime) EvalValue(ctx context.Context, fn *Function) (core.Value, error) {
	out, err := rt.call(ctx, fn.returnValue())

	if err != nil {
		return values.None, err
	}

	return rt.resolver.ToValue(ctx, out)
}

func (rt *Runtime) EvalElement(ctx context.Context, fn *Function) (drivers.HTMLElement, error) {
	ref, err := rt.EvalRef(ctx, fn)

	if err != nil {
		return nil, err
	}

	return rt.resolver.ToElement(ctx, ref)
}

func (rt *Runtime) EvalElements(ctx context.Context, fn *Function) (*values.Array, error) {
	ref, err := rt.EvalRef(ctx, fn)

	if err != nil {
		return nil, err
	}

	val, err := rt.resolver.ToValue(ctx, ref)

	if err != nil {
		return nil, err
	}

	arr, ok := val.(*values.Array)

	if ok {
		return arr, nil
	}

	return values.NewArrayWith(val), nil
}

func (rt *Runtime) call(ctx context.Context, fn *Function) (runtime.RemoteObject, error) {
	log := rt.logger.With().
		Str("expression", fn.String()).
		Str("returns", fn.returnType.String()).
		Bool("is-async", fn.async).
		Str("owner", string(fn.ownerID)).
		Array("arguments", fn.args).
		Logger()

	log.Trace().Msg("executing expression...")

	repl, err := rt.client.Runtime.CallFunctionOn(ctx, fn.build(rt.contextID))

	if err != nil {
		log.Trace().Err(err).Msg("failed executing expression")

		return runtime.RemoteObject{}, errors.Wrap(err, "runtime call")
	}

	if err := parseRuntimeException(repl.ExceptionDetails); err != nil {
		log.Trace().Err(err).Msg("expression has failed with runtime exception")

		return runtime.RemoteObject{}, err
	}

	var className string

	if repl.Result.ClassName != nil {
		className = *repl.Result.ClassName
	}

	var subtype string

	if repl.Result.Subtype != nil {
		subtype = *repl.Result.Subtype
	}

	log.Trace().
		Str("return-type", repl.Result.Type).
		Str("return-sub-type", subtype).
		Str("return-class-name", className).
		Str("return-value", string(repl.Result.Value)).
		Msg("succeeded executing expression")

	return repl.Result, nil
}
