package vm

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

type (
	RunSafetyMode uint8

	options struct {
		shapeCacheLimit         int
		fastObjectDictThreshold int
		runSafetyMode           RunSafetyMode
	}

	Option func(*options)
)

const (
	RunSafetyStrict RunSafetyMode = iota
	RunSafetyFast
)

func newOptions(opts []Option) options {
	cfg := options{
		shapeCacheLimit:         defaultShapeCacheLimit,
		fastObjectDictThreshold: defaultFastObjectDictThreshold,
		runSafetyMode:           RunSafetyStrict,
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

// WithRunSafetyMode controls panic handling policy during Run.
// RunSafetyStrict wraps panics into runtime errors, while RunSafetyFast lets panics propagate.
func WithRunSafetyMode(mode RunSafetyMode) Option {
	return func(cfg *options) {
		switch mode {
		case RunSafetyStrict, RunSafetyFast:
			cfg.runSafetyMode = mode
		}
	}
}
