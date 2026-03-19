package mem

import (
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// ResourceKey is the canonical identity used in all ownership-tracking maps.
	// Exactly one field is non-zero:
	//   - ID is set (non-zero) when the value implements runtime.Resource.
	//     Two values with the same ResourceID share this key regardless of which
	//     Go object holds them, satisfying the stable-identity dedup contract.
	//   - Closer is set (non-nil) for plain comparable io.Closer values; identity
	//     is the interface value itself (compatibility path).
	//
	// ResourceKey is a comparable struct and is safe to use as a map key.
	ResourceKey struct {
		Closer io.Closer
		ID     uint64
	}

	// Untracked is implemented by internal VM scaffolding values that never
	// participate in resource ownership tracking even if they are produced by the
	// VM and shuffled through registers.
	Untracked interface {
		VMUntracked()
	}
)

// CanTrackValue cheaply rejects values that are known to never participate in
// direct VM ownership tracking. It is intentionally conservative: returning
// true only means the caller should fall back to the full interface checks in
// ResourceKeyOf.
func CanTrackValue(val runtime.Value) bool {
	if val == nil || val == runtime.None {
		return false
	}

	switch val.(type) {
	case runtime.Boolean, runtime.Int, runtime.Float, runtime.String, *runtime.Array, *runtime.Object:
		return false
	default:
		_, ok := val.(Untracked)

		return !ok
	}
}

// ResourceKeyOf derives the ownership key and associated closer from val.
//
// Primary path (runtime.Resource): returns ResourceKey{ID: val.ResourceID()}
// and the resource as the closer. Two values with the same ResourceID produce
// equal keys.
//
// Compatibility path (plain comparable io.Closer): returns
// ResourceKey{Closer: val} and the closer. Identity is the interface value.
// Plain non-comparable closers must implement runtime.Resource if they need VM
// ownership tracking.
//
// Returns (zero, nil, false) if val does not implement io.Closer or cannot be
// tracked safely.
func ResourceKeyOf(val runtime.Value) (ResourceKey, io.Closer, bool) {
	if !CanTrackValue(val) {
		return ResourceKey{}, nil, false
	}

	closer, ok := val.(io.Closer)
	if !ok || closer == nil {
		return ResourceKey{}, nil, false
	}

	key, ok := closerResourceKey(closer)
	if !ok {
		return ResourceKey{}, nil, false
	}

	return key, closer, true
}

func closerResourceKey(closer io.Closer) (ResourceKey, bool) {
	if closer == nil {
		return ResourceKey{}, false
	}

	if res, ok := closer.(runtime.Resource); ok {
		return ResourceKey{ID: res.ResourceID()}, true
	}

	typ := reflect.TypeOf(closer)
	if typ == nil || !typ.Comparable() {
		return ResourceKey{}, false
	}

	return ResourceKey{Closer: closer}, true
}
