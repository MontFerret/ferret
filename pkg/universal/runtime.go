package universal

import (
	"context"
	"sync"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Runtime adapts a borrowed Native engine to api.Runtime. It must be constructed
// with New and must not be copied after first use.
type Runtime struct {
	native    *engine.Engine
	admission admission
	closeOnce sync.Once
}

var _ api.Runtime = (*Runtime)(nil)

// New borrows native without transferring ownership of it or its resources.
// The caller must close the adapter and its children before closing native.
// New panics when native is nil.
func New(native *engine.Engine) *Runtime {
	if native == nil {
		panic("universal: nil native engine")
	}

	return &Runtime{native: native}
}

// Run delegates convenience execution and transient session/plan cleanup to
// Native. Available encoded output is retained alongside execution or cleanup errors.
func (r *Runtime) Run(ctx context.Context, src api.Source, setters ...api.SessionOption) (api.Output, error) {
	if err := r.admission.begin(ctx, "runtime"); err != nil {
		return api.Output{}, err
	}

	defer r.admission.end()

	opts, err := newSessionOptions(ctx, setters)
	if err != nil {
		return api.Output{}, wrapDiagnosticError(err)
	}

	output, err := r.native.Run(ctx, source.New(src.Name, src.Content), opts.native...)

	return convertOutput(output), wrapDiagnosticError(err)
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

// Close rejects new calls and waits for admitted calls under their caller
// contexts. It neither cancels those calls nor closes the borrowed engine or
// published children. Concurrent and repeated calls return nil after settling.
func (r *Runtime) Close() error {
	r.closeOnce.Do(func() {
		r.admission.close()
		r.admission.wait()
	})

	return nil
}

func (r *Runtime) compile(ctx context.Context, src api.Source, debug bool, setters []api.PlanOption) (api.Plan, error) {
	if err := r.admission.begin(ctx, "runtime"); err != nil {
		return nil, err
	}

	defer r.admission.end()

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
