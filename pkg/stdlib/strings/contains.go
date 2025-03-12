package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Boolean} [returnIndex=False] - Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} - A value indicating whether a specified substring occurs within a string.
func Contains(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 3); err != nil {
		return core.False, err
	}

	text := core.SafeCastString(args[0], core.EmptyString)
	search := core.SafeCastString(args[1], core.EmptyString)
	returnIndex := core.False

	if len(args) > 2 {
		returnIndex = core.SafeCastBoolean(args[2], core.False)
	}

	if returnIndex == core.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}
