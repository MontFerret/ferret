package runtime

import (
	"fmt"
	"reflect"
)

type (
	Type string

	// Typed is an interface that can be implemented by any value to provide its type information.
	Typed interface {
		Type() Type
	}
)

// NewType creates a new Type from a given package and name.
func NewType(pkg, name string) Type {
	if pkg == "" {
		return Type(name)
	}

	if name == "" {
		panic("name is empty")
	}

	return Type(pkg + "." + name)
}

func (t Type) String() string {
	return string(t)
}

const (
	// Actual types
	TypeNone     = Type("None")
	TypeBoolean  = Type("Boolean")
	TypeInt      = Type("Int")
	TypeFloat    = Type("Float")
	TypeString   = Type("String")
	TypeDateTime = Type("DateTime")
	TypeArray    = Type("Array")
	TypeObject   = Type("Object")
	TypeBinary   = Type("Binary")
	TypeRegexp   = Type("Regexp")
	TypeQuery    = Type("Query")

	// Interfaces
	TypeCollection     = Type("Collection")
	TypeList           = Type("List")
	TypeMap            = Type("Map")
	TypeIndexReadable  = Type("IndexReadable")
	TypeIndexRemovable = Type("IndexRemovable")
	TypeIndexWritable  = Type("IndexWritable")
	TypeKeyReadable    = Type("KeyReadable")
	TypeKeyWritable    = Type("KeyWritable")
	TypeKeyRemovable   = Type("KeyRemovable")
	TypeValueRemovable = Type("ValueRemovable")
	TypeAppendable     = Type("Appendable")
	TypeContainable    = Type("Containable")
	TypeIterable       = Type("Iterable")
	TypeIterator       = Type("Iterator")
	TypeMeasurable     = Type("Measurable")
	TypeComparable     = Type("Comparable")
	TypeCloneable      = Type("Cloneable")
	TypeSortable       = Type("Sortable")
	TypeDispatchable   = Type("Dispatchable")
	TypeObservable     = Type("Observable")
	TypeQueryable      = Type("Queryable")
)

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
	case Query:
		return TypeQuery
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
		return ReflectTypeOf(input)
	}
}

// ReflectTypeOf creates a new Type from a given value using reflection.
func ReflectTypeOf(input any) Type {
	return typeOfBuiltin(reflect.TypeOf(input))
}

func typeOfBuiltin(t reflect.Type) Type {
	var name, pkg string

	switch t.Kind() {
	case reflect.Ptr:
		return typeOfBuiltin(t.Elem())
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

	return NewType(pkg, name)
}

// IsType reports whether a value matches the required type.
// It uses TypeOf first (respecting Typed overrides) and then checks interface conformance.
func IsType(value Value, required Type) bool {
	actual := TypeOf(value)

	return isTypeWithActual(value, actual, required)
}

func isTypeWithActual(value Value, actual Type, required Type) bool {
	if actual == required {
		return true
	}

	switch required {
	case TypeCollection:
		_, ok := value.(Collection)
		return ok
	case TypeList:
		_, ok := value.(List)
		return ok
	case TypeMap:
		_, ok := value.(Map)
		return ok
	case TypeIndexReadable:
		_, ok := value.(IndexReadable)
		return ok
	case TypeIndexWritable:
		_, ok := value.(IndexWritable)
		return ok
	case TypeIndexRemovable:
		_, ok := value.(IndexRemovable)
		return ok
	case TypeKeyReadable:
		_, ok := value.(KeyReadable)
		return ok
	case TypeKeyWritable:
		_, ok := value.(KeyWritable)
		return ok
	case TypeKeyRemovable:
		_, ok := value.(KeyRemovable)
		return ok
	case TypeValueRemovable:
		_, ok := value.(ValueRemovable)
		return ok
	case TypeAppendable:
		_, ok := value.(Appendable)
		return ok
	case TypeContainable:
		_, ok := value.(Containable)
		return ok
	case TypeIterable:
		_, ok := value.(Iterable)
		return ok
	case TypeIterator:
		_, ok := value.(Iterator)
		return ok
	case TypeMeasurable:
		_, ok := value.(Measurable)
		return ok
	case TypeComparable:
		_, ok := value.(Comparable)
		return ok
	case TypeCloneable:
		_, ok := value.(Cloneable)
		return ok
	case TypeSortable:
		_, ok := value.(Sortable)
		return ok
	case TypeDispatchable:
		_, ok := value.(Dispatchable)
		return ok
	case TypeObservable:
		_, ok := value.(Observable)
		return ok
	case TypeQueryable:
		_, ok := value.(Queryable)
		return ok
	default:
		return false
	}
}

// ValidateType checks whether a value matches any of the required types and returns an error if it doesn't.
// It uses IsType, so interface conformance and Typed overrides are accepted.
func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := TypeOf(value)

	for _, t := range required {
		if isTypeWithActual(value, tid, t) {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(tid, required...)
	}

	return nil
}
