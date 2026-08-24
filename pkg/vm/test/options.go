package test

import "github.com/ziflex/go-options"

type (
	// Options contains testing-only VM construction settings.
	Options struct {
		BenchmarkMode bool
	}

	// Option configures testing-only VM state.
	Option = options.Option[Options]
)

func defaultOptions() Options {
	return Options{}
}

// NewOptions resolves VM testing options.
func NewOptions(opts []Option) (Options, error) {
	cfg, err := options.ApplyTo(defaultOptions(), opts...)
	if err != nil {
		return Options{}, err
	}

	return cfg, nil
}

// WithBenchmarkMode enables reusable benchmark results.
func WithBenchmarkMode() Option {
	return func(opts *Options) error {
		opts.BenchmarkMode = true

		return nil
	}
}
