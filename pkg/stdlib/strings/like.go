package strings

import (
	"context"
	"unicode"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Like checks whether the pattern search is contained in the string text, using wildcard matching.
// @param text (String) - The string to search in.
// @param search (String) - A search pattern that can contain the wildcard characters.
// @param caseInsensitive (Boolean) - If set to true, the matching will be case-insensitive. The default is false.
// @return (Boolean) - Returns true if the pattern is contained in text, and false otherwise.
func Like(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.False, err
	}

	str := []rune(args[0].String())
	pattern := []rune(args[1].String())

	if len(pattern) == 0 {
		return values.NewBoolean(len(str) == 0), nil
	}

	lookup := make([][]bool, len(str)+1)

	for i := range lookup {
		lookup[i] = make([]bool, len(pattern)+1)
	}

	lookup[0][0] = true

	for j := 1; j < len(pattern)+1; j++ {
		if pattern[j-1] == '%' {
			lookup[0][j] = lookup[0][j-1]
		}
	}

	for i := 1; i < len(str)+1; i++ {
		for j := 1; j < len(pattern)+1; j++ {
			switch {
			case pattern[j-1] == '%':
				lookup[i][j] = lookup[i][j-1] || lookup[i-1][j]
			case pattern[j-1] == '_' || str[i-1] == pattern[j-1]:
				lookup[i][j] = lookup[i-1][j-1]
			case len(args) > 2:
				isEq := unicode.ToLower(str[i-1]) == unicode.ToLower(pattern[j-1])
				if args[2] == values.True && isEq {
					lookup[i][j] = lookup[i-1][j-1]
				}
			default:
				lookup[i][j] = false
			}
		}
	}

	matched := lookup[len(str)][len(pattern)]

	return values.NewBoolean(matched), nil
}
