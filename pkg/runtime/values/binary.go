package values

import (
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"io"
	"io/ioutil"
)

type Binary struct {
	values []byte
}

func NewBinary(values []byte) *Binary {
	return &Binary{values}
}

func NewBinaryFrom(stream io.Reader) (*Binary, error) {
	values, err := ioutil.ReadAll(stream)

	if err != nil {
		return nil, err
	}

	return &Binary{values}, nil
}

func (b *Binary) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.values)
}

func (b *Binary) Type() core.Type {
	return core.BinaryType
}

func (b *Binary) String() string {
	return string(b.values)
}

func (b *Binary) Compare(other core.Value) int {
	// TODO: Lame comparison, need to think more about it
	switch other.Type() {
	case core.BooleanType:
		b2 := other.(*Binary)

		if b2.Length() == b.Length() {
			return 0
		}

		if b.Length() > b2.Length() {
			return 1
		}

		return -1
	default:
		return 1
	}
}

func (b *Binary) Unwrap() interface{} {
	return b.values
}

func (b *Binary) Hash() int {
	h := sha512.New()

	out, err := h.Write(b.values)

	if err != nil {
		return 0
	}

	return out
}

func (b *Binary) Clone() core.Value {
	c := make([]byte, len(b.values))

	copy(c, b.values)

	return NewBinary(c)
}

func (b *Binary) Length() Int {
	return NewInt(len(b.values))
}
