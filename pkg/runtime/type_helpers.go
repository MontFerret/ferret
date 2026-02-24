package runtime

import (
	"fmt"
	"math"
	"reflect"
)

// TypeName returns the stable name of a Type for equality/encoding purposes.
func TypeName(t Type) string {
	if t == nil {
		return ""
	}

	return t.Name()
}

func typeRank(value Value) int64 {
	if value == None || value == nil {
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
	case Binary:
		return 6
	case List:
		return 7
	case Map:
		return 8
	default:
		return math.MaxInt // unknown types last
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

// TypeOf returns the Type of a given Value, respecting any Typed overrides.
// It checks for Typed first, then known concrete types and interfaces, and if no match is found, it falls back to HostTypeOf to create a Type based on the Go type of the value using reflection.
// This allows the runtime to recognize and work with types defined in the host environment, even if they don't implement any specific interfaces or aren't known concrete types within the runtime.
// For example, if a value is of a Go struct type that doesn't implement Typed or any known interfaces, HostTypeOf will create a Type with the name of that struct type, allowing it to be used in type checks and error messages within the runtime.
func TypeOf(input Value) Type {
	if input == None || input == nil {
		return TypeNone
	}

	if typed, ok := input.(Typed); ok {
		return typed.Type()
	}

	switch input.(type) {
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
	default:
		return HostTypeOf(input)
	}
}

// HostTypeOf creates a new Type from a given value using reflection, for types not known to the runtime.
// It handles pointers, slices, and maps by unwrapping them and including their element types in the name.
// For example, a value of type *[]int would yield a Type with the name "[]int".
// This function is used as a fallback in TypeOf when a value doesn't match any known concrete types or interfaces.
// Note that the resulting Type may not have any special behavior beyond its name, and Is checks will likely fail unless the value is exactly the same type.
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

	return newHostType(pkg, name, t)
}

// IsSameType checks if two Types are the same, considering nil as a valid type that can be compared.
// It returns true if both types are nil or if they are the same non-nil type, and false otherwise.
// This function is used in IsType to quickly check for type equality before falling back to Is checks, and it allows for nil types to be compared without causing a panic.
// For example, IsSameType(nil, nil) returns true, IsSameType(nil, TypeInt) returns false, and IsSameType(TypeInt, TypeInt) returns true.
// Note that IsSameType only checks for direct type equality and does not consider interface conformance or Typed overrides, which is why IsType uses it as a first check before calling required.Is(value).
func IsSameType(a, b Type) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a == b
}

// IsType reports whether a value matches the required type.
// It uses TypeOf first (respecting Typed overrides) and then checks required.Is.
func IsType(value Value, required Type) bool {
	if required == nil {
		return false
	}

	actual := TypeOf(value)

	if IsSameType(actual, required) {
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
		if IsSameType(tid, t) || (t != nil && t.Is(value)) {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(tid, required...)
	}

	return nil
}
