package core

var cloneableTypes = map[Type]bool{
	ArrayType:  true,
	ObjectType: true,
}

type Cloneable interface {
	Value
	Clone() Cloneable
}
