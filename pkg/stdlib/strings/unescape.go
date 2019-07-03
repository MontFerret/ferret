package strings

import (
	"context"
	"html"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// UnescapeHTML unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&#xE1;".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
// @param (String) - Uri to escape.
// @returns String - Escaped string.
func UnescapeHTML(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return values.NewString(html.UnescapeString(args[0].String())), nil
}
