package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/vm/test"
)

type (
	PanicPolicy uint8

	options struct {
		shapeCacheLimit         int
		fastObjectDictThreshold int
		panicPolicy             PanicPolicy
	}

	configurator struct {
		testing test.Testing[*Result]
		options
	}

	Option func(*configurator)
)

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

const (
	PanicRecover PanicPolicy = iota
	PanicPropagate
)

func newOptions(opts []Option) (options, test.Testing[*Result]) {
	cfg := configurator{
		options: options{
			shapeCacheLimit:         defaultShapeCacheLimit,
			fastObjectDictThreshold: defaultFastObjectDictThreshold,
			panicPolicy:             PanicRecover,
		},
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg.options, cfg.testing
}

// WithShapeCacheLimit overrides the per-VM shape transition cache size.
func WithShapeCacheLimit(limit int) Option {
	return func(cfg *configurator) {
		if limit > 0 {
			cfg.shapeCacheLimit = limit
		}
	}
}

// WithFastObjectDictThreshold overrides the key count after which FastObject switches to dict mode.
func WithFastObjectDictThreshold(limit int) Option {
	return func(cfg *configurator) {
		if limit > 0 {
			cfg.fastObjectDictThreshold = limit
		}
	}
}

// WithPanicPolicy controls panic handling policy during Run.
// PanicRecover wraps panics into runtime errors, while PanicPropagate lets panics propagate.
func WithPanicPolicy(mode PanicPolicy) Option {
	return func(cfg *configurator) {
		switch mode {
		case PanicRecover, PanicPropagate:
			cfg.panicPolicy = mode
		}
	}
}

// WithTesting configures a testing instance for the VM, which is used to support test/benchmark-only features like the benchmark result mode.
// This is not intended for public use and may be removed in the future as test/benchmark features are integrated into the public API.
func WithTesting(opts ...test.Option) Option {
	return func(cfg *configurator) {
		cfg.testing = test.NewTesting[*Result](opts)
	}
}
