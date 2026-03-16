package test

type (
	Options struct {
		BenchmarkMode bool
	}

	Option func(*Options)
)

func NewOptions(opts []Option) Options {
	cfg := &Options{}

	for _, opt := range opts {
		opt(cfg)
	}

	return *cfg
}

func WithBenchmarkMode() Option {
	return func(opts *Options) {
		opts.BenchmarkMode = true
	}
}
