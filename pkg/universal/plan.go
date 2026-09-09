package universal

import (
	"context"
	"sync"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type plan struct {
	native    *engine.Plan
	closing   context.Context
	cancel    context.CancelFunc
	closeErr  error
	admission admission
	closeOnce sync.Once
}

var _ api.Plan = (*plan)(nil)

func (p *plan) Params() []string {
	return p.native.Params()
}

func (p *plan) NewSession(ctx context.Context, setters ...api.SessionOption) (api.Session, error) {
	construction, done, err := p.begin(ctx)
	if err != nil {
		return nil, err
	}

	defer done()

	opts, err := newSessionOptions(construction, setters)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	native, err := p.native.NewSession(construction, opts.native...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &session{native: native}, nil
}

func (p *plan) NewDebugSession(ctx context.Context, setters ...api.SessionOption) (apidebugger.Session, error) {
	construction, done, err := p.begin(ctx)
	if err != nil {
		return nil, err
	}

	defer done()

	opts, err := newSessionOptions(construction, setters)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	native, err := p.native.NewDebugSession(construction, opts.native...)
	if err != nil {
		return nil, wrapDiagnosticError(err)
	}

	return &debugSession{native: native}, nil
}

func (p *plan) Close() error {
	p.closeOnce.Do(func() {
		p.admission.close()
		// Wake capacity waiters before waiting, without closing the Native pool
		// or running its hooks while admitted constructors still use it.
		p.cancel()
		p.admission.wait()
		p.closeErr = wrapDiagnosticError(p.native.Close())
	})

	return p.closeErr
}

func (p *plan) begin(ctx context.Context) (context.Context, func(), error) {
	if err := p.admission.begin(ctx, "plan"); err != nil {
		return nil, nil, err
	}

	construction, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(p.closing, cancel)

	return construction, func() {
		// Native uses this context only for construction. Run/Start establish
		// their own lifetimes, so published sessions never inherit plan closure.
		stop()
		cancel()
		p.admission.end()
	}, nil
}

func newPlan(native *engine.Plan) *plan {
	closing, cancel := context.WithCancel(context.Background())

	return &plan{native: native, closing: closing, cancel: cancel}
}
