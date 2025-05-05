package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
	"strings"
	"time"
)

func ToNumberOnly(ctx context.Context, input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int, runtime.Float:
		return value
	case runtime.String:
		if strings.Contains(value.String(), ".") {
			if val, err := runtime.ToFloat(ctx, value); err == nil {
				return val
			}

			return runtime.ZeroFloat
		}

		if val, err := runtime.ToInt(ctx, value); err == nil {
			return val
		}

		return runtime.ZeroFloat
	case runtime.Iterable:
		iterator, err := value.Iterate(ctx)

		if err != nil {
			return runtime.ZeroInt
		}

		i := runtime.ZeroInt
		f := runtime.ZeroFloat

		for hasNext, err := iterator.HasNext(ctx); hasNext && err == nil; {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				continue
			}

			out := ToNumberOnly(ctx, val)

			switch item := out.(type) {
			case runtime.Int:
				i += item
			case runtime.Float:
				f += item
			}
		}

		if f == 0 {
			return i
		}

		return runtime.Float(i) + f
	default:
		if val, err := runtime.ToFloat(ctx, value); err == nil {
			return val
		}

		return runtime.ZeroInt
	}
}

func ToRegexp(input runtime.Value) (*Regexp, error) {
	switch r := input.(type) {
	case *Regexp:
		return r, nil
	case runtime.String:
		return NewRegexp(r)
	default:
		return nil, runtime.TypeError(input, runtime.TypeString, "regexp")
	}
}

func Sleep(ctx context.Context, duration runtime.Int) error {
	timer := time.NewTimer(time.Millisecond * time.Duration(duration))

	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	return nil
}
