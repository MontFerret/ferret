package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// IS_HTML_ELEMENT checks whether value is a HTMLElement value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is HTMLElement, otherwise false.
func IsHTMLElement(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], drivers.HTMLElementType), nil
}
