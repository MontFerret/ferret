package net

import ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"

type (
	Option func(*options)

	options struct {
		http ferrethttp.Client
	}
)

// WithHTTPClient sets the HTTP client used by a Network.
func WithHTTPClient(client ferrethttp.Client) Option {
	return func(opts *options) {
		if client == nil {
			return
		}

		opts.http = client
	}
}
