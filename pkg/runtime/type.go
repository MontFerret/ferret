package runtime

import "reflect"

const (
	TypeNone     = "none"
	TypeBoolean  = "boolean"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeString   = "string"
	TypeDateTime = "date_time"
	TypeList     = "list"
	TypeMap      = "map"
	TypeBinary   = "binary"

	// Interfaces
	TypeIterable   = "iterable"
	TypeIterator   = "iterator"
	TypeMeasurable = "measurable"
	TypeComparable = "comparable"
	TypeCloneable  = "cloneable"
	TypeSortable   = "sortable"
	TypeObservable = "observable"
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

func Reflect(input Value) string {
	if input == None || input == nil {
		return TypeNone
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
	case Observable:
		return TypeObservable
	default:
		return reflect.TypeOf(input).String()
	}
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
