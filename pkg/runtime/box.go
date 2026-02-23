package runtime

import (
	"fmt"
	"hash/fnv"
	"reflect"
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

func (b *Box[T]) Type() Type {
	t := reflect.TypeOf(b.Value)

	return NewType(t.PkgPath(), t.Name())
}

func (b *Box[T]) String() string {
	return fmt.Sprintf("Box[%s]: %s", b.Type(), b.Value)
}

func (b *Box[T]) Unwrap() any {
	return b.Value
}

func (b *Box[T]) Hash() uint64 {
	h := fnv.New64a()

	_, _ = h.Write([]byte("box:"))
	_, _ = h.Write([]byte(b.String()))

	return h.Sum64()
}

func (b *Box[T]) Copy() Value {
	return &Box[T]{Value: b.Value}
}
