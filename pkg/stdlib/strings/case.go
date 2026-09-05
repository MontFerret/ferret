package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// lower converts strings to their lower-case counterparts. All other characters are returned unchanged.
// @param str {String} The source string.
// @return {String} The string in lower case.
func Lower(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	str, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	text := strings.ToLower(string(str))

	return runtime.NewString(text), nil
}

// upper converts strings to their upper-case counterparts. All other characters are returned unchanged.
// @param str {String} The source string.
// @return {String} The string in upper case.
func Upper(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	str, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	text := strings.ToUpper(string(str))

	return runtime.NewString(text), nil
}
