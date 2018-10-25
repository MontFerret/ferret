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

// WaitNavigation waits for document to navigate to a new url.
// Stops the execution until the navigation ends or operation times out.
<<<<<<< HEAD
// @param doc (HTMLDocument) - Driver HTMLDocument.
// @param timeout (Int, optional) - Optional timeout. Default 5000 ms.
func WaitNavigation(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
// @param doc (HTMLDocument) - Dynamic HTMLDocument.
// @param timeout (Int, optional) - Optional timeout. Default 5000 ms.
func WaitNavigation(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(defaultTimeout)

	if len(args) > 1 {
		err = core.ValidateType(args[1], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[1].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.None, doc.WaitForNavigation(ctx)
}
