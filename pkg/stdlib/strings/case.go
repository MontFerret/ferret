package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// LOWER converts strings to their lower-case counterparts. All other characters are returned unchanged.
// @param {String} str - The source string.
// @return {String} - THis string in lower case.
func Lower(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text := strings.ToLower(arg.String())

	return runtime.NewString(text), nil
}

// UPPER converts strings to their upper-case counterparts. All other characters are returned unchanged.
// @param {String} str - The source string.
// @return {String} - THis string in upper case.
func Upper(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text := strings.ToUpper(arg.String())

	return runtime.NewString(text), nil
}
