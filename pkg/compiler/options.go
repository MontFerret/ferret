package compiler

import "github.com/MontFerret/ferret/pkg/compiler/internal/optimization"

type (
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
