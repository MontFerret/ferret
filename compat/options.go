package compat

import (
	ferret "github.com/MontFerret/ferret/v2"
)

// Option is a functional option for configuring a compat Instance.
type Option func(*instanceOptions) error

// instanceOptions holds configuration for Instance construction.
type instanceOptions struct {
	engineOpts []ferret.Option
	noStdlib   bool
}

func newInstanceOptions(setters []Option) (*instanceOptions, error) {
	o := &instanceOptions{}

	for _, s := range setters {
		if s == nil {
			continue
		}

		if err := s(o); err != nil {
			return nil, err
		}
	}

	return o, nil
}

// WithoutStdlib disables the Ferret standard library.
//
// When this option is used, no built-in functions (MATH_*, STRING_*, etc.) are
// registered. Users must supply all required functions manually via
// Instance.Functions().
func WithoutStdlib() Option {
	return func(o *instanceOptions) error {
		o.noStdlib = true
		return nil
	}
}
