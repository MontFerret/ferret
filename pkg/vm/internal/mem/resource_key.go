package mem

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ResourceKey is the canonical identity used in all ownership-tracking maps.
// Exactly one field is non-zero:
//   - ID is set (non-zero) when the value implements runtime.Resource.
//     Two values with the same ResourceID share this key regardless of which
//     Go object holds them, satisfying the stable-identity dedup contract.
//   - Closer is set (non-nil) for plain io.Closer values; identity is the
//     pointer itself (compatibility path).
//
// ResourceKey is a comparable struct and is safe to use as a map key.
type ResourceKey struct {
	Closer io.Closer
	ID     uint64
}

// ResourceKeyOf derives the ownership key and associated closer from val.
//
// Primary path (runtime.Resource): returns ResourceKey{ID: val.ResourceID()}
// and the resource as the closer. Two values with the same ResourceID produce
// equal keys.
//
// Compatibility path (plain io.Closer): returns ResourceKey{Closer: val} and
// the closer. Identity is the pointer value.
//
// Returns (zero, nil, false) if val does not implement io.Closer.
func ResourceKeyOf(val runtime.Value) (ResourceKey, io.Closer, bool) {
	if res, ok := val.(runtime.Resource); ok {
		return ResourceKey{ID: res.ResourceID()}, res, true
	}

	if closer, ok := val.(io.Closer); ok && closer != nil {
		return ResourceKey{Closer: closer}, closer, true
	}

	return ResourceKey{}, nil, false
}
