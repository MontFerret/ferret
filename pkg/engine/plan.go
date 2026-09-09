package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
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

// NewSession creates a session for executing the plan with optional per-run
// settings. For a valid plan, all non-nil options are applied in order, and
// failures from multiple options are joined before session resources are
// acquired. The context must be non-nil. Close wakes pending capacity acquisition;
// constructors already holding capacity may finish or fail on local pool closure.
// A returned session remains caller-owned even if Close races its construction.
func (p *Plan) NewSession(ctx context.Context, setters ...SessionOption) (*Session, error) {
	opts, err := p.sessionOptions(ctx, false, setters)
	if err != nil {
		return nil, err
	}

	execution, err := enginesession.NewExecution(ctx, opts.config, p.host, p.sessionHooks, p.limiter, p.pool, p.closed)
	if err != nil {
		return nil, err
	}

	session := &Session{execution: execution, hooks: p.sessionHooks}
	// Construction transferred ownership. Only request cancellation rolls back
	// the completed session; parent closure does not revoke that ownership.
	if err := ctx.Err(); err != nil {
		if closeErr := session.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return nil, err
	}

	return session, nil
}

// NewDebugSession creates a retained-state source-level debugging session. All
// non-nil options are applied in order after the plan's debug metadata is
// validated, and failures from multiple options are joined before session
// resources are acquired. Context, admission, and ownership follow NewSession.
func (p *Plan) NewDebugSession(ctx context.Context, setters ...SessionOption) (*debugger.Session, error) {
	opts, err := p.sessionOptions(ctx, true, setters)
	if err != nil {
		return nil, err
	}

	session, err := enginesession.NewDebugSession(ctx, opts.config, p.host, p.sessionHooks, p.limiter, p.prog, p.closed)
	if err != nil {
		return nil, err
	}

	if err := ctx.Err(); err != nil {
		if closeErr := session.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return nil, err
	}

	return session, nil
}

func (p *Plan) sessionOptions(ctx context.Context, debug bool, setters []SessionOption) (sessionConfig, error) {
	if p == nil {
		return sessionConfig{}, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	if ctx == nil {
		return sessionConfig{}, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return sessionConfig{}, err
	}

	select {
	case <-p.closed:
		return sessionConfig{}, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	default:
	}

	if debug && len(p.prog.Metadata.DebugPoints) == 0 {
		return sessionConfig{}, runtime.Error(runtime.ErrInvalidOperation, "plan was not compiled for debugging")
	}

	opts, err := newSessionConfig(setters)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return sessionConfig{}, err
	}

	if err := ctx.Err(); err != nil {
		return sessionConfig{}, err
	}

	return opts, nil
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
