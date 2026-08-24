package net

import (
	stdhttp "net/http"

	"github.com/ziflex/go-options"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type (
	config struct {
		httpClient    ferrethttp.Client
		httpTransport stdhttp.RoundTripper
		httpPolicies  []ferrethttp.PolicyOption
	}

	Option = options.Option[config]
)

func defaultConfig() config {
	return config{}
}

// WithHTTPClient sets the HTTP client used by a Network.
func WithHTTPClient(client ferrethttp.Client) Option {
	return func(config *config) error {
		if client == nil {
			return nil
		}

		config.httpClient = client

		return nil
	}
}

// WithHTTPPolicies sets the HTTP policies used by a Network.
func WithHTTPPolicies(policies ...ferrethttp.PolicyOption) Option {
	return func(config *config) error {
		if len(policies) == 0 {
			return nil
		}

		config.httpPolicies = append(config.httpPolicies, policies...)

		return nil
	}
}

// WithHTTPTransport sets the standard-library transport used by the
// policy-aware HTTP client. The transport is ignored when WithHTTPClient
// supplies a client. A nil transport makes this option a no-op.
//
// Custom transports remain responsible for proxy behavior, DNS and
// concrete-address enforcement, and response-header limits.
func WithHTTPTransport(transport stdhttp.RoundTripper, policies ...ferrethttp.PolicyOption) Option {
	return func(config *config) error {
		if transport == nil {
			return nil
		}

		config.httpTransport = transport

		if len(policies) > 0 {
			config.httpPolicies = append(config.httpPolicies, policies...)
		}

		return nil
	}
}
