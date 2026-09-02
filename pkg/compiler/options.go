package compiler

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
)

const (
	// None disables optimizer passes.
	None = optimization.None
	// Basic enables constant propagation, liveness analysis, and peephole optimization.
	Basic = optimization.Basic
	// Full is the default and adds register coalescing to the Basic pipeline.
	Full = optimization.Full
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
		Level: Full,
	}
}

// WithOptimizationLevel configures the compiler optimizer pipeline. None disables
// optimizer passes, Basic enables the reduced pipeline, and Full enables the
// complete supported pipeline.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return options.New(
		func(config *config, level OptimizationLevel) {
			config.Level = level
		},
	).
		Named("optimization level").
		Validators(
			options.OneOf(
				None,
				Basic,
				Full,
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
		config.Level = None

		return nil
	}
}
