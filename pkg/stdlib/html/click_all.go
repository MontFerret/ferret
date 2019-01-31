package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ClickAll dispatches click event on all matched element
// @param source (Document) - Document.
// @param selector (String) - Selector.
// @returns (Boolean) - Returns true if matched at least one element.
func ClickAll(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.False, err
	}

	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, drivers.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(drivers.DHTMLDocument)

	if !ok {
		return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return doc.ClickBySelectorAll(values.NewString(selector))
}
