package session

import (
	"context"
	"slices"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// outputValueAdopter discovers resources already present in concrete stored values.
// Host capabilities and lazy iterables remain with the encoder's normal traversal.
type outputValueAdopter struct {
	result  *vm.Result
	seen    map[runtime.Value]struct{}
	pending []runtime.Value
}

func (a *outputValueAdopter) Adopt(value runtime.Value, stored bool) {
	if !stored {
		a.result.AdoptValue(value)

		return
	}

	switch value.(type) {
	case *runtime.Array, *runtime.Object:
	default:
		a.result.AdoptValue(value)

		return
	}

	// Hosts and hooks may update the same container between discoveries. Reset
	// cycle protection so later stored children still acquire cleanup ownership.
	defer func() { clear(a.seen) }()
	a.pending = append(a.pending, value)

	for len(a.pending) > 0 {
		last := len(a.pending) - 1
		item := a.pending[last]
		a.pending[last] = nil
		a.pending = a.pending[:last]
		start := len(a.pending)

		// Only these concrete, materialized containers have non-blocking storage
		// access independent of cancellation. Their ForEach methods only return
		// predicate errors; this always-successful predicate cannot fail. Never
		// inspect a host collection or a lazy iterable.
		switch container := item.(type) {
		case *runtime.Array:
			if container == nil || a.visited(container) {
				continue
			}

			_ = container.ForEach(context.Background(), func(_ context.Context, child runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
				a.pending = append(a.pending, child)

				return true, nil
			})
		case *runtime.Object:
			if container == nil || a.visited(container) {
				continue
			}

			_ = container.ForEach(context.Background(), func(_ context.Context, child, _ runtime.Value) (runtime.Boolean, error) {
				a.pending = append(a.pending, child)

				return true, nil
			})
		default:
			a.result.AdoptValue(item)
		}

		// Keep the discovered children's cleanup order consistent with traversal.
		slices.Reverse(a.pending[start:])
	}
}

func (a *outputValueAdopter) visited(container runtime.Value) bool {
	if a.seen == nil {
		a.seen = make(map[runtime.Value]struct{})
	}

	if _, ok := a.seen[container]; ok {
		return true
	}

	a.seen[container] = struct{}{}

	return false
}
