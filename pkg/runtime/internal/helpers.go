package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"go/types"
)

func ToRegexp(input core.Value) (*Regexp, error) {
	switch r := input.(type) {
	case *Regexp:
		return r, nil
	case core.String:
		return NewRegexp(r)
	default:
		return nil, TypeError(input, types.String, types.Regexp)
	}
}
