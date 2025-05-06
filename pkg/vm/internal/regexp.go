package internal

import (
	"context"
	"hash/fnv"
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/wI2L/jettison"
)

type Regexp regexp.Regexp

func NewRegexp(pattern runtime.String) (*Regexp, error) {
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

func (r *Regexp) Copy() runtime.Value {
	copied, err := NewRegexp(runtime.String(r.String()))

	// it should never happen
	if err != nil {
		panic(err)
	}

	return copied
}

func (r *Regexp) Compare(_ context.Context, other runtime.Value) (int64, error) {
	otherRegexp, ok := other.(*Regexp)

	if !ok {
		return runtime.CompareTypes(r, other), nil
	}

	return int64(strings.Compare(r.String(), otherRegexp.String())), nil
}

func (r *Regexp) Match(value runtime.Value) runtime.Boolean {
	return runtime.Boolean((*regexp.Regexp)(r).MatchString(value.String()))
}

func (r *Regexp) MatchString(str runtime.String) runtime.Boolean {
	return runtime.Boolean((*regexp.Regexp)(r).MatchString(string(str)))
}
