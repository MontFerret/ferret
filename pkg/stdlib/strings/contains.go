package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Boolean} [returnIndex=False] - Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} - A value indicating whether a specified substring occurs within a string.
func Contains(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 3); err != nil {
		return values.False, err
	}

	text := values.SafeCastString(args[0], values.EmptyString)
	search := values.SafeCastString(args[1], values.EmptyString)
	returnIndex := values.False

	if len(args) > 2 {
		returnIndex = values.SafeCastBoolean(args[2], values.False)
	}

	if returnIndex == values.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}
