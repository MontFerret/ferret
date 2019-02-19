package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Navigate navigates a document to a new resource.
// The operation blocks the execution until the page gets loaded.
// Which means there is no need in WAIT_NAVIGATION function.
// @param doc (Document) - Target document.
// @param url (String) - Target url to navigate.
// @param timeout (Int, optional) - Optional timeout. Default is 5000.
func Navigate(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(defaultTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	return values.None, doc.Navigate(args[1].(values.String), timeout)
}
