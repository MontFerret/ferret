package static

import "github.com/sethgrid/pester"

type (
	Option func(opts *pester.Client)
)

func WithDefaultBackoff() Option {
	return func(opts *pester.Client) {
		opts.Backoff = pester.DefaultBackoff
	}
}

func WithLinearBackoff() Option {
	return func(opts *pester.Client) {
		opts.Backoff = pester.LinearBackoff
	}
}

func WithExponentialBackoff() Option {
	return func(opts *pester.Client) {
		opts.Backoff = pester.ExponentialBackoff
	}
}

func WithMaxRetries(value int) Option {
	return func(opts *pester.Client) {
		opts.MaxRetries = value
	}
}

func WithConcurrency(value int) Option {
	return func(opts *pester.Client) {
		opts.Concurrency = value
	}
}
