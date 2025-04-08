package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"strings"
)

func ToNumberOnly(ctx context.Context, input core.Value) core.Value {
	switch value := input.(type) {
	case core.Int, core.Float:
		return value
	case core.String:
		if strings.Contains(value.String(), ".") {
			if val, err := core.ToFloat(ctx, value); err == nil {
				return val
			}

			return core.ZeroFloat
		}

		if val, err := core.ToInt(ctx, value); err == nil {
			return val
		}

		return core.ZeroFloat
	case core.Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return core.ZeroInt
		}

		i := core.ZeroInt
		f := core.ZeroFloat

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				continue
			}

			out := ToNumberOnly(ctx, val)

			switch item := out.(type) {
			case core.Int:
				i += item
			case core.Float:
				f += item
			}
		}

		if f == 0 {
			return i
		}

		return core.Float(i) + f
	default:
		if val, err := core.ToFloat(ctx, value); err == nil {
			return val
		}

		return core.ZeroInt
	}
}

func ToRegexp(input core.Value) (*Regexp, error) {
	switch r := input.(type) {
	case *Regexp:
		return r, nil
	case core.String:
		return NewRegexp(r)
	default:
		return nil, core.TypeError(input, core.TypeString, "regexp")
	}
}
