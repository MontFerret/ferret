package runtime

import (
	"hash/fnv"
	"io"

	"github.com/wI2L/jettison"
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

func (b Binary) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts([]byte(b),
		jettison.NoStringEscaping(),
		jettison.NoCompact(),
	)
}

func (b Binary) String() string {
	return string(b)
}

func (b Binary) Unwrap() interface{} {
	return []byte(b)
}

func (b Binary) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeBinary))
	h.Write([]byte(":"))
	h.Write(b)

	return h.Sum64()
}

func (b Binary) Copy(_ Context) (Value, error) {
	c := make([]byte, len(b))

	copy(c, b)

	return NewBinary(c), nil
}

func (b Binary) Length(_ Context) (Int, error) {
	return Int(len(b)), nil
}

func (b Binary) Compare(_ Context, other Value) int64 {
	otherBin, ok := other.(Binary)

	if !ok {
		return CompareTypes(b, other)
	}

	size := len(b)
	otherSize := len(otherBin)

	if size > otherSize {
		return 1
	} else if size < otherSize {
		return -1
	}

	for i := 0; i < size; i++ {
		if b[i] > otherBin[i] {
			return 1
		} else if b[i] < otherBin[i] {
			return -1
		}
	}

	return 0
}
