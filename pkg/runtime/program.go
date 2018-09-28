package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver"
)

type Program struct {
	exp core.Expression
}

func NewProgram(exp core.Expression) *Program {
	return &Program{exp}
}

func (p *Program) Run(ctx context.Context, setters ...Option) ([]byte, error) {
	scope, closeFn := core.NewRootScope()

	defer closeFn()

	opts := newOptions()

	for _, setter := range setters {
		setter(opts)
	}

	ctx = opts.withContext(ctx)
	ctx = driver.WithDynamicDriver(ctx, opts.cdp)
	ctx = driver.WithStaticDriver(ctx)

	out, err := p.exp.Exec(ctx, scope)

	if err != nil {
		js, _ := values.None.MarshalJSON()

		return js, err
	}

	return out.MarshalJSON()
}

func (p *Program) MustRun(ctx context.Context, setters ...Option) []byte {
	out, err := p.Run(ctx, setters...)

	if err != nil {
		panic(err)
	}

	return out
}
