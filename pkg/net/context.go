package net

import (
	"context"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type contextKey struct{}

var ctxKey = contextKey{}

// WithNetwork adds a Network instance to the provided context, creating a new context if the input context is nil.
func WithNetwork(ctx context.Context, net Network) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ctxKey, net)
}

// NetworkFrom retrieves a Network instance from the provided context, returning an error if not found.
func NetworkFrom(ctx context.Context) (Network, error) {
	if ctx == nil {
		return nil, ErrNotFound
	}

	val := ctx.Value(ctxKey)
	if val == nil {
		return nil, ErrNotFound
	}

	net, ok := val.(Network)

	if !ok {
		return nil, ErrNotFound
	}

	return net, nil
}

// HTTPClientFrom retrieves the HTTP client from a Network instance stored in context.
func HTTPClientFrom(ctx context.Context) (ferrethttp.Client, error) {
	net, err := NetworkFrom(ctx)
	if err != nil {
		return nil, err
	}

	client := net.HTTP()
	if client == nil {
		return nil, ErrHTTPClientNotFound
	}

	return client, nil
}
