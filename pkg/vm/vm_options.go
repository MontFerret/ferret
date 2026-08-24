package vm

import (
	"github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/vm/test"
)

type (
	PanicPolicy uint8

	config struct {
		testing                 test.Testing[*Result]
		shapeCacheLimit         int
		fastObjectDictThreshold int
		panicPolicy             PanicPolicy
	}

	// Option configures a VM during construction.
	Option = options.Option[config]
)

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

const (
	PanicRecover PanicPolicy = iota
	PanicPropagate
)

func defaultConfig() config {
	return config{
		shapeCacheLimit:         defaultShapeCacheLimit,
		fastObjectDictThreshold: defaultFastObjectDictThreshold,
		panicPolicy:             PanicRecover,
	}
}

// WithShapeCacheLimit overrides the per-VM shape transition cache size.
func WithShapeCacheLimit(limit int) Option {
	return func(cfg *config) error {
		if limit > 0 {
			cfg.shapeCacheLimit = limit
		}

		return nil
	}
}

// WithFastObjectDictThreshold overrides the key count after which FastObject switches to dict mode.
func WithFastObjectDictThreshold(limit int) Option {
	return func(cfg *config) error {
		if limit > 0 {
			cfg.fastObjectDictThreshold = limit
		}

		return nil
	}
}

// WithPanicPolicy controls panic handling policy during Run.
// PanicRecover wraps panics into runtime errors, while PanicPropagate lets panics propagate.
func WithPanicPolicy(mode PanicPolicy) Option {
	return func(cfg *config) error {
		switch mode {
		case PanicRecover, PanicPropagate:
			cfg.panicPolicy = mode
		}

		return nil
	}
}

// WithTesting configures a testing instance for the VM, which is used to support test/benchmark-only features like the benchmark result mode.
// This is not intended for public use and may be removed in the future as test/benchmark features are integrated into the public API.
func WithTesting(opts ...test.Option) Option {
	return func(cfg *config) error {
		testing, err := test.NewTesting[*Result](opts)
		if err != nil {
			return err
		}

		cfg.testing = testing

		return nil
	}
}
