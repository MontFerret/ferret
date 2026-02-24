package runtime

import (
	"hash/fnv"
	"strings"
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

	h.Write([]byte(TypeBoolean.Name()))
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

func (t Boolean) Type() Type {
	return TypeBoolean
}

func (t Boolean) String() string {
	if t {
		return "true"
	}

	return "false"
}

func (t Boolean) Hash() uint64 {
	if t {
		return hashTrue
	}

	return hashFalse
}

func (t Boolean) Copy() Value {
	return t
}

func (t Boolean) Compare(other Value) int {
	otherBool, ok := other.(Boolean)

	if !ok {
		return CompareTypes(t, other)
	}

	if t == otherBool {
		return 0
	}

	if !t && otherBool {
		return -1
	}

	return +1
}

func (t Boolean) Unwrap() any {
	return bool(t)
}
