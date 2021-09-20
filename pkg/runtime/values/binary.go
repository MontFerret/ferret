package values

import (
	"hash/fnv"
	"io"
	"io/ioutil"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Binary []byte

func NewBinary(values []byte) Binary {
	return Binary(values)
}

func NewBinaryFrom(stream io.Reader) (Binary, error) {
	values, err := ioutil.ReadAll(stream)

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
	if other.Type() == types.Binary {
		// TODO: Lame comparison, need to think more about it
		b2 := other.(*Binary)

		if b2.Length() == b.Length() {
			return 0
		}

		if b.Length() > b2.Length() {
			return 1
		}

		return -1
	}

	return types.Compare(types.Binary, other.Type())
}

func (b Binary) Unwrap() interface{} {
	return []byte(b)
}

func (b Binary) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(b.Type().String()))
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
