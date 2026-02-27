package strings

import (
	"context"
	"html"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// UNESCAPE_HTML unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&#xE1;".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
// @param {String} uri - Uri to escape.
// @return {String} - Escaped string.
func UnescapeHTML(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.NewString(html.UnescapeString(arg.String())), nil
}
