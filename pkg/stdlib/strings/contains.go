package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Boolean} [returnIndex=False] - Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} - A value indicating whether a specified substring occurs within a string.
func Contains(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.False, err
	}

	var text values.String
	var search values.String
	returnIndex := values.False

	arg1 := args[0]
	arg2 := args[1]

	if arg1.Type() == types.String {
		text = arg1.(values.String)
	} else {
		text = values.NewString(arg1.String())
	}

	if arg2.Type() == types.String {
		search = arg2.(values.String)
	} else {
		search = values.NewString(arg2.String())
	}

	if len(args) > 2 {
		arg3 := args[2]

		if arg3.Type() == types.Boolean {
			returnIndex = arg3.(values.Boolean)
		}
	}

	if returnIndex == values.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}
