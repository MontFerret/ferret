package values

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type String string

const (
	EmptyString = String("")
	SpaceString = String(" ")
)

func NewString(input string) String {
	if input == "" {
		return EmptyString
	}

	return String(input)
}

func NewStringFromRunes(input []rune) String {
	if len(input) == 0 {
		return EmptyString
	}

	return String(input)
}

func ParseString(input interface{}) (String, error) {
	if core.IsNil(input) {
		return EmptyString, nil
	}

	str, ok := input.(string)

	if ok {
		if str != "" {
			return String(str), nil
		}

		return EmptyString, nil
	}

	stringer, ok := input.(fmt.Stringer)

	if ok {
		return String(stringer.String()), nil
	}

	return EmptyString, core.Error(core.ErrInvalidType, "expected 'string'")
}

func MustParseString(input interface{}) String {
	res, err := ParseString(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t String) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(string(t), jettison.NoHTMLEscaping())
}

func (t String) Type() core.Type {
	return types.String
}

func (t String) String() string {
	return string(t)
}

func (t String) Compare(other core.Value) int64 {
	if other.Type() == types.String {
		return int64(strings.Compare(string(t), other.Unwrap().(string)))
	}

	return types.Compare(types.String, other.Type())
}

func (t String) Unwrap() interface{} {
	return string(t)
}

func (t String) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(t))

	return h.Sum64()
}

func (t String) Copy() core.Value {
	return t
}

func (t String) Length() Int {
	return Int(len([]rune(string(t))))
}

func (t String) Contains(other String) Boolean {
	return t.IndexOf(other) > -1
}

func (t String) IndexOf(other String) Int {
	return Int(strings.Index(string(t), string(other)))
}

func (t String) Concat(other core.Value) String {
	return String(string(t) + other.String())
}

func (t String) At(index Int) String {
	return String([]rune(t)[index])
}
