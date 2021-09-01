package runtime

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Program struct {
	src    string
	body   core.Expression
	params map[string]struct{}
}

func NewProgram(src string, body core.Expression, params map[string]struct{}) (*Program, error) {
	if src == "" {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if body == nil {
		return nil, core.Error(core.ErrMissedArgument, "body")
	}

	return &Program{src, body, params}, nil
}

func (p *Program) Source() string {
	return p.src
}

func (p *Program) Params() []string {
	res := make([]string, 0, len(p.params))

	for name := range p.params {
		res = append(res, name)
	}

	return res
}

func (p *Program) Run(ctx context.Context, setters ...Option) (result []byte, err error) {
	opts := NewOptions(setters)

	err = p.validateParams(opts)

	if err != nil {
		return nil, err
	}

	ctx = opts.WithContext(ctx)
	logger := logging.FromContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = errors.WithStack(err)
			default:
				err = errors.New("unknown panic")
			}

			logger.Error().
				Timestamp().
				Err(err).
				Str("stack", fmt.Sprintf("%+v", err)).
				Msg("panic")

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

func (p *Program) validateParams(opts *Options) error {
	if len(p.params) == 0 {
		return nil
	}

	// There might be no errors.
	// Thus, we allocate this slice lazily, on a first error.
	var missedParams []string

	for n := range p.params {
		_, exists := opts.params[n]

		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(p.params))
			}

			missedParams = append(missedParams, "@"+n)
		}
	}

	if len(missedParams) > 0 {
		return core.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return nil
}
