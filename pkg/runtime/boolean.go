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

func ParseBoolean(input any) (Boolean, error) {
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

func MustParseBoolean(input any) Boolean {
	res, err := ParseBoolean(input)

	if err != nil {
		panic(err)
	}

	return res
}

func ToBoolean(input Value) Boolean {
	if input == None {
		return False
	}

	switch val := input.(type) {
	case Boolean:
		return val
	case String:
		return val != ""
	case Int:
		return val != 0
	case Float:
		return val != 0
	case Duration:
		return val != 0
	case DateTime:
		return Boolean(!val.IsZero())
	default:
		return True
	}
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

func (t Boolean) Unwrap() any {
	return bool(t)
}
