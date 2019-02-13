package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// IsHTMLElement checks whether value is a HTMLElement value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) - Returns true if value is HTMLElement, otherwise false.
func IsHTMLElement(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], types.HTMLElement), nil
}
