package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// FindFirst returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param src (String) - The source string.
// @param search (String) - The string to seek.
// @param start (Int, optional) - Limit the search to a subset of the text, beginning at start.
// @param end (Int, optional) - Limit the search to a subset of the text, ending at end
// @returns (Int) - The character position of the match.
// If search is not contained in text, -1 is returned. If search is empty, start is returned.
func FindFirst(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.NewInt(-1), err
	}

	argsCount := len(args)

	text := args[0].String()
	runes := []rune(text)
	search := args[1].String()
	start := values.NewInt(0)
	end := values.NewInt(int64(len(text)))

	if argsCount == 3 {
		arg3 := args[2]

		if arg3.Type() == core.IntType {
			start = arg3.(values.Int)
		}
	}

	if argsCount == 4 {
		arg4 := args[3]

		if arg4.Type() == core.IntType {
			end = arg4.(values.Int)
		}
	}

	found := int64(strings.Index(string(runes[start:end]), search))

	if found > -1 {
		return values.NewInt(found + int64(start)), nil
	}

	return values.NewInt(found), nil
}

// FindLast returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param src (String) - The source string.
// @param search (String) - The string to seek.
// @param start (Int, optional) - Limit the search to a subset of the text, beginning at start.
// @param end (Int, optional) - Limit the search to a subset of the text, ending at end
// @returns (Int) - The character position of the match.
// If search is not contained in text, -1 is returned. If search is empty, start is returned.
func FindLast(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.NewInt(-1), err
	}

	argsCount := len(args)

	text := args[0].String()
	runes := []rune(text)
	search := args[1].String()
	start := values.NewInt(0)
	end := values.NewInt(int64(len(text)))

	if argsCount == 3 {
		arg3 := args[2]

		if arg3.Type() == core.IntType {
			start = arg3.(values.Int)
		}
	}

	if argsCount == 4 {
		arg4 := args[3]

		if arg4.Type() == core.IntType {
			end = arg4.(values.Int)
		}
	}

	found := int64(strings.LastIndex(string(runes[start:end]), search))

	if found > -1 {
		return values.NewInt(found + int64(start)), nil
	}

	return values.NewInt(found), nil
}
