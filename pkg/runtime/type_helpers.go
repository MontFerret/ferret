package runtime

import (
	"fmt"
	"reflect"
)

// TypeName returns the stable name of a Type for equality/encoding purposes.
func TypeName(t Type) string {
	if t == nil {
		return ""
	}

	return t.Name()
}

// SameType reports whether two Types share the same stable name.
func SameType(a, b Type) bool {
	return TypeName(a) == TypeName(b)
}

func typeRank(value Value) int64 {
	if value == None {
		return 0
	}

	switch value.(type) {
	case Boolean:
		return 1
	case Int:
		return 2
	case Float:
		return 3
	case String:
		return 4
	case DateTime:
		return 5
	case List:
		return 6
	case Map:
		return 7
	case Binary:
		return 8
	default:
		return -1
	}
}

// CompareTypes compares the types of two values and returns -1 if a < b, 0 if a == b, and 1 if a > b.
func CompareTypes(a, b Value) int64 {
	aRank := typeRank(a)
	bRank := typeRank(b)

	if aRank == bRank {
		return 0
	}

	if aRank < bRank {
		return -1
	}

	return 1
}

// TypeOf returns the runtime Type of a given value.
// It prefers known concrete types and interfaces; if none match and the value
// implements Typed, its Type is returned. Otherwise it falls back to reflection.
func TypeOf(input Value) Type {
	if input == None || input == nil {
		return TypeNone
	}

	switch v := input.(type) {
	case Boolean:
		return TypeBoolean
	case Int:
		return TypeInt
	case Float:
		return TypeFloat
	case String:
		return TypeString
	case DateTime:
		return TypeDateTime
	case *Array:
		return TypeArray
	case *Object:
		return TypeObject
	case ObjectLike:
		return TypeObject
	case List:
		return TypeList
	case Map:
		return TypeMap
	case Binary:
		return TypeBinary
	case Iterable:
		return TypeIterable
	case Iterator:
		return TypeIterator
	case Measurable:
		return TypeMeasurable
	case Dispatchable:
		return TypeDispatchable
	case Observable:
		return TypeObservable
	case Queryable:
		return TypeQueryable
	case Typed:
		return v.Type()
	default:
		return HostTypeOf(input)
	}
}

// HostTypeOf creates a new Type from a given value using reflection, ignoring pointers and slices/maps of known types.
// This is used as a fallback for values that don't match any known concrete types or interfaces and don't implement Typed.
// Unlike TypeOf, it doesn't check for Typed overrides and always uses reflection, so it's suitable for host values that may not implement Typed.
func HostTypeOf(input any) Type {
	return typeFromReflect(reflect.TypeOf(input))
}

func typeFromReflect(t reflect.Type) Type {
	var name, pkg string

	switch t.Kind() {
	case reflect.Ptr:
		return typeFromReflect(t.Elem())
	case reflect.Slice:
		elem := t.Elem()
		name = fmt.Sprintf("[]%s", elem.Name())
		pkg = elem.PkgPath()
	case reflect.Map:
		key := t.Key()
		value := t.Elem()
		name = fmt.Sprintf("map[%s]%s", key.Name(), value.Name())
		pkg = key.PkgPath()
	default:
		name = t.Name()
		pkg = t.PkgPath()
	}

	return NewHostType(pkg, name)
}

// IsType reports whether a value matches the required type.
// It uses TypeOf first (respecting Typed overrides) and then checks required.Is.
func IsType(value Value, required Type) bool {
	if required == nil {
		return false
	}

	actual := TypeOf(value)

	if SameType(actual, required) {
		return true
	}

	return required.Is(value)
}

// ValidateType checks whether a value matches any of the required types and returns an error if it doesn't.
// It uses IsType, so interface conformance and Typed overrides are accepted.
func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := TypeOf(value)

	for _, t := range required {
		if SameType(tid, t) || (t != nil && t.Is(value)) {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(tid, required...)
	}

	return nil
}
