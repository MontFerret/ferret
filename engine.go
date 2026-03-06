package ferret

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

type Engine struct {
	compiler *compiler.Compiler
	host     *host
	hooks    *hookRegistry
}

func New(setters ...Option) (*Engine, error) {
	opts, err := newOptions(setters)
	if err != nil {
		return nil, err
	}

	boot := newBootstrap(opts)

	for _, m := range opts.modules {
		if err := m.Register(boot); err != nil {
			return nil, err
		}
	}

	h, err := boot.host.Build()
	if err != nil {
		return nil, err
	}

	hooks := boot.hooks.clone()

	if err := hooks.engine.runInitHooks(); err != nil {
		return nil, fmt.Errorf("init hooks: %w", err)
	}

	return &Engine{
		compiler: compiler.New(opts.compiler...),
		host:     h,
		hooks:    hooks,
	}, nil
}

func (e *Engine) Compile(ctx context.Context, src *file.Source) (*Plan, error) {
	if err := e.hooks.plan.runBeforeCompileHooks(ctx); err != nil {
		return nil, fmt.Errorf("before compile hooks: %w", err)
	}

	prog, err := e.compiler.Compile(src)

	if hookErr := e.hooks.plan.runAfterCompileHooks(ctx, err); hookErr != nil {
		return nil, errors.Join(err, fmt.Errorf("after compile hooks: %w", err))
	}

	if err != nil {
		return nil, err
	}

	return &Plan{
		prog:         prog,
		host:         e.host,
		hooks:        e.hooks.plan,
		sessionHooks: e.hooks.session,
	}, nil
}

func (e *Engine) Run(ctx context.Context, src *file.Source, opts ...SessionOption) (Result, error) {
	plan, err := e.Compile(ctx, src)

	if err != nil {
		return nil, err
	}

	session, err := plan.NewSession(opts...)

	if err != nil {
		return nil, err
	}

	defer func() {
		// TODO: Log errors
		session.Close()
		plan.Close()
	}()

	return session.Run(ctx)
}

func (e *Engine) Close() error {
	if err := e.hooks.engine.runCloseHooks(); err != nil {
		return fmt.Errorf("close hooks: %w", err)
	}

	return nil
}
