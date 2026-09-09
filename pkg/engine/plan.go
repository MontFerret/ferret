package engine

import (
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Plan wraps a compiled program together with the host state needed to execute it.
// Callers own directly created sessions and must close them before the plan.
type Plan struct {
	closeErr     error
	hooks        *host.PlanHooks
	sessionHooks *host.SessionHooks
	prog         *bytecode.Program
	host         *host.Host
	limiter      *enginesession.Limiter
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

// Marshal serializes the plan's compiled program using the provided program options.
// It rejects a closed plan; serialization admitted before Close may finish afterward.
func (p *Plan) Marshal(opts ...artifact.Option) ([]byte, error) {
	select {
	case <-p.closed:
		return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	default:
	}

	return artifact.Marshal(p.prog, opts...)
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

		if hookErr := p.hooks.RunClose(); hookErr != nil {
			p.closeErr = fmt.Errorf("close hooks: %w", hookErr)
		}

		if poolErr := p.pool.Close(); poolErr != nil {
			p.closeErr = errors.Join(p.closeErr, fmt.Errorf("close pool: %w", poolErr))
		}
	})

	return p.closeErr
}
