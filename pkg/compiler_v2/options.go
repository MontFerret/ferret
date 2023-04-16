package compiler_v2

type (
	Option func(opts *Options)

	Options struct {
		noStdlib bool
	}
)

func WithoutStdlib() Option {
	return func(opts *Options) {
		opts.noStdlib = true
	}
}
