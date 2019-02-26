package runtime

import (
	"context"
	"runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type Program struct {
	src  string
	body core.Expression
}

func NewProgram(src string, body core.Expression) (*Program, error) {
	if src == "" {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if body == nil {
		return nil, core.Error(core.ErrMissedArgument, "body")
	}

	return &Program{src, body}, nil
}

func (p *Program) Source() string {
	return p.src
}

func (p *Program) Run(ctx context.Context, setters ...Option) (result []byte, err error) {
	ctx = NewOptions(setters).WithContext(ctx)

	logger := logging.FromContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			b := make([]byte, 0, 20)
			runtime.Stack(b, true)

			logger.Error().
				Timestamp().
				Err(err).
				Str("stack", string(b)).
				Msg("Panic")

			result = nil
		}
	}()

	scope, closeFn := core.NewRootScope()

	defer func() {
		if err := closeFn(); err != nil {
			logger.Error().
				Timestamp().
				Err(err).
				Msg("closing root scope")
		}
	}()

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
