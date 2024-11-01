package values

import (
	"hash/fnv"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Binary []byte

func NewBinary(values []byte) Binary {
	return Binary(values)
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

func (b Binary) Type() core.Type {
	return types.Binary
}

func (b Binary) String() string {
	return string(b)
}

func (b Binary) Compare(other core.Value) int64 {
	otherBin, ok := other.(Binary)

	if !ok {
		return types.Compare(types.Binary, core.Reflect(other))
	}

	// TODO: Lame comparison, need to think more about it
	// Maybe we should do a byte by byte comparison?
	if otherBin.Length() == b.Length() {
		return 0
	}

	if b.Length() > otherBin.Length() {
		return 1
	}

	return -1
}

func (b Binary) Unwrap() interface{} {
	return []byte(b)
}

func (b Binary) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Binary.String()))
	h.Write([]byte(":"))
	h.Write(b)

	return h.Sum64()
}

func (b Binary) Copy() core.Value {
	c := make([]byte, len(b))

	copy(c, b)

	return NewBinary(c)
}

func (b Binary) Length() Int {
	return NewInt(len(b))
}
