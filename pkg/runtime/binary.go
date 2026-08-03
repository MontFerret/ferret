package runtime

import (
	"context"
	"hash/fnv"
	"io"
)

type Binary []byte

func NewBinary(values []byte) Binary {
	return values
}

func NewBinaryFrom(stream io.Reader) (Binary, error) {
	values, err := io.ReadAll(stream)

	if err != nil {
		return nil, err
	}

	return values, nil
}

func (b Binary) Type() Type {
	return TypeBinary
}

func (b Binary) String() string {
	return string(b)
}

func (b Binary) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeBinary.Name()))
	h.Write([]byte(":"))
	h.Write(b)

	return h.Sum64()
}

func (b Binary) Copy() Value {
	c := make([]byte, len(b))

	copy(c, b)

	return NewBinary(c)
}

func (b Binary) Length(_ context.Context) (Int, error) {
	return Int(len(b)), nil
}

func (b Binary) Unwrap() any {
	return []byte(b)
}
