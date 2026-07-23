package net

import (
	stdhttp "net/http"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type (
	Option func(*options)

	options struct {
		httpClient    ferrethttp.Client
		httpTransport stdhttp.RoundTripper
		httpPolicies  []ferrethttp.PolicyOption
	}
)

// WithHTTPClient sets the HTTP client used by a Network.
func WithHTTPClient(client ferrethttp.Client) Option {
	return func(opts *options) {
		if client == nil {
			return
		}

		opts.httpClient = client
	}
}

// WithHTTPPolicies sets the HTTP policies used by a Network.
func WithHTTPPolicies(policies ...ferrethttp.PolicyOption) Option {
	return func(opts *options) {
		if len(policies) == 0 {
			return
		}

		opts.httpPolicies = append(opts.httpPolicies, policies...)
	}
}

// WithHTTPTransport sets the standard-library transport used by the
// policy-aware HTTP client. The transport is ignored when WithHTTPClient
// supplies a client. A nil transport makes this option a no-op.
//
// Custom transports remain responsible for proxy behavior, DNS and
// concrete-address enforcement, and response-header limits.
func WithHTTPTransport(transport stdhttp.RoundTripper, policies ...ferrethttp.PolicyOption) Option {
	return func(opts *options) {
		if transport == nil {
			return
		}

		opts.httpTransport = transport

		if len(policies) > 0 {
			opts.httpPolicies = append(opts.httpPolicies, policies...)
		}
	}
}
