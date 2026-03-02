package vm

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

type (
	RunSafetyMode uint8

	vmConfig struct {
		shapeCacheLimit         int
		fastObjectDictThreshold int
		runSafetyMode           RunSafetyMode
	}

	VMOption func(*vmConfig)
)

const (
	RunSafetyStrict RunSafetyMode = iota
	RunSafetyFast
)

func defaultVMConfig() vmConfig {
	return vmConfig{
		shapeCacheLimit:         defaultShapeCacheLimit,
		fastObjectDictThreshold: defaultFastObjectDictThreshold,
		runSafetyMode:           RunSafetyStrict,
	}
}

// WithShapeCacheLimit overrides the per-VM shape transition cache size.
func WithShapeCacheLimit(limit int) VMOption {
	return func(cfg *vmConfig) {
		if limit > 0 {
			cfg.shapeCacheLimit = limit
		}
	}
}

// WithFastObjectDictThreshold overrides the key count after which FastObject switches to dict mode.
func WithFastObjectDictThreshold(limit int) VMOption {
	return func(cfg *vmConfig) {
		if limit > 0 {
			cfg.fastObjectDictThreshold = limit
		}
	}
}

// WithRunSafetyMode controls panic handling policy during Run.
// RunSafetyStrict wraps panics into runtime errors, while RunSafetyFast lets panics propagate.
func WithRunSafetyMode(mode RunSafetyMode) VMOption {
	return func(cfg *vmConfig) {
		switch mode {
		case RunSafetyStrict, RunSafetyFast:
			cfg.runSafetyMode = mode
		}
	}
}
