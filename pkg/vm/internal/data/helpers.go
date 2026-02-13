package data

import (
	"bytes"
	"context"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func ToRegexp(input runtime.Value) (*Regexp, error) {
	switch r := input.(type) {
	case *Regexp:
		return r, nil
	case runtime.String:
		return NewRegexp(r)
	default:
		return nil, runtime.TypeErrorOf(input, runtime.TypeString, runtime.TypeRegexp)
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

// Stringify converts a Value to a String. If the input is an Iterable, it concatenates
func Stringify(ctx context.Context, input runtime.Value) (string, error) {
	if input == nil {
		return "", nil
	}

	switch val := input.(type) {
	case runtime.Iterable:
		var b bytes.Buffer

		defer b.Reset()

		err := runtime.ForEach(ctx, val, func(ctx context.Context, value, key runtime.Value) (runtime.Boolean, error) {
			keyStr, err := Stringify(ctx, key)

			if err != nil {
				return runtime.False, err
			}

			valStr, err := Stringify(ctx, value)

			if err != nil {
				return runtime.False, err
			}

			b.WriteString(keyStr)
			b.Write([]byte(":"))
			b.WriteString(valStr)
			b.Write([]byte(";"))

			return runtime.True, nil
		})

		if err != nil {
			return "", err
		}

		return b.String(), nil
	default:
		return val.String(), nil
	}
}
