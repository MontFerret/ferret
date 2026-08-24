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

// WithShapeCacheLimit overrides the per-VM shape transition cache size. The
// limit must be positive; invalid values make VM construction fail.
func WithShapeCacheLimit(limit int) Option {
	return options.New(func(cfg *config, limit int) {
		cfg.shapeCacheLimit = limit
	}).
		Value(limit).
		Named("shape cache limit").
		Validators(options.Positive[int]()).
		Build()
}

// WithFastObjectDictThreshold overrides the key count after which FastObject
// switches to dict mode. The threshold must be positive; invalid values make VM
// construction fail.
func WithFastObjectDictThreshold(limit int) Option {
	return options.New(func(cfg *config, limit int) {
		cfg.fastObjectDictThreshold = limit
	}).
		Value(limit).
		Named("fast object dict threshold").
		Validators(options.Positive[int]()).
		Build()
}

// WithPanicPolicy controls panic handling policy during Run.
// PanicRecover wraps panics into runtime errors, while PanicPropagate lets
// panics propagate. Unsupported values make VM construction fail.
func WithPanicPolicy(mode PanicPolicy) Option {
	return options.New(func(cfg *config, mode PanicPolicy) {
		cfg.panicPolicy = mode
	}).
		Value(mode).
		Named("panic policy").
		Validators(options.OneOf(PanicRecover, PanicPropagate)).
		Build()
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
