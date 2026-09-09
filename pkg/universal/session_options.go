package universal

import (
	"context"
	"errors"
	"fmt"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type sessionOptions struct {
	native []engine.SessionOption
}

var _ api.SessionOptions = (*sessionOptions)(nil)

func (o *sessionOptions) SetParam(name string, value any) error {
	if name == "" {
		return fmt.Errorf("param name cannot be empty")
	}

	if value == nil {
		return fmt.Errorf("param value cannot be nil")
	}

	return o.SetParams(map[string]any{name: value})
}

func (o *sessionOptions) SetParams(params map[string]any) error {
	if len(params) == 0 {
		return nil
	}

	converted, err := runtime.NewParamsFrom(params)
	if err != nil {
		return fmt.Errorf("convert session params: %w", err)
	}

	o.native = append(o.native, engine.WithSessionRuntimeParams(converted))

	return nil
}

func (o *sessionOptions) SetOutputContentType(value string) error {
	return gooptions.New(func(opts *sessionOptions, value string) {
		opts.native = append(opts.native, engine.WithOutputContentType(value))
	}).Value(value).Named("output content type").Validators(gooptions.NotBlank[string]()).Build()(o)
}

func (o *sessionOptions) SetFSRoot(value string) error {
	return gooptions.New(func(opts *sessionOptions, value string) {
		opts.native = append(opts.native, engine.WithSessionFSRoot(value))
	}).Value(value).Named("fs root").Validators(gooptions.NotBlank[string]()).Build()(o)
}

func newSessionOptions(ctx context.Context, setters []api.SessionOption) (*sessionOptions, error) {
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
