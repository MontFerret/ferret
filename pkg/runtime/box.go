package runtime

import (
	"fmt"
	"hash/fnv"
	"reflect"

	"github.com/wI2L/jettison"
)

// Box is a generic container that holds a single value of type T.
// Useful for wrapping values that do not implement the Value interface, allowing them to be used in contexts where a Value is expected.
type Box[T any] struct {
	Value T
}

// NewBox creates a new Box containing the provided value.
func NewBox[T any](value T) *Box[T] {
	return &Box[T]{
		Value: value,
	}
}

func (v *Box[T]) Type() Type {
	t := reflect.TypeOf(v.Value)

	return NewType(t.PkgPath(), t.Name())
}

func (v *Box[T]) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(v.Value, jettison.NoHTMLEscaping())
}

func (v *Box[T]) String() string {
	return fmt.Sprintf("Box[%s]", v.Type())
}

func (v *Box[T]) Unwrap() any {
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
