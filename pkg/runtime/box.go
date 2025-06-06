package runtime

import (
	"hash/fnv"

	"github.com/wI2L/jettison"
)

// Box is a generic wrapper for any value type.
type Box[T any] struct {
	Value T
}

func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		Value: value,
	}
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
	h := fnv.New64a()

	_, _ = h.Write([]byte("box:"))

	data, err := v.MarshalJSON()

	if err != nil {
		// TODO: Panic?
		return 0
	}

	_, _ = h.Write(data)

	return h.Sum64()
}

func (v *Box[T]) Copy() Value {
	return &Box[T]{Value: v.Value}
}
