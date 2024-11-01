package values

import (
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Regexp regexp.Regexp

func NewRegexp(pattern String) (*Regexp, error) {
	r, err := regexp.Compile(string(pattern))

	if err != nil {
		return nil, err
	}

	return (*Regexp)(r), nil
}

func (r *Regexp) Type() core.Type {
	return types.Regexp
}

func (r *Regexp) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(r.String(), jettison.NoHTMLEscaping())
}

func (r *Regexp) String() string {
	return (*regexp.Regexp)(r).String()
}

func (r *Regexp) Compare(other core.Value) int64 {
	otherRegexp, ok := other.(*Regexp)

	if !ok {
		return types.Compare(types.Regexp, core.Reflect(other))
	}

	return int64(strings.Compare(r.String(), otherRegexp.String()))
}

func (r *Regexp) Unwrap() interface{} {
	return (*regexp.Regexp)(r)
}

func (r *Regexp) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Regexp.String()))
	h.Write([]byte(":"))
	h.Write([]byte(r.String()))

	return h.Sum64()
}

func (r *Regexp) Copy() core.Value {
	copied, err := NewRegexp(String(r.String()))

	// it should never happen
	if err != nil {
		panic(err)
	}

	return copied
}

func (r *Regexp) Match(value core.Value) Boolean {
	return Boolean((*regexp.Regexp)(r).MatchString(value.String()))
}

func (r *Regexp) MatchString(str String) Boolean {
	return Boolean((*regexp.Regexp)(r).MatchString(string(str)))
}
