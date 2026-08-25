package compiler

import (
	"fmt"

	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
)

const (
	// O0 represents no optimization level, where the compiler performs minimal or no optimizations.
	O0 = optimization.LevelNone
	// O1 represents basic optimization level, providing a balance between performance and resource usage.
	O1 = optimization.LevelBasic
)

type (
	// OptimizationLevel controls the compiler optimization pipeline.
	OptimizationLevel = optimization.Level

	config struct {
		Level     optimization.Level
		DebugInfo bool
	}

	// Option configures a Compiler during construction.
	Option = options.Option[config]
)

func defaultConfig() config {
	return config{
		Level: optimization.LevelBasic,
	}
}

// WithDebugInfo emits source-level debugger metadata and disables optimization
// so debugger-visible register bindings remain stable.
func WithDebugInfo() Option {
	return func(config *config) error {
		config.DebugInfo = true
		config.Level = optimization.LevelNone

		return nil
	}
}

// WithOptimizationLevel configures the compiler optimization level. It accepts
// levels from 0 through 3.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return options.New(
		func(config *config, level optimization.Level) {
			config.Level = level
		},
	).
		Named("optimization level").
		Validators(
			options.OneOf(
				optimization.LevelNone,
				optimization.LevelBasic,
				optimization.LevelFull,
				optimization.LevelAggressive,
			),
		).
		Value(level).
		Build()
}

func validateOptimizationLevel(level OptimizationLevel) error {
	switch level {
	case optimization.LevelNone,
		optimization.LevelBasic,
		optimization.LevelFull,
		optimization.LevelAggressive:
		return nil
	default:
		return fmt.Errorf("invalid optimization level: %d", level)
	}
}
