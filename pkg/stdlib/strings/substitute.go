package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Substitute replaces search values in the string value.
// @params text (String) - The string to modify
// @params search (String) - The string representing a search pattern
// @params replace (String) - The string representing a replace value
// @param limit (Int) - The cap the number of replacements to this value.
// @return (String) - Returns a string with replace substring.
func Substitute(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	search := args[1].String()
	replace := ""
	limit := -1

	if len(args) > 2 {
		replace = args[2].String()
	}

	if len(args) > 3 {
		if args[3].Type() == types.Int {
			limit = int(args[3].(values.Int))
		}
	}

	out := strings.Replace(text, search, replace, limit)

	return values.NewString(out), nil
}
