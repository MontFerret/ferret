package html

import (
	"context"

<<<<<<< HEAD
	"github.com/MontFerret/ferret/pkg/drivers"
=======
	"github.com/MontFerret/ferret/pkg/html/dynamic"
>>>>>>> 9f24172... rewrite comments
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ClickAll dispatches click event on all matched element
// @param source (Document) - Document.
// @param selector (String) - Selector.
// @returns (Boolean) - Returns true if matched at least one element.
<<<<<<< HEAD
func ClickAll(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func ClickAll(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
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

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	return doc.ClickBySelectorAll(ctx, values.NewString(selector))
}
