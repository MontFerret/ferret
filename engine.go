package ferret

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
)

type Engine struct {
	compiler *compiler.Compiler
}

func NewEngine(setters ...Option) *Engine {
	opts := NewOptions(setters)

	return &Engine{
		compiler: compiler.New(opts.compiler...),
		//drivers:  drivers.NewContainer(),
	}
}

func (e *Engine) Compile(src *file.Source) (*Plan, error) {
	prog, err := e.compiler.Compile(src)

	if err != nil {
		return nil, err
	}

	return &Plan{
		prog: prog,
	}, nil
}

func (e *Engine) Run(ctx context.Context, src *file.Source, opts ...SessionOption) (Result, error) {
	plan, err := e.Compile(src)

	if err != nil {
		return nil, err
	}

	session := plan.NewSession()
	defer session.Close()

	return session.Run(ctx, opts)
}
