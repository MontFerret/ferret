package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
)

/*
 * Navigates a document to a new resource.
 * The operation blocks the execution until the page gets loaded.
 * Which means there is no need in WAIT_NAVIGATION function.
 * @param doc (Document) - Target document.
 * @param url (String) - Target url to navigate.
 */
func Navigate(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

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

	return values.None, doc.Navigate(args[1].(values.String))
}
