package strings

import (
	"context"
	"regexp"
	"strings"

	"github.com/gobwas/glob"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var (
	deprecatedLikeSyntax = regexp.MustCompile("[%_]")
)

// LIKE checks whether the pattern search is contained in the string text, using wildcard matching.
// @param {String} str - The string to search in.
// @param {String} search - A search pattern that can contain the wildcard characters.
// @param {Boolean} caseInsensitive - If set to true, the matching will be case-insensitive. The default is false.
// @return {Boolean} - Returns true if the pattern is contained in text, and false otherwise.
func Like(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.False, err
	}

	str := args[0].String()
	pattern := args[1].String()

	if len(pattern) == 0 {
		return values.NewBoolean(len(str) == 0), nil
	}

	// TODO: Remove me in next releases
	replaced := deprecatedLikeSyntax.ReplaceAllFunc([]byte(pattern), func(b []byte) []byte {
		str := string(b)

		switch str {
		case "%":
			return []byte("*")
		case "_":
			return []byte("?")
		default:
			return b
		}
	})

	pattern = string(replaced)

	if len(args) > 2 {
		if values.ToBoolean(args[2]) {
			str = strings.ToLower(str)
			pattern = strings.ToLower(pattern)
		}
	}

	g, err := glob.Compile(pattern)

	if err != nil {
		return nil, err
	}

	return values.NewBoolean(g.Match(str)), nil
}
