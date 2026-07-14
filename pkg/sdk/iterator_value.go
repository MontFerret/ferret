package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IteratorValue exposes an Iterator host target as both a Ferret value and iterable.
type IteratorValue[T runtime.Iterator] struct {
	*HostValue[T]
}

// NewIteratorValue creates an iterator host value with a derived type.
func NewIteratorValue[T runtime.Iterator](target T) *IteratorValue[T] {
	return &IteratorValue[T]{HostValue: NewHostValue(target)}
}

// NewIteratorValueWithType creates an iterator host value with an explicit type.
func NewIteratorValueWithType[T runtime.Iterator](typeName runtime.Type, target T) *IteratorValue[T] {
	return &IteratorValue[T]{HostValue: NewHostValueWithType(typeName, target)}
}

// Iterate returns the wrapped iterator.
func (v *IteratorValue[T]) Iterate(_ context.Context) (runtime.Iterator, error) {
	return v.Target(), nil
}

// Next delegates to the wrapped iterator.
func (v *IteratorValue[T]) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	return v.Target().Next(ctx)
}

// Copy preserves the iterator wrapper's target, type, and identity.
func (v *IteratorValue[T]) Copy() runtime.Value {
	if v == nil {
		return runtime.None
	}

	return &IteratorValue[T]{HostValue: v.HostValue.copyValue()}
}
