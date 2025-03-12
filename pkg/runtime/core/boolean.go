package core

import (
	"context"
	"hash/fnv"
	"strings"

	"github.com/wI2L/jettison"
)

type Boolean bool

const (
	False = Boolean(false)
	True  = Boolean(true)
)

var (
	hashTrue  = booleanHash(True)
	hashFalse = booleanHash(False)
)

func booleanHash(val Boolean) uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeBoolean))
	h.Write([]byte(val.String()))

	return h.Sum64()
}

func NewBoolean(input bool) Boolean {
	return Boolean(input)
}

func ParseBoolean(input interface{}) (Boolean, error) {
	b, ok := input.(bool)

	if ok {
		if b {
			return True, nil
		}

		return False, nil
	}

	s, ok := input.(string)

	if ok {
		return strings.ToLower(s) == "true", nil
	}

	return False, Error(ErrInvalidType, "expected 'bool'")
}

func MustParseBoolean(input interface{}) Boolean {
	res, err := ParseBoolean(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Boolean) Type() string {
	return TypeBoolean
}

func (t Boolean) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(bool(t), jettison.NoHTMLEscaping())
}

func (t Boolean) String() string {
	if t {
		return "true"
	}

	return "false"
}

func (t Boolean) Unwrap() interface{} {
	return bool(t)
}

func (t Boolean) Hash() uint64 {
	if t {
		return hashTrue
	} else {
		return hashFalse
	}
}

func (t Boolean) Copy() Value {
	return t
}

func (t Boolean) Compare(_ context.Context, other Value) (int64, error) {
	otherBool, ok := other.(Boolean)

	if !ok {
		return CompareTypes(t, other), nil
	}

	if t == otherBool {
		return 0, nil
	}

	if !t && otherBool {
		return -1, nil
	}

	return +1, nil
}
