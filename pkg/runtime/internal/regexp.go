package internal

import (
	"context"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Regexp regexp.Regexp

func NewRegexp(pattern core.String) (*Regexp, error) {
	r, err := regexp.Compile(string(pattern))

	if err != nil {
		return nil, err
	}

	return (*Regexp)(r), nil
}

func (r *Regexp) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(r.String(), jettison.NoHTMLEscaping())
}

func (r *Regexp) String() string {
	return (*regexp.Regexp)(r).String()
}

func (r *Regexp) Unwrap() interface{} {
	return (*regexp.Regexp)(r)
}

func (r *Regexp) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("regexp"))
	h.Write([]byte(":"))
	h.Write([]byte(r.String()))

	return h.Sum64()
}

func (r *Regexp) Copy() core.Value {
	copied, err := NewRegexp(core.String(r.String()))

	// it should never happen
	if err != nil {
		panic(err)
	}

	return copied
}

func (r *Regexp) Compare(_ context.Context, other core.Value) (int64, error) {
	otherRegexp, ok := other.(*Regexp)

	if !ok {
		return core.CompareTypes(r, other), nil
	}

	return int64(strings.Compare(r.String(), otherRegexp.String())), nil
}

func (r *Regexp) Match(value core.Value) core.Boolean {
	return core.Boolean((*regexp.Regexp)(r).MatchString(value.String()))
}

func (r *Regexp) MatchString(str core.String) core.Boolean {
	return core.Boolean((*regexp.Regexp)(r).MatchString(string(str)))
}
