package net

import (
	"fmt"

	"github.com/ziflex/go-options"

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
// It returns option-application errors unchanged and an error when the default
// client cannot be initialized.
func New(setters ...Option) (Network, error) {
	cfg, err := options.ApplyTo(defaultConfig(), setters...)
	if err != nil {
		return nil, err
	}

	if cfg.httpClient == nil {
		var client ferrethttp.Client

		if cfg.httpTransport == nil {
			client, err = ferrethttp.New(cfg.httpPolicies...)
		} else {
			client, err = ferrethttp.NewWithTransport(
				cfg.httpTransport,
				cfg.httpPolicies...,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("http client: %w", err)
		}

		cfg.httpClient = client
	}

	return &defaultNetwork{
		http: cfg.httpClient,
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
