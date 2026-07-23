package net

import (
	"fmt"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type (
	// Network defines an interface to provide access to HTTP operations via an HTTP client.
	Network interface {
		HTTP() ferrethttp.Client
	}

	defaultNetwork struct {
		http ferrethttp.Client
	}
)

// New constructs a Network with a default HTTP client unless one is supplied.
// It returns an error when the default client cannot be initialized.
func New(setters ...Option) (Network, error) {
	opts := options{}

	for _, option := range setters {
		if option == nil {
			continue
		}

		option(&opts)
	}

	if opts.httpClient == nil {
		client, err := ferrethttp.New(opts.httpPolicies...)

		if err != nil {
			return nil, fmt.Errorf("http client: %w", err)
		}

		opts.httpClient = client
	}

	return &defaultNetwork{
		http: opts.httpClient,
	}, nil
}

func (n *defaultNetwork) HTTP() ferrethttp.Client {
	return n.http
}

func (n *defaultNetwork) CloseIdleConnections() {
	if closer, ok := n.http.(ferrethttp.IdleConnectionCloser); ok && closer != nil {
		closer.CloseIdleConnections()
	}
}
