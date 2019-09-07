package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// XPATH evaluates the XPath expression.
// @param source (HTMLPage | HTMLDocument | HTMLElement) - Target HTML object.
// @param expression (String) - XPath expression.
// @returns (Value) - Returns result of a given XPath expression.
func XPath(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	element, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	expr := values.ToString(args[1])

	return element.XPath(ctx, expr)
}
