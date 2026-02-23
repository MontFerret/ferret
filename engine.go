package ferret

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Engine struct {
	compiler  *compiler.Compiler
	functions runtime.Functions
	params    map[string]runtime.Value
	logging   runtime.LogSettings
	registry  *encoding.Registry
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
		registry:  opts.registry,
	}, nil
}

func (e *Engine) Codecs() *encoding.Registry {
	return e.registry
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
	}, e.registry), nil
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
