package session

import (
	"context"
	"slices"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// outputValueAdopter discovers only resources already present in yielded values.
// Host capabilities and lazy iterables remain with the encoder's normal traversal.
type outputValueAdopter struct {
	ctx     context.Context
	result  *vm.Result
	seen    map[runtime.Value]struct{}
	pending []runtime.Value
}

func (a *outputValueAdopter) Adopt(value runtime.Value) {
	switch value.(type) {
	case *runtime.Array, *runtime.Object:
	default:
		a.result.AdoptValue(value)

		return
	}

	// A host may update and yield the same container again. Cycle protection is
	// per yield so that later stored children still acquire cleanup ownership.
	defer func() { clear(a.seen) }()
	a.pending = append(a.pending, value)
	for len(a.pending) > 0 {
		last := len(a.pending) - 1
		item := a.pending[last]
		a.pending[last] = nil
		a.pending = a.pending[:last]
		a.result.AdoptValue(item)
		start := len(a.pending)

		// Only these concrete, materialized containers have non-blocking storage
		// access independent of cancellation. Never inspect a host collection.
		switch container := item.(type) {
		case *runtime.Array:
			if container == nil || a.visited(container) {
				continue
			}

			_ = container.ForEach(a.ctx, func(_ context.Context, child runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
				a.pending = append(a.pending, child)

				return true, nil
			})
		case *runtime.Object:
			if container == nil || a.visited(container) {
				continue
			}

			_ = container.ForEach(a.ctx, func(_ context.Context, child, _ runtime.Value) (runtime.Boolean, error) {
				a.pending = append(a.pending, child)

				return true, nil
			})
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
