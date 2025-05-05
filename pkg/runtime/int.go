package runtime

import (
	"encoding/binary"
	"hash/fnv"
	"strconv"

	"github.com/wI2L/jettison"
)

type Int int64

const (
	ZeroInt = Int(0)
)

func NewInt(input int) Int {
	return Int(int64(input))
}

func NewInt64(input int64) Int {
	return Int(input)
}

func ParseInt(input interface{}) (Int, error) {
	if IsNil(input) {
		return ZeroInt, nil
	}

	switch val := input.(type) {
	case int:
		return Int(val), nil
	case int64:
		return Int(val), nil
	case int32:
		return Int(val), nil
	case int16:
		return Int(val), nil
	case int8:
		return Int(val), nil
	case string:
		i, err := strconv.Atoi(val)

		if err == nil {
			if i == 0 {
				return ZeroInt, nil
			}

			return Int(i), nil
		}

		return ZeroInt, err
	default:
		return ZeroInt, Error(ErrInvalidType, "expected 'int'")
	}
}

func MustParseInt(input interface{}) Int {
	res, err := ParseInt(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (i Int) Type() string {
	return TypeInt
}

func (i Int) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(int64(i), jettison.NoHTMLEscaping())
}

func (i Int) String() string {
	return strconv.Itoa(int(i))
}

func (i Int) Unwrap() interface{} {
	return int(i)
}

func (i Int) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeInt))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(i))
	h.Write(bytes)

	return h.Sum64()
}

func (i Int) Copy() Value {
	return i
}

func (i Int) Compare(other Value) int64 {
	switch otherVal := other.(type) {
	case Int:
		if i == otherVal {
			return 0
		}

		if i < otherVal {
			return -1
		}

		return +1

	case Float:
		f := Float(i)

		if f == otherVal {
			return 0
		}

		if f < otherVal {
			return -1
		}

		return +1

	default:
		return CompareTypes(i, other)
	}
}
