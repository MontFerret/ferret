package engine

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	enginesession "github.com/MontFerret/ferret/v2/pkg/engine/internal/session"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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
