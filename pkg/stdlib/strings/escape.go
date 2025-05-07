package strings

import (
	"context"
	"html"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ESCAPE_HTML escapes special characters like "<" to become "&lt;". It
// escapes only five such characters: <, >, &, ' and ".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
// @param {String} uri - Uri to escape.
// @return {String} - Escaped string.
func EscapeHTML(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(html.EscapeString(args[0].String())), nil
}
