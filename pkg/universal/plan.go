package universal

import (
	"context"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type plan struct {
	native *engine.Plan
}

var _ api.Plan = (*plan)(nil)

func (p *plan) Params() []string {
	return p.native.Params()
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
