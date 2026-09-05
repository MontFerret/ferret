package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// contains reports whether text contains search.
// @param text {String} The source string.
// @param search {String} The literal string to match.
// @return {Boolean} Whether the literal string matches.
func Contains(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, search, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return text.Contains(search), nil
}

// starts_with reports whether text begins with prefix.
// @param text {String} The source string.
// @param prefix {String} The literal string to match.
// @return {Boolean} Whether the literal string matches.
func StartsWith(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, search, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Boolean(strings.HasPrefix(string(text), string(search))), nil
}

// ends_with reports whether text ends with suffix.
// @param text {String} The source string.
// @param suffix {String} The literal string to match.
// @return {Boolean} Whether the literal string matches.
func EndsWith(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, search, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Boolean(strings.HasSuffix(string(text), string(search))), nil
}
