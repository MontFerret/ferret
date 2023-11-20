package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Int} [start] - Limit the search to a subset of the text, beginning at start.
// @param {Int} [end] - Limit the search to a subset of the text, ending at end
// @return {Int} - The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
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
	end := values.NewInt(len(text))

	if argsCount == 3 {
		arg3, ok := args[2].(values.Int)

		if ok {
			start = arg3
		}
	}

	if argsCount == 4 {
		arg4, ok := args[3].(values.Int)

		if ok {
			end = arg4
		}
	}

	found := strings.Index(string(runes[start:end]), search)

	if found > -1 {
		return values.NewInt(found + int(start)), nil
	}

	return values.NewInt(found), nil
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param {String} src - The source string.
// @param {String} search - The string to seek.
// @param {Int} [start] - Limit the search to a subset of the text, beginning at start.
// @param {Int} [end] - Limit the search to a subset of the text, ending at end
// @return {Int} - The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
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
	end := values.NewInt(len(text))

	if argsCount == 3 {
		arg3, ok := args[2].(values.Int)

		if ok {
			start = arg3
		}
	}

	if argsCount == 4 {
		arg4, ok := args[3].(values.Int)

		if ok {
			end = arg4
		}
	}

	found := strings.LastIndex(string(runes[start:end]), search)

	if found > -1 {
		return values.NewInt(found + int(start)), nil
	}

	return values.NewInt(found), nil
}
