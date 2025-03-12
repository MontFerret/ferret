package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// SUBSTITUTE replaces search values in the string value.
// @param {String} str - The string to modify
// @param {String} search - The string representing a search pattern
// @param {String} replace - The string representing a replace value
// @param {Int} limit - The cap the number of replacements to this value.
// @return {String} - Returns a string with replace substring.
func Substitute(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	search := args[1].String()
	replace := ""
	limit := -1

	if len(args) > 2 {
		replace = args[2].String()
	}

	if len(args) > 3 {
		arg3, ok := args[3].(core.Int)

		if ok {
			limit = int(arg3)
		}
	}

	out := strings.Replace(text, search, replace, limit)

	return core.NewString(out), nil
}
