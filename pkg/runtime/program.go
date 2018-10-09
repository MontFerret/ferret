package runtime

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Program struct {
	src  string
	body core.Expression
}

func NewProgram(src string, body core.Expression) (*Program, error) {
	if src == "" {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if core.IsNil(body) {
		return nil, core.Error(core.ErrMissedArgument, "body")
	}

	return &Program{src, body}, nil
}

func (p *Program) Source() string {
	return p.src
}

func (p *Program) Run(ctx context.Context, setters ...Option) ([]byte, error) {
	scope, closeFn := core.NewRootScope()

	defer closeFn()

	opts := newOptions()

	for _, setter := range setters {
		setter(opts)
	}

	ctx = opts.withContext(ctx)
	ctx = html.WithDynamicDriver(ctx)
	ctx = html.WithStaticDriver(ctx)

	out, err := p.body.Exec(ctx, scope)

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
