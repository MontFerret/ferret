package universal

import (
	"context"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Runtime adapts a Native engine to api.Runtime. Construct an owned engine with
// New, or borrow an existing engine with Wrap.
type Runtime struct {
	native     *engine.Engine
	ownsEngine bool
}

var _ api.Runtime = (*Runtime)(nil)

// New constructs and owns a Native engine using its options. Native handles
// construction rollback; failures return a nil runtime and a projected error.
// Callers must settle work and close sessions and plans before the runtime.
func New(opts ...engine.Option) (*Runtime, error) {
	native, err := engine.New(opts...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &Runtime{native: native, ownsEngine: true}, nil
}

// Wrap borrows native without transferring ownership of it or its resources.
// The caller must settle outstanding work and close its children before native.
// Wrap panics when native is nil.
func Wrap(native *engine.Engine) *Runtime {
	if native == nil {
		panic("universal: nil native engine")
	}

	return &Runtime{native: native}
}

// Run delegates convenience execution and transient session/plan cleanup to
// Native. Available encoded output is retained alongside execution or cleanup errors.
func (r *Runtime) Run(ctx context.Context, src api.Source, setters ...api.SessionOption) (api.Output, error) {
	opts, err := newSessionOptions(ctx, setters)
	if err != nil {
		return api.Output{}, wrapDiagnosticError(err)
	}

	output, err := r.native.Run(ctx, source.New(src.Name, src.Content), opts.native...)

	return outputValue(output), wrapDiagnosticError(err)
}

// Compile translates portable compilation options and creates a reusable plan.
// Omitted optimization inherits the Native engine's default.
func (r *Runtime) Compile(ctx context.Context, src api.Source, setters ...api.PlanOption) (api.Plan, error) {
	return r.compile(ctx, src, false, setters)
}

// CompileDebug creates a reusable plan with Native debugger metadata. Explicit
// optimization must be api.OptimizationNone.
func (r *Runtime) CompileDebug(ctx context.Context, src api.Source, setters ...api.PlanOption) (api.Plan, error) {
	return r.compile(ctx, src, true, setters)
}

// Close delegates to Native for an engine owned through New. Native retains the
// cleanup result and rejects subsequent work; it does not settle outstanding work
// or close caller-owned descendants. For Wrap, Close returns nil and leaves both
// the adapter and borrowed engine usable.
func (r *Runtime) Close() error {
	if r.ownsEngine {
		return wrapDiagnosticError(r.native.Close())
	}

	return nil
}

func (r *Runtime) compile(ctx context.Context, src api.Source, debug bool, setters []api.PlanOption) (api.Plan, error) {
	opts, err := newPlanOptions(ctx, debug, setters)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	compile := r.native.Compile
	if debug {
		compile = r.native.CompileDebug
	}

	native, err := compile(ctx, source.New(src.Name, src.Content), opts.native...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return newPlan(native), nil
}
