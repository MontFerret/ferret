package core

const (
	TypeNone     = "none"
	TypeBoolean  = "boolean"
	TypeInt      = "int"
	TypeFloat    = "float"
	TypeString   = "string"
	TypeDateTime = "date_time"
	TypeList     = "list"
	TypeSet      = "set"
	TypeMap      = "map"
	TypeBinary   = "binary"
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
	case Set:
		return 7
	case Map:
		return 8
	case Binary:
		return 9
	default:
		return -1
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
