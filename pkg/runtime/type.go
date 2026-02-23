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

func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := TypeOf(value)

	for _, t := range required {
		if tid == t {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(tid, required...)
	}

	return nil
}

// PairValueType is a supporting
// structure that used in validateValueTypePairs.
type PairValueType struct {
	Value Value
	Types []Type
}

// ValidateValueTypePairs validate pairs of
// Values and Types.
// Returns error when type didn't match
func ValidateValueTypePairs(pairs ...PairValueType) error {
	var err error

	for idx, pair := range pairs {
		err = ValidateType(pair.Value, pair.Types...)

		if err != nil {
			return fmt.Errorf("pair %d: %w", idx, err)
		}
	}

	return nil
}
