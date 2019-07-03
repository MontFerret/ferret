package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ClickAll dispatches click event on all matched element
// @param source (Open) - Open.
// @param selector (String) - Selector.
// @returns (Boolean) - Returns true if matched at least one element.
func ClickAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.False, err
	}

	doc, err := drivers.ToDocument(args[0])

	if err != nil {
		return values.None, err
	}

	selector := args[1].String()

	return doc.ClickBySelectorAll(ctx, values.NewString(selector))
}
