package runtime

import (
	"fmt"
	"reflect"
)

type (
	Type interface {
		fmt.Stringer
		// Name returns the name of the type represented by this Type.
		Name() string
		// Is checks if the provided value is of the type represented by this Type.
		Is(Value) bool
	}

	// TypeMatcher is a function type that takes a Value and returns a boolean indicating whether the value matches a specific type.
	TypeMatcher func(Value) bool

	// Typed is an interface that can be implemented by any value to provide its type information.
	Typed interface {
		Type() Type
	}

	// runtimeType is a struct that represents a runtime-defined type with a name and an assertion function.
	// name is the string representation of the type's name.
	// assert is a TypeMatcher function used to determine if a value matches this type.
	runtimeType struct {
		name   string
		assert TypeMatcher
	}

	// hostType is a struct that represents a type defined in the host environment (e.g., Go types).
	// name is the string representation of the type's name, typically in the format "package.TypeName".
	hostType struct {
		name string
		typ  reflect.Type
	}
)

const (
	UnknownTypeName = "Unknown"
)

func NewType(name string, assert TypeMatcher) Type {
	return &runtimeType{name: name, assert: assert}
}

func (rt *runtimeType) Name() string {
	if rt == nil {
		return UnknownTypeName
	}

	return rt.name
}

func (rt *runtimeType) String() string {
	if rt == nil {
		return UnknownTypeName
	}

	return rt.name
}

func (rt *runtimeType) Is(v Value) bool {
	if rt == nil || rt.assert == nil {
		return false
	}

	return rt.assert(v)
}

// hewHostType creates a new Type for a host-defined type (e.g., a Go type).
// The pkg parameter is the package path of the type, and the name parameter is the name of the type.
func hewHostType(pkg, name string, typ reflect.Type) Type {
	if pkg == "" {
		return &hostType{name: name}
	}

	if name == "" {
		panic("name is empty")
	}

	return &hostType{name: pkg + "." + name, typ: typ}
}

func (ht *hostType) Name() string {
	if ht == nil {
		return UnknownTypeName
	}

	return ht.name
}

func (ht *hostType) String() string {
	if ht == nil {
		return UnknownTypeName
	}

	return ht.name
}

func (ht *hostType) Is(v Value) bool {
	if ht.typ == nil {
		return false
	}

	// If v is nil, reflect.TypeOf(v) is nil.
	// Only "matches" if the target is a nil-able kind (ptr/map/slice/func/chan/interface)
	// and v is a typed nil - but `any(nil)` is untyped nil, so treat as false here.
	if v == nil {
		return false
	}

	t := reflect.TypeOf(v)

	// Exact type match
	if t == ht.typ {
		return true
	}

	// If target is an interface, this is the key: does v implement it?
	// More generally, this answers: can a value of type t be assigned to ht.typ?
	if t.AssignableTo(ht.typ) {
		return true
	}

	// If v is *T and target is T, check element type.
	if t.Kind() == reflect.Ptr && t.Elem().AssignableTo(ht.typ) {
		return true
	}

	return false
}
