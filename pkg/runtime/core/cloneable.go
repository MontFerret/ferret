package core

// Cloneable represents an interface of a value that can be cloned.
// The difference between Copy and Clone is that Copy returns a shallow copy of the value
// and Clone returns a deep copy of the value.
type Cloneable interface {
	Value
	Clone() Cloneable
}
