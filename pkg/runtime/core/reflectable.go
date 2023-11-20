package core

import "reflect"

// Reflectable Represents a value can reflect its own type.
type Reflectable interface {
	Type() Type
}

func Reflect(input Value) Type {
	reflectable, ok := input.(Reflectable)

	if ok {
		return reflectable.Type()
	}

	typeDesc := reflect.TypeOf(input)

	return NewType(typeDesc.PkgPath(), typeDesc.Name())
}
