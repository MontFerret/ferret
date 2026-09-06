package ferret

import (
	"errors"
	"fmt"
	"strings"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type sessionOptions struct {
	logger            []logging.Option
	outputContentType string
	fsRoot            string
	env               []vm.EnvironmentOption
	debugFormat       debugger.FormatOptions
}

var _ api.SessionOptions = (*sessionOptions)(nil)

func defaultSessionOptions() sessionOptions {
	return sessionOptions{outputContentType: encodingjson.ContentType, debugFormat: debugger.DefaultFormatOptions()}
}

func newSessionOptions(setters []SessionOption) (sessionOptions, error) {
	if len(setters) == 0 {
		return defaultSessionOptions(), nil
	}

	opts := defaultSessionOptions()
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(&opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	if err := errors.Join(failures...); err != nil {
		return sessionOptions{}, err
	}

	return opts, nil
}

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

	o.env = append(o.env, vm.WithParams(converted))

	return nil
}

func (o *sessionOptions) SetOutputContentType(value string) error {
	return gooptions.New(func(o *sessionOptions, value string) { o.outputContentType = strings.TrimSpace(value) }).
		Value(value).Named("output content type").Validators(gooptions.NotBlank[string]()).Build()(o)
}

func (o *sessionOptions) SetFSRoot(value string) error {
	return gooptions.New(func(o *sessionOptions, value string) { o.fsRoot = strings.TrimSpace(value) }).
		Value(value).Named("fs root").Validators(gooptions.NotBlank[string]()).Build()(o)
}
