package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/wI2L/jettison"
)

type Box[T any] struct {
	Value T
}

func (v *Box[T]) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(v.Value, jettison.NoHTMLEscaping())
}

func (v *Box[T]) String() string {
	return "[Box]"
}

func (v *Box[T]) Unwrap() interface{} {
	return v.Value
}

func (v *Box[T]) Hash() uint64 {
	panic("not supported")
}

func (v *Box[T]) Copy() runtime.Value {
	return &Box[T]{Value: v.Value}
}
