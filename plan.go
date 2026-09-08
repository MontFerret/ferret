package ferret

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Plan wraps a compiled program together with the host state needed to execute it.
// Callers own directly created sessions and must close them before the plan.
type Plan struct {
	closeErr     error
	hooks        planHooks
	sessionHooks sessionHooks
	prog         *bytecode.Program
	host         *host
	limiter      *sessionLimiter
	pool         *vm.Pool
	closed       chan struct{}
	closeOnce    sync.Once
}

// Params returns the list of parameter names declared in the query.
func (p *Plan) Params() []string {
	if p == nil || p.prog == nil {
		return nil
	}

	// Don't expose the underlying slice to callers.
	// External mutation can corrupt the plan/program state.
	params := make([]string, len(p.prog.Params))
	copy(params, p.prog.Params)

	return params
}

// NewSession creates a session for executing the plan with optional per-run
// settings. For a valid plan, all non-nil options are applied in order, and
// failures from multiple options are joined before session resources are
// acquired. The context must be non-nil. Close wakes pending capacity acquisition;
// constructors already holding capacity may finish or fail on local pool closure.
// A returned session remains caller-owned even if Close races its construction.
func (p *Plan) NewSession(ctx context.Context, setters ...SessionOption) (*Session, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{}, buildSession)
}

// NewDebugSession creates a retained-state source-level debugging session. All
// non-nil options are applied in order after the plan's debug metadata is
// validated, and failures from multiple options are joined before session
// resources are acquired. Context, admission, and ownership follow NewSession.
func (p *Plan) NewDebugSession(ctx context.Context, setters ...SessionOption) (*DebugSession, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{requiresDebugInfo: true}, buildDebugSession)
}

// Marshal serializes the plan's compiled program using the provided program options.
// It rejects a closed plan; serialization admitted before Close may finish afterward.
func (p *Plan) Marshal(opts ...ProgramOption) ([]byte, error) {
	select {
	case <-p.closed:
		return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	default:
	}

	return MarshalProgram(p.prog, opts...)
}

// Close rejects new sessions, wakes capacity waiters, runs cleanup hooks, and
// closes the VM pool without waiting for constructors or closing borrowed VMs.
// Callers must settle outstanding work and close sessions first. Concurrent and
// repeated calls retain the cleanup result. A plan close hook must not recursively
// call Close on the same plan.
func (p *Plan) Close() error {
	if p == nil {
		return nil
	}

	p.closeOnce.Do(func() {
		close(p.closed)

		if hookErr := p.hooks.runCloseHooks(); hookErr != nil {
			p.closeErr = fmt.Errorf("close hooks: %w", hookErr)
		}

		if poolErr := p.pool.Close(); poolErr != nil {
			p.closeErr = errors.Join(p.closeErr, fmt.Errorf("close pool: %w", poolErr))
		}
	})

	return p.closeErr
}
