package ferret

import (
	"context"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Instance struct {
	compiler *compiler.Compiler
	drivers  *drivers.Container
}

func New(setters ...Option) *Instance {
	opts := NewOptions(setters)

	return &Instance{
		compiler: compiler.New(opts.compiler...),
		drivers:  drivers.NewContainer(),
	}
}

func (i *Instance) Functions() core.Namespace {
	return i.compiler
}

func (i *Instance) Drivers() *drivers.Container {
	return i.drivers
}

func (i *Instance) Compile(query string) (*runtime.Program, error) {
	return i.compiler.Compile(query)
}

func (i *Instance) MustCompile(query string) *runtime.Program {
	return i.compiler.MustCompile(query)
}

func (i *Instance) Exec(ctx context.Context, query string, opts ...runtime.Option) ([]byte, error) {
	p, err := i.Compile(query)

	if err != nil {
		return nil, err
	}

	ctx = i.drivers.WithContext(ctx)

	return p.Run(ctx, opts...)
}

func (i *Instance) MustExec(ctx context.Context, query string, opts ...runtime.Option) []byte {
	out, err := i.Exec(ctx, query, opts...)

	if err != nil {
		panic(err)
	}

	return out
}

func (i *Instance) Run(ctx context.Context, program *runtime.Program, opts ...runtime.Option) ([]byte, error) {
	if program == nil {
		return nil, core.Error(core.ErrInvalidArgument, "program")
	}

	ctx = i.drivers.WithContext(ctx)

	return program.Run(ctx, opts...)
}

func (i *Instance) MustRun(ctx context.Context, program *runtime.Program, opts ...runtime.Option) []byte {
	out, err := i.Run(ctx, program, opts...)

	if err != nil {
		panic(err)
	}

	return out
}
