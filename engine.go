package ferret

import (
	"context"

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

	return &Engine{
		compiler: compiler.New(opts.compiler...),
		host:     h,
		hooks:    boot.hooks.clone(),
	}, nil
}

func (e *Engine) Compile(src *file.Source) (*Plan, error) {
	prog, err := e.compiler.Compile(src)

	if err != nil {
		return nil, err
	}

	return &Plan{
		prog:  prog,
		host:  e.host,
		hooks: e.hooks.session,
	}, nil
}

func (e *Engine) Run(ctx context.Context, src *file.Source, opts ...SessionOption) (Result, error) {
	plan, err := e.Compile(src)

	if err != nil {
		return nil, err
	}

	session, err := plan.NewSession(opts...)

	if err != nil {
		return nil, err
	}

	defer session.Close()

	return session.Run(ctx)
}

func (e *Engine) Close() error {
	return e.hooks.engine.runCloseHooks()
}
