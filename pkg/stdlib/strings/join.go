package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// join concatenates a list of strings with one separator between adjacent elements, preserving empty strings.
// @param values {String[]} The strings to join in iteration order. Nested lists and non-strings are invalid.
// @param separator {String} The literal separator.
// @return {String} The joined string, or an empty string for an empty list.
func Join(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	values, separator, err := runtime.CastArgs2[runtime.List, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	var out strings.Builder
	first := true
	err = values.ForEach(ctx, func(c context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		text, ok := value.(runtime.String)
		if !ok {
			return false, runtime.ArgError(runtime.Errorf(runtime.ErrInvalidArgument, "element %d must be a String", index), 0)
		}

		if !first {
			out.WriteString(string(separator))
		}

		first = false
		out.WriteString(string(text))

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return runtime.String(out.String()), nil
}
