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

func New() *Instance {
	return &Instance{
		compiler: compiler.New(),
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

func (i *Instance) Exec(ctx context.Context, query string, opts ...runtime.Option) ([]byte, error) {
	p, err := i.Compile(query)

	if err != nil {
		return nil, err
	}

	for _, drv := range i.drivers.GetAll() {
		ctx = drivers.WithContext(ctx, drv)
	}

	return p.Run(ctx, opts...)
}
