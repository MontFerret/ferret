package engine

import (
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	options struct {
		compiler []compiler.Option
		env      *vm.Environment
	}

	Option func(env *options)
)

func newOptions(setters []Option) *options {
	opts := &options{}

	for _, opt := range setters {
		opt(opts)
	}

	return opts
}

func WithEnvironment(setter *vm.Environment) Option {
	return func(o *options) {
		o.env = setter
	}
}
