package core

type Cloneable interface {
	Value
	Clone() Cloneable
}
