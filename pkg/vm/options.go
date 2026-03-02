package vm

const (
	defaultShapeCacheLimit         = 128
	defaultFastObjectDictThreshold = 64
)

type VMOption func(*vmConfig)

type vmConfig struct {
	shapeCacheLimit         int
	fastObjectDictThreshold int
}

func defaultVMConfig() vmConfig {
	return vmConfig{
		shapeCacheLimit:         defaultShapeCacheLimit,
		fastObjectDictThreshold: defaultFastObjectDictThreshold,
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
