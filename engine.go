package ferret

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Engine struct {
	compiler *compiler.Compiler
	host     *host
	hooks    *hookRegistry
	limiter  *sessionLimiter
	idleCap  int
	totalCap int
}

func New(setters ...Option) (*Engine, error) {
	opts, err := newOptions(setters)
	if err != nil {
		return nil, err
	}

	boot := newBootstrap(opts)

	for _, m := range opts.modules {
		if err := m.Register(boot); err != nil {
			if closeErr := boot.hooks.engine.runCloseHooks(); closeErr != nil {
				return nil, errors.Join(err, fmt.Errorf("close hooks: %w", closeErr))
			}

			return nil, err
		}
	}

	h, err := boot.host.Build()
	if err != nil {
		if closeErr := boot.hooks.engine.runCloseHooks(); closeErr != nil {
			return nil, errors.Join(err, fmt.Errorf("close hooks: %w", closeErr))
		}

		return nil, err
	}

	hooks := boot.hooks.clone()

	// Run init hooks after bootstrap is finalized and before returning the engine.
	if err := hooks.engine.runInitHooks(); err != nil {
		initErr := fmt.Errorf("init hooks: %w", err)

		if closeErr := hooks.engine.runCloseHooks(); closeErr != nil {
			return nil, errors.Join(initErr, fmt.Errorf("close hooks: %w", closeErr))
		}

		return nil, initErr
	}

	return &Engine{
		compiler: compiler.New(opts.compiler...),
		host:     h,
		hooks:    hooks,
		limiter:  newSessionLimiter(opts.maxActiveSessions),
		idleCap:  opts.maxIdleVMsPerPlan,
		totalCap: opts.maxVMsPerPlan,
	}, nil
}

func (e *Engine) Compile(ctx context.Context, src *file.Source) (*Plan, error) {
	if err := e.hooks.plan.runBeforeCompileHooks(ctx); err != nil {
		return nil, fmt.Errorf("before compile hooks: %w", err)
	}

	prog, err := e.compiler.Compile(src)

	// After-compile hooks always run and receive the compilation error (if any).
	if hookErr := e.hooks.plan.runAfterCompileHooks(ctx, err); hookErr != nil {
		return nil, errors.Join(err, fmt.Errorf("after compile hooks: %w", hookErr))
	}

	if err != nil {
		return nil, err
	}

	return &Plan{
		prog:         prog,
		host:         e.host,
		hooks:        e.hooks.plan,
		sessionHooks: e.hooks.session,
		limiter:      e.limiter,
		pool:         vm.NewPoolWithLimits(prog, e.idleCap, e.totalCap),
	}, nil
}

func (e *Engine) Run(ctx context.Context, src *file.Source, opts ...SessionOption) (*Output, error) {
	plan, err := e.Compile(ctx, src)

	if err != nil {
		return nil, err
	}

	var session *Session

	defer func() {
		logger := runtime.NewLogger(e.host.logging)

		if session != nil {
			if closeErr := session.Close(); closeErr != nil {
				logger.Error().
					Err(closeErr).
					Str("phase", "session").
					Str("operation", "close").
					Msg("deferred cleanup failed")
			}
		}

		if closeErr := plan.Close(); closeErr != nil {
			logger.Error().
				Err(closeErr).
				Str("phase", "plan").
				Str("operation", "close").
				Msg("deferred cleanup failed")
		}
	}()

	session, err = plan.NewSession(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return session.Run(ctx)
}

func (e *Engine) Close() error {
	if err := e.hooks.engine.runCloseHooks(); err != nil {
		return fmt.Errorf("close hooks: %w", err)
	}

	return nil
}
