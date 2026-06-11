package compiler

import "github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"

const (
	// O0 represents no optimization level, where the compiler performs minimal or no optimizations.
	O0 = optimization.LevelNone
	// O1 represents basic optimization level, providing a balance between performance and resource usage.
	O1 = optimization.LevelBasic
)

type (
	OptimizationLevel = optimization.Level

	Option func(opts *options)

	options struct {
		Level     optimization.Level
		DebugInfo bool
	}
)

func WithOptimizationLevel(level optimization.Level) Option {
	return func(opts *options) {
		opts.Level = level
	}
}

// WithDebugInfo emits source-level debugger metadata and disables optimization
// so debugger-visible register bindings remain stable.
func WithDebugInfo() Option {
	return func(opts *options) {
		opts.DebugInfo = true
		opts.Level = optimization.LevelNone
	}
}
