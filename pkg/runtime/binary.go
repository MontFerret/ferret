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

// ToBinary attempts to convert an arbitrary Value into a Binary type.
// If the input is already a Binary, it returns it directly.
// For other types, it converts the input to its string representation and then to a byte slice to create a new Binary.
// This allows for flexible conversion of various Value types into a Binary format, using their string representation as the basis for the binary data.
func ToBinary(input Value) Binary {
	bin, ok := input.(Binary)

	if ok {
		return bin
	}

	return NewBinary([]byte(input.String()))
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
