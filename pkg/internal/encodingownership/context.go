// Package encodingownership carries synchronous resource adoption during native
// output materialization without making ownership part of the public codec API.
package encodingownership

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	valueAdopterKey struct{}

	operationContext struct {
		context.Context
		adopt func(runtime.Value)
	}
)

// WithValueAdopter lends adopt to codecs for this materialization only. Codecs
// must not retain it or call it after returning. It registers ownership without
// encoding, closing resources, or observing cancellation.
func WithValueAdopter(ctx context.Context, adopt func(runtime.Value)) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, valueAdopterKey{}, adopt)
}

// ValueAdopterFromContext returns the operation's adoption callback, if present.
func ValueAdopterFromContext(ctx context.Context) func(runtime.Value) {
	if ctx == nil {
		return nil
	}

	adopt, _ := ctx.Value(valueAdopterKey{}).(func(runtime.Value))

	return adopt
}

// NewOperationContext caches adoption once while preserving the context's
// cancellation lifetime. Operations without an adopter need no wrapper.
func NewOperationContext(ctx context.Context) context.Context {
	adopt := ValueAdopterFromContext(ctx)
	if adopt == nil {
		return ctx
	}

	return &operationContext{Context: ctx, adopt: adopt}
}

// OperationValueAdopter returns only an operation's cached callback. Codec
// traversal uses its own operation context, independent of host predicate contexts.
func OperationValueAdopter(ctx context.Context) func(runtime.Value) {
	operation, ok := ctx.(*operationContext)
	if !ok {
		return nil
	}

	return operation.adopt
}
