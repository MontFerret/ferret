package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WaitElement waits for element to appear in the DOM.
// Stops the execution until it finds an element or operation times out.
// @param doc (HTMLDocument) - Driver HTMLDocument.
// @param selector (String) - Target element's selector.
// @param timeout (Int, optional) - Optional timeout. Default 5000 ms.
func WaitElement(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	selector := args[1].String()
	timeout := values.NewInt(defaultTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	ctx, fn := context.WithTimeout(ctx, time.Duration(timeout))
	defer fn()

	return values.None, doc.WaitForSelector(ctx, values.NewString(selector))
}
