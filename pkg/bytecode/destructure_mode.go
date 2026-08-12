package bytecode

// DestructureMode identifies the capability required by a destructuring assertion.
type DestructureMode byte

const (
	DestructureModeInvalid DestructureMode = iota
	DestructureModeObject
	DestructureModeArray
)

func (m DestructureMode) String() string {
	switch m {
	case DestructureModeObject:
		return "Object"
	case DestructureModeArray:
		return "Array"
	default:
		return "Invalid"
	}
}

func (m DestructureMode) Valid() bool {
	return m == DestructureModeObject || m == DestructureModeArray
}
