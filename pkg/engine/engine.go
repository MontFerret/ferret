package engine

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type Engine struct {
	compiler  *compiler.Compiler
	functions runtime.Functions
	params    map[string]runtime.Value
	logging   runtime.LogSettings
}

func New(setters ...Option) (*Engine, error) {
	opts, err := newOptions(setters)

	if err != nil {
		return nil, err
	}

	return &Engine{
		compiler:  compiler.New(opts.compiler...),
		functions: opts.functions.Build(),
		params:    opts.params,
		logging:   opts.logging,
	}, nil
}

func (e *Engine) Compile(src *file.Source) (*Plan, error) {
	prog, err := e.compiler.Compile(src)

	if err != nil {
		return nil, err
	}

	return newPlan(prog, &vm.Environment{
		Functions: e.functions,
		Params:    e.params,
		Logging:   e.logging,
	}), nil
}

func (e *Engine) MustCompile(src *file.Source) *Plan {
	program, err := e.Compile(src)

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
