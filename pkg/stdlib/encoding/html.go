package encoding

import (
	"context"
	"html"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// html_escape escapes the five HTML-sensitive characters: ampersand, angle brackets, and quotes.
// @param text {String} The input string.
// @return {String} The converted string.
func HTMLEscape(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(html.EscapeString(string(text))), nil
}

// html_unescape decodes HTML character references.
// @param text {String} The input string.
// @return {String} The converted string.
func HTMLUnescape(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(html.UnescapeString(string(text))), nil
}
