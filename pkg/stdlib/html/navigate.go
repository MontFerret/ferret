package html

import (
	"context"
<<<<<<< HEAD
=======

	"github.com/MontFerret/ferret/pkg/html/dynamic"
>>>>>>> 9f24172... rewrite comments
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
<<<<<<< HEAD
func Navigate(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func Navigate(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

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

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.None, doc.Navigate(ctx, args[1].(values.String))
}
