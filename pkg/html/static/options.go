package static

import (
	"github.com/sethgrid/pester"
)

type (
	Option  func(opts *Options)
	Options struct {
		backoff     pester.BackoffStrategy
		maxRetries  int
		concurrency int
		proxy       string
		userAgent   string
	}
)

func WithDefaultBackoff() Option {
	return func(opts *Options) {
		opts.backoff = pester.DefaultBackoff
	}
}

func WithLinearBackoff() Option {
	return func(opts *Options) {
		opts.backoff = pester.LinearBackoff
	}
}

func WithExponentialBackoff() Option {
	return func(opts *Options) {
		opts.backoff = pester.ExponentialBackoff
	}
}

func WithMaxRetries(value int) Option {
	return func(opts *Options) {
		opts.maxRetries = value
	}
}

func WithConcurrency(value int) Option {
	return func(opts *Options) {
		opts.concurrency = value
	}
}

func WithProxy(address string) Option {
	return func(opts *Options) {
		opts.proxy = address
	}
}

func WithUserAgent(value string) Option {
	return func(opts *Options) {
		opts.userAgent = value
	}
}
