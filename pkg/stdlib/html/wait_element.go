package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WaitElement waits for element to appear in the DOM.
// Stops the execution until it finds an element or operation times out.
// @param doc (HTMLDocument) - Dynamic HTMLDocument.
// @param selector (String) - Target element's selector.
// @param timeout (Int, optional) - Optional timeout. Default 5000 ms.
func WaitElement(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	arg := args[0]
	selector := args[1].String()
	timeout := values.NewInt(defaultTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	err = core.ValidateType(arg, drivers.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg.(drivers.DHTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, doc.WaitForSelector(values.NewString(selector), timeout)
}
