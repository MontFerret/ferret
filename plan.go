package ferret

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Plan wraps a compiled program together with the host state needed to execute it.
type Plan struct {
	prog         *bytecode.Program
	host         *host
	hooks        planHooks
	sessionHooks sessionHooks
	limiter      *sessionLimiter
	pool         *vm.Pool
	mu           sync.RWMutex
	closed       bool
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

// NewSession creates a session for executing the plan with optional per-run settings.
func (p *Plan) NewSession(ctx context.Context, setters ...SessionOption) (*Session, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{}, buildSession)
}

// NewDebugSession creates a retained-state source-level debugging session.
func (p *Plan) NewDebugSession(ctx context.Context, setters ...SessionOption) (*DebugSession, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{requiresDebugInfo: true}, buildDebugSession)
}

// Close runs plan cleanup hooks and closes the plan's VM pool.
func (p *Plan) Close() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}

	p.closed = true
	// Snapshot close dependencies before unlocking so later Session.Close calls can
	// finish independently of the Plan's mutable state.
	hooks := p.hooks
	pool := p.pool
	p.mu.Unlock()

	var err error

	// Plan close hooks follow the hook registry close semantics (LIFO with error aggregation).
	if hookErr := hooks.runCloseHooks(); hookErr != nil {
		err = fmt.Errorf("close hooks: %w", hookErr)
	}

	if poolErr := pool.Close(); poolErr != nil {
		err = errors.Join(err, fmt.Errorf("close pool: %w", poolErr))
	}

	return err
}
