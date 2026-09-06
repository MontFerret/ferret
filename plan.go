package ferret

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Plan wraps a compiled program together with the host state needed to execute it.
type Plan struct {
	prog            *bytecode.Program
	host            *host
	hooks           planHooks
	sessionHooks    sessionHooks
	limiter         *sessionLimiter
	pool            *vm.Pool
	lifecycle       creationLifecycle
	engineLifecycle *creationLifecycle
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
// acquired.
func (p *Plan) NewSession(ctx context.Context, setters ...SessionOption) (*Session, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{}, buildSession)
}

// NewDebugSession creates a retained-state source-level debugging session. All
// non-nil options are applied in order after the plan's debug metadata is
// validated, and failures from multiple options are joined before session
// resources are acquired.
func (p *Plan) NewDebugSession(ctx context.Context, setters ...SessionOption) (*DebugSession, error) {
	return newPlanSession(p, ctx, setters, planSessionSetup{requiresDebugInfo: true}, buildDebugSession)
}

// Marshal serializes the plan's compiled program using the provided program options.
func (p *Plan) Marshal(opts ...ProgramOption) ([]byte, error) {
	if err := p.lifecycle.check(); err != nil {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "plan is closed")
	}

	return MarshalProgram(p.prog, opts...)
}

// Close rejects new sessions, settles admitted construction, runs cleanup hooks,
// and closes the VM pool. Concurrent/repeated calls retain the cleanup result.
// Callers must close sessions first; this method does not close descendants.
func (p *Plan) Close() error {
	if p == nil {
		return nil
	}

	return p.lifecycle.close(func() error {
		var err error

		if hookErr := p.hooks.runCloseHooks(); hookErr != nil {
			err = fmt.Errorf("close hooks: %w", hookErr)
		}

		if poolErr := p.pool.Close(); poolErr != nil {
			err = errors.Join(err, fmt.Errorf("close pool: %w", poolErr))
		}

		return err
	})
}
