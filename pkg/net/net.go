package net

import ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"

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
func New(setters ...Option) Network {
	opts := options{}

	for _, option := range setters {
		if option == nil {
			continue
		}

		option(&opts)
	}

	if opts.http == nil {
		opts.http = ferrethttp.New()
	}

	return &defaultNetwork{
		http: opts.http,
	}
}

func (n *defaultNetwork) HTTP() ferrethttp.Client {
	return n.http
}

func (n *defaultNetwork) CloseIdleConnections() {
	if closer, ok := n.http.(ferrethttp.IdleConnectionCloser); ok {
		closer.CloseIdleConnections()
	}
}
