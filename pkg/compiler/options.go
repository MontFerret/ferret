package compiler

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/ziflex/go-options"
)

const (
	// O0 represents no optimization level, where the compiler performs minimal or no optimizations.
	O0 = optimization.LevelNone
	// O1 represents basic optimization level, providing a balance between performance and resource usage.
	O1 = optimization.LevelBasic
)

type (
	OptimizationLevel = optimization.Level

	Option = options.Option[config]

	config struct {
		Level     optimization.Level
		DebugInfo bool
	}
)

func WithOptimizationLevel(level optimization.Level) Option {
	return func(opts *config, report options.Report) {
		if level < optimization.LevelNone || level > optimization.LevelBasic {
			report(options.ValidationError{
				Field:  "Level",
				Value:  fmt.Sprintf("%d", level),
				Reason: "invalid optimization level",
			})

			return
		}

		opts.Level = level
	}
}

// WithDebugInfo emits source-level debugger metadata and disables optimization
// so debugger-visible register bindings remain stable.
func WithDebugInfo() Option {
	return func(opts *config, _ options.Report) {
		opts.DebugInfo = true
		opts.Level = optimization.LevelNone
	}
}
