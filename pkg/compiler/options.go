package compiler

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
)

const (
	// OptimizationNone disables optimizer passes.
	OptimizationNone = optimization.None
	// OptimizationBasic enables constant propagation, liveness analysis, and peephole optimization.
	OptimizationBasic = optimization.Basic
	// OptimizationFull is the default and adds register coalescing to the Basic pipeline.
	OptimizationFull = optimization.Full
)

type (
	config struct {
		Level     optimization.Level
		DebugInfo bool
	}

	// OptimizationLevel controls the compiler optimization pipeline.
	OptimizationLevel = optimization.Level

	// Option configures a Compiler during construction.
	Option = options.Option[config]
)

func defaultConfig() config {
	return config{
		Level: OptimizationFull,
	}
}

// WithOptimizationLevel configures the compiler optimization level. It accepts
// levels from 0 through 2.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return options.New(
		func(config *config, level OptimizationLevel) {
			config.Level = level
		},
	).
		Named("optimization level").
		Validators(
			options.OneOf(
				OptimizationNone,
				OptimizationBasic,
				OptimizationFull,
			),
		).
		Value(level).
		Build()
}

// WithDebugInfo emits source-level debugger metadata and disables optimization
// so debugger-visible register bindings remain stable.
func WithDebugInfo() Option {
	return func(config *config) error {
		config.DebugInfo = true
		config.Level = OptimizationNone

		return nil
	}
}
