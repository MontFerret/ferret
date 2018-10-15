package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// WaitClassAll waits for a class to appear on all matched elements.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Parent document.
// @param selector (String) - String of CSS selector.
// @param class (String) - String of target CSS class.
// @param timeout (Int, optional) - Optional timeout.
func WaitClassAll(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)

	if err != nil {
		return values.None, err
	}

	// document only
	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	// selector
	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	// class
	err = core.ValidateType(args[2], core.StringType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	selector := args[1].(values.String)
	class := args[2].(values.String)
	timeout := values.NewInt(defaultTimeout)

	if len(args) == 4 {
		err = core.ValidateType(args[3], core.IntType)

		if err != nil {
			return values.None, err
		}

		timeout = args[3].(values.Int)
	}

	return values.None, doc.WaitForClassAll(selector, class, timeout)
}
