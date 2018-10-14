package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// WaitNavigation waits for document to navigate to a new url.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Dynamic HTMLDocument.
// @param timeout (Int, optional) - Optional timeout. Default 5000 ms.
func WaitNavigation(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	timeout := values.NewInt(defaultTimeout)

	if len(args) > 1 {
		err = core.ValidateType(args[1], core.IntType)

		if err != nil {
			return values.None, err
		}

		timeout = args[1].(values.Int)
	}

	return values.None, doc.WaitForNavigation(timeout)
}
