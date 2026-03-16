package vm

type (
	PanicPolicy uint8

	options struct {
		shapeCacheLimit         int
		fastObjectDictThreshold int
		benchmarkResultMode     bool
		panicPolicy             PanicPolicy
	}

	Option func(*options)
)

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

const (
	PanicRecover PanicPolicy = iota
	PanicPropagate
)

func newOptions(opts []Option) options {
	cfg := options{
		shapeCacheLimit:         defaultShapeCacheLimit,
		fastObjectDictThreshold: defaultFastObjectDictThreshold,
		panicPolicy:             PanicRecover,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// WithShapeCacheLimit overrides the per-VM shape transition cache size.
func WithShapeCacheLimit(limit int) Option {
	return func(cfg *options) {
		if limit > 0 {
			cfg.shapeCacheLimit = limit
		}
	}
}

// WithFastObjectDictThreshold overrides the key count after which FastObject switches to dict mode.
func WithFastObjectDictThreshold(limit int) Option {
	return func(cfg *options) {
		if limit > 0 {
			cfg.fastObjectDictThreshold = limit
		}
	}
}

// WithPanicPolicy controls panic handling policy during Run.
// PanicRecover wraps panics into runtime errors, while PanicPropagate lets panics propagate.
func WithPanicPolicy(mode PanicPolicy) Option {
	return func(cfg *options) {
		switch mode {
		case PanicRecover, PanicPropagate:
			cfg.panicPolicy = mode
		}
	}
}

// WithBenchmarkResultMode enables a benchmark/test-only result mode that
// reuses a single Result handle across runs. This keeps benchmark hot paths
// comparable to pre-result-handle runs while preserving the public Run API.
func WithBenchmarkResultMode() Option {
	return func(cfg *options) {
		cfg.benchmarkResultMode = true
	}
}
