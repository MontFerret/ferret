package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// IterableValue exposes an Iterable host target as a Ferret value.
type IterableValue[T runtime.Iterable] struct {
	*HostValue[T]
}

// NewIterableValue creates an iterable host value with a derived type.
func NewIterableValue[T runtime.Iterable](target T) *IterableValue[T] {
	return &IterableValue[T]{HostValue: NewHostValue(target)}
}

// NewIterableValueWithType creates an iterable host value with an explicit type.
func NewIterableValueWithType[T runtime.Iterable](typeName runtime.Type, target T) *IterableValue[T] {
	return &IterableValue[T]{HostValue: NewHostValueWithType(typeName, target)}
}

// Iterate delegates iteration to the wrapped target.
func (v *IterableValue[T]) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return v.Target().Iterate(ctx)
}

// Copy preserves the iterable wrapper's target, type, and identity.
func (v *IterableValue[T]) Copy() runtime.Value {
	if v == nil {
		return runtime.None
	}

	return &IterableValue[T]{HostValue: v.HostValue.copyValue()}
}
