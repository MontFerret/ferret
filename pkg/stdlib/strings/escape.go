package strings

import (
	"context"
	"html"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ESCAPE_HTML escapes special characters like "<" to become "&lt;". It
// escapes only five such characters: <, >, &, ' and ".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
// @param {String} uri - Uri to escape.
// @return {String} - Escaped string.
func EscapeHTML(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.NewString(html.EscapeString(arg.String())), nil
}
