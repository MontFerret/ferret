package mem

import (
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ManagedResource wraps a runtime.Value that also implements io.Closer,
// providing a guaranteed pointer-comparable identity for the VM's ownership
// tracking maps. Only non-comparable closers need wrapping; pointer-type
// closers (the vast majority) pass through unchanged and serve as their
// own map keys.
type ManagedResource struct {
	inner  runtime.Value
	closer io.Closer
}

// AdoptCloser ensures that val can safely participate in the VM's
// ownership tracking maps. If val implements io.Closer and its dynamic
// type is already comparable, it is returned unchanged. If val is a
// non-comparable closer, it is wrapped in a *ManagedResource to provide
// pointer-comparable identity. Non-closers pass through unchanged.
//
// The reflect.Comparable check runs once here at the production boundary,
// removing it from the hot-path functions TrackedCloserOf and CloserSet.Add
// that are called repeatedly during execution.
func AdoptCloser(val runtime.Value) runtime.Value {
	if val == nil {
		return val
	}

	if _, ok := val.(*ManagedResource); ok {
		return val
	}

	closer, ok := val.(io.Closer)
	if !ok || closer == nil {
		return val
	}

	// Pointer-type closers are always comparable and can serve as their
	// own map key without wrapping. Only wrap non-comparable closers.
	typ := reflect.TypeOf(closer)
	if typ != nil && typ.Comparable() {
		return val
	}

	return &ManagedResource{
		inner:  val,
		closer: closer,
	}
}

// Close delegates to the underlying closer.
func (m *ManagedResource) Close() error {
	return m.closer.Close()
}

// String delegates to the inner value.
func (m *ManagedResource) String() string {
	return m.inner.String()
}

// Hash delegates to the inner value.
func (m *ManagedResource) Hash() uint64 {
	return m.inner.Hash()
}

// Copy returns a copy of the inner value. Copies do not carry resource
// ownership, so the result is not wrapped in a ManagedResource.
func (m *ManagedResource) Copy() runtime.Value {
	return m.inner.Copy()
}

// Unwrap returns the original runtime.Value before wrapping.
func (m *ManagedResource) Unwrap() runtime.Value {
	return m.inner
}
