package compiler

import "github.com/MontFerret/ferret/pkg/compiler/internal/optimization"

const (
	// O0 represents no optimization level, where the compiler performs minimal or no optimizations.
	O0 = optimization.LevelNone
	// O1 represents basic optimization level, providing a balance between performance and resource usage.
	O1 = optimization.LevelBasic
	// O2 represents full optimization level, aiming for optimal performance with moderate resource usage.
	O2 = optimization.LevelFull
	// O3 represents aggressive optimization level, prioritizing maximum performance, potentially at the cost of higher resource consumption.
	O3 = optimization.LevelAggressive
)

type (
	OptimizationLevel = optimization.Level

	Option func(opts *options)

	options struct {
		Level optimization.Level
	}
)

func WithOptimizationLevel(level optimization.Level) Option {
	return func(opts *options) {
		opts.Level = level
	}
}
