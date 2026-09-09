package engine

import (
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type sessionConfig struct {
	logger            []logging.Option
	outputContentType string
	fsRoot            string
	env               []vm.EnvironmentOption
	debugFormat       debugger.FormatOptions
}

func defaultSessionConfig() sessionConfig {
	return sessionConfig{outputContentType: encodingjson.ContentType, debugFormat: debugger.DefaultFormatOptions()}
}

func newSessionConfig(setters []SessionOption) (sessionConfig, error) {
	if len(setters) == 0 {
		return defaultSessionConfig(), nil
	}

	opts := defaultSessionConfig()
	var failures []error
	for _, setter := range setters {
		if setter != nil {
			if err := setter(&opts); err != nil {
				failures = append(failures, err)
			}
		}
	}

	if err := errors.Join(failures...); err != nil {
		return sessionConfig{}, err
	}

	return opts, nil
}

func (cfg *sessionConfig) setParam(name string, value any) error {
	if name == "" {
		return fmt.Errorf("param name cannot be empty")
	}

	if value == nil {
		return fmt.Errorf("param value cannot be nil")
	}

	return cfg.setParams(map[string]any{name: value})
}

func (cfg *sessionConfig) setParams(params map[string]any) error {
	if len(params) == 0 {
		return nil
	}

	converted, err := runtime.NewParamsFrom(params)
	if err != nil {
		return fmt.Errorf("convert session params: %w", err)
	}

	cfg.env = append(cfg.env, vm.WithParams(converted))

	return nil
}
