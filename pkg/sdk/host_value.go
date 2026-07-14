package sdk

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var hostValueIdentity atomic.Uint64

// HostValue exposes opaque host data as a Ferret value without claiming optional capabilities.
// Its hash uses stable wrapper identity rather than host pointer reflection.
type HostValue[T any] struct {
	target   T
	typeName runtime.Type
	identity uint64
}

// NewHostValue creates an opaque host value with a type derived from T.
func NewHostValue[T any](target T) *HostValue[T] {
	return &HostValue[T]{
		target:   target,
		identity: hostValueIdentity.Add(1),
	}
}

// NewHostValueWithType creates an opaque host value with explicit Ferret type information.
func NewHostValueWithType[T any](typeName runtime.Type, target T) *HostValue[T] {
	value := NewHostValue(target)
	value.typeName = typeName

	return value
}

// Target returns the wrapped host value.
func (v *HostValue[T]) Target() T {
	if v == nil {
		var zero T
		return zero
	}

	return v.target
}

// Type returns the explicit type, the target's runtime type, or its derived Go host type.
func (v *HostValue[T]) Type() runtime.Type {
	if v == nil {
		return runtime.TypeNone
	}

	if v.typeName != nil {
		return v.typeName
	}

	if typed, ok := any(v.target).(runtime.Typed); ok {
		return typed.Type()
	}
	if reflect.TypeOf(v.target) == nil {
		return runtime.TypeNone
	}

	return runtime.HostTypeOf(v.target)
}

// Unwrap returns the underlying host representation.
func (v *HostValue[T]) Unwrap() any {
	if v == nil {
		return nil
	}

	if unwrappable, ok := any(v.target).(runtime.Unwrappable); ok {
		return unwrappable.Unwrap()
	}

	return v.target
}

// MarshalJSON delegates to the target when possible and otherwise encodes the host value.
func (v *HostValue[T]) MarshalJSON() ([]byte, error) {
	if v == nil || v.targetNil() {
		return []byte("null"), nil
	}

	if marshaler, ok := any(v.target).(json.Marshaler); ok {
		return marshaler.MarshalJSON()
	}

	return jettison.MarshalOpts(v.target, jettison.NoHTMLEscaping())
}

// String returns the target's string representation without exposing capabilities.
func (v *HostValue[T]) String() string {
	if v == nil || v.targetNil() {
		return runtime.None.String()
	}

	if stringer, ok := any(v.target).(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprint(v.target)
}

// Hash returns the stable identity shared by copies of this wrapper.
func (v *HostValue[T]) Hash() uint64 {
	if v == nil {
		return runtime.None.Hash()
	}

	return v.identity
}

// Copy creates a shallow wrapper copy that preserves target, type, and identity.
func (v *HostValue[T]) Copy() runtime.Value {
	if v == nil {
		return runtime.None
	}

	return v.copyValue()
}

func (v *HostValue[T]) copyValue() *HostValue[T] {
	if v == nil {
		return nil
	}

	return &HostValue[T]{
		target:   v.target,
		typeName: v.typeName,
		identity: v.identity,
	}
}

func (v *HostValue[T]) setTarget(target T) {
	v.target = target
}

func (v *HostValue[T]) targetNil() bool {
	if v == nil {
		return true
	}

	target := reflect.ValueOf(v.target)
	if !target.IsValid() {
		return true
	}

	switch target.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return target.IsNil()
	default:
		return false
	}
}
