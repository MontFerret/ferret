package engine

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
)

type Engine struct {
	compiler *compiler.Compiler
	opts     *options
}

func New(setters ...Option) *Engine {
	opts := newOptions(setters)

	return &Engine{
		compiler: compiler.New(opts.compiler...),
		opts:     opts,
		//drivers:  drivers.NewContainer(),
	}
}

func (e *Engine) Compile(src *file.Source, opts ...PlanOption) (*Plan, error) {
	prog, err := e.compiler.Compile(src)

	if err != nil {
		return nil, err
	}

	if e.opts.env != nil {
		opts = append(opts, WithPlanEnvironment(e.opts.env))
	}

	return newPlan(prog, opts), nil
}

func (e *Engine) MustCompile(src *file.Source, opts ...PlanOption) *Plan {
	program, err := e.Compile(src, opts...)

	if err != nil {
		panic(err)
	}

	return program
}

func (e *Engine) Run(ctx context.Context, src *file.Source, opts ...SessionOption) (Result, error) {
	plan, err := e.Compile(src)

	if err != nil {
		return nil, err
	}

	session := plan.NewSession(opts...)
	defer session.Close()

	return session.Run(ctx)
}
