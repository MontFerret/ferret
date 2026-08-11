package runtime

import (
	"context"
	"hash/fnv"
	"regexp"
	"strings"
)

// Regexp is a compiled regular expression used by Ferret matching operations.
type Regexp regexp.Regexp

// NewRegexp compiles pattern into a Ferret regular expression value.
func NewRegexp(pattern string) (*Regexp, error) {
	r, err := regexp.Compile(pattern)

	if err != nil {
		return nil, err
	}

	return (*Regexp)(r), nil
}

// String returns the source pattern.
func (r *Regexp) String() string {
	return (*regexp.Regexp)(r).String()
}

// Hash returns a deterministic hash of the source pattern.
func (r *Regexp) Hash() uint64 {
	h := fnv.New64a()

	// Preserve the historical namespace after moving Regexp from the VM package.
	h.Write([]byte("vm.regexp"))
	h.Write([]byte(":"))
	h.Write([]byte(r.String()))

	return h.Sum64()
}

// Copy returns an independently compiled expression with the same pattern.
func (r *Regexp) Copy() Value {
	copied, err := NewRegexp(r.String())
	if err != nil {
		return r
	}

	return copied
}

// Type returns the stable built-in identity used by validation and diagnostics.
func (r *Regexp) Type() Type {
	return TypeRegexp
}

// Equal reports whether other has the same source pattern.
func (r *Regexp) Equal(_ context.Context, other Value) (bool, error) {
	otherRegexp, ok := other.(*Regexp)
	if !ok {
		return false, nil
	}

	return r.String() == otherRegexp.String(), nil
}

// Compare orders regular expressions lexicographically by source pattern.
func (r *Regexp) Compare(_ context.Context, other Value) (Ordering, error) {
	otherRegexp, ok := other.(*Regexp)
	if !ok {
		return Equal, incompatibleComparisonError(r, other)
	}

	return Ordering(strings.Compare(r.String(), otherRegexp.String())), nil
}

// Match reports whether the string form of value matches the expression.
func (r *Regexp) Match(value Value) Boolean {
	return Boolean((*regexp.Regexp)(r).MatchString(value.String()))
}

// MatchString reports whether str matches the expression.
func (r *Regexp) MatchString(str string) Boolean {
	return Boolean((*regexp.Regexp)(r).MatchString(str))
}
