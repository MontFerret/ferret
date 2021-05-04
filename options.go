package ferret

import "github.com/MontFerret/ferret/pkg/compiler"

type (
	Options struct {
		compiler []compiler.Option
	}

	Option func(opts *Options)
)

func NewOptions(setters []Option) *Options {
	res := &Options{
		compiler: make([]compiler.Option, 0, 2),
	}

	for _, setter := range setters {
		setter(res)
	}

	return res
}

func WithoutStdlib() Option {
	return func(opts *Options) {
		opts.compiler = append(opts.compiler, compiler.WithoutStdlib())
	}
}
