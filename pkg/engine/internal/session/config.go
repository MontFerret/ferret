package session

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Config holds applied session settings before acquisition. Native option
// validation and ordering remain with the engine's private option target.
type Config struct {
	OutputContentType string
	FSRoot            string
	Logger            []logging.Option
	Environment       []vm.EnvironmentOption
	DebugFormat       debugger.FormatOptions
}

// NewConfig supplies defaults shared by ordinary and debug sessions.
func NewConfig() Config {
	return Config{OutputContentType: encodingjson.ContentType, DebugFormat: debugger.DefaultFormatOptions()}
}

// SetParam converts a named host value into an environment override.
func (cfg *Config) SetParam(name string, value any) error {
	if name == "" {
		return fmt.Errorf("param name cannot be empty")
	}

	if value == nil {
		return fmt.Errorf("param value cannot be nil")
	}

	return cfg.SetParams(map[string]any{name: value})
}

// SetParams appends converted overrides without changing earlier options.
func (cfg *Config) SetParams(params map[string]any) error {
	if len(params) == 0 {
		return nil
	}

	converted, err := runtime.NewParamsFrom(params)
	if err != nil {
		return fmt.Errorf("convert session params: %w", err)
	}

	cfg.Environment = append(cfg.Environment, vm.WithParams(converted))

	return nil
}
