// Package encodingownership carries synchronous resource adoption during native
// output materialization without making ownership part of the public codec API.
package encodingownership

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type valueAdopterKey struct{}

// WithValueAdopter lends adopt to codecs for this materialization only. Codecs
// must not retain it or call it after returning. It registers ownership without
// encoding, closing resources, or observing cancellation. The callback's second
// argument requests discovery of concrete stored descendants rather than direct
// registration.
func WithValueAdopter(ctx context.Context, adopt func(runtime.Value, bool)) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, valueAdopterKey{}, adopt)
}

// ValueAdopterFromContext returns the operation's adoption callback, if present.
func ValueAdopterFromContext(ctx context.Context) func(runtime.Value, bool) {
	if ctx == nil {
		return nil
	}

	adopt, _ := ctx.Value(valueAdopterKey{}).(func(runtime.Value, bool))

	return adopt
}
