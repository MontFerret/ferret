package net

import (
	stdhttp "net/http"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
	"github.com/ziflex/go-options"
)

type (
	Option = options.Option[config]

	config struct {
		httpClient    ferrethttp.Client
		httpTransport stdhttp.RoundTripper
		httpPolicies  []ferrethttp.PolicyOption
	}
)

// WithHTTPClient sets the HTTP client used by a Network.
func WithHTTPClient(client ferrethttp.Client) Option {
	return func(opts *config, report options.Report) {
		if client == nil {
			report(options.ValidationError{
				Field:  "client",
				Reason: "nil client is not allowed",
			})

			return
		}

		opts.httpClient = client
	}
}

// WithHTTPPolicies sets the HTTP policies used by a Network.
func WithHTTPPolicies(policies ...ferrethttp.PolicyOption) Option {
	return func(opts *config, _ options.Report) {
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
	return func(opts *config, report options.Report) {
		if transport == nil {
			report(options.ValidationError{
				Field:  "transport",
				Reason: "nil transport is not allowed",
			})

			return
		}

		opts.httpTransport = transport

		if len(policies) > 0 {
			opts.httpPolicies = append(opts.httpPolicies, policies...)
		}
	}
}
