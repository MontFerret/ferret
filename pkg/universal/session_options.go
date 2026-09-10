package universal

import (
	"context"
	"errors"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type sessionOptions struct {
	native []engine.SessionOption
}

var _ api.SessionOptions = (*sessionOptions)(nil)

func (o *sessionOptions) SetParam(name string, value any) error {
	o.native = append(o.native, engine.WithSessionParam(name, value))

	return nil
}

func (o *sessionOptions) SetParams(params map[string]any) error {
	o.native = append(o.native, engine.WithSessionParams(params))

	return nil
}

func (o *sessionOptions) SetOutputContentType(value string) error {
	o.native = append(o.native, engine.WithOutputContentType(value))

	return nil
}

func (o *sessionOptions) SetFSRoot(value string) error {
	o.native = append(o.native, engine.WithSessionFSRoot(value))

	return nil
}

func newSessionOptions(ctx context.Context, setters []api.SessionOption) (*sessionOptions, error) {
	if err := checkOptionContext(ctx); err != nil {
		return nil, err
	}

	opts := &sessionOptions{}
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	if err := ctx.Err(); err != nil {
		failures = append(failures, err)
	}

	if err := errors.Join(failures...); err != nil {
		return nil, err
	}

	return opts, nil
}
