package net

import ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"

type (
	Option func(*options)

	options struct {
		httpClient   ferrethttp.Client
		httpPolicies []ferrethttp.PolicyOption
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
