package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Navigates a document to a new resource.
 * The operation blocks the execution until the page gets loaded.
 * Which means there is no need in WAIT_NAVIGATION function.
 * @param doc (Document) - Target document.
 * @param url (String) - Target url to navigate.
 * @param timeout (Int, optional) - Optional timeout. Default is 5000.
 */
func Navigate(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	timeout := values.NewInt(defaultTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], core.IntType)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	return values.None, doc.Navigate(args[1].(values.String), timeout)
}
