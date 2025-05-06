package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
)

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Boolean} [returnIndex=False] - Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} - A value indicating whether a specified substring occurs within a string.
func Contains(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.False, err
	}

	text := runtime.SafeCastString(args[0], runtime.EmptyString)
	search := runtime.SafeCastString(args[1], runtime.EmptyString)
	returnIndex := runtime.False

	if len(args) > 2 {
		returnIndex = runtime.SafeCastBoolean(args[2], runtime.False)
	}

	if returnIndex == runtime.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}
