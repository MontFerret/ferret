package uapi

import (
	"context"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type plan struct {
	native *engine.Plan
}

var _ api.Plan = (*plan)(nil)

// Params returns Native's caller-owned parameter snapshot, including after Close.
// The portable context must be non-nil and not already canceled.
func (p *plan) Params(ctx context.Context) ([]string, error) {
	if ctx == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return p.native.Params(), nil
}

func (p *plan) NewSession(ctx context.Context, setters ...api.SessionOption) (api.Session, error) {
	opts, err := newSessionOptions(setters)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	native, err := p.native.NewSession(ctx, opts.native...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &session{native: native}, nil
}

func (p *plan) NewDebugSession(ctx context.Context, setters ...api.SessionOption) (apidebugger.Session, error) {
	opts, err := newSessionOptions(setters)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	native, err := p.native.NewDebugSession(ctx, opts.native...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &debugSession{native: native}, nil
}

func (p *plan) Close() error {
	return wrapDiagnosticError(p.native.Close())
}

func newPlan(native *engine.Plan) *plan {
	return &plan{native: native}
}
