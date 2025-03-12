package core

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/wI2L/jettison"
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
	if IsNil(input) {
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

	return EmptyString, Error(ErrInvalidType, "expected '"+TypeString+"'")
}

func MustParseString(input interface{}) String {
	res, err := ParseString(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (s String) Type() string {
	return TypeString
}

func (s String) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(string(s), jettison.NoHTMLEscaping())
}

func (s String) String() string {
	return string(s)
}

func (s String) Unwrap() interface{} {
	return string(s)
}

func (s String) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeString))
	h.Write([]byte(s))

	return h.Sum64()
}

func (s String) Copy() Value {
	return s
}

func (s String) Compare(_ context.Context, other Value) (int64, error) {
	otherString, ok := other.(String)

	if !ok {
		return CompareTypes(s, other), nil
	}

	return int64(strings.Compare(string(s), otherString.Unwrap().(string))), nil
}

func (s String) Length(_ context.Context) (int, error) {
	return len([]rune(s)), nil
}

func (s String) Contains(other String) Boolean {
	return s.IndexOf(other) > -1
}

func (s String) IndexOf(other String) Int {
	return Int(strings.Index(string(s), string(other)))
}

func (s String) Concat(other Value) String {
	return String(string(s) + other.String())
}

func (s String) At(index Int) String {
	return String([]rune(s)[index])
}
