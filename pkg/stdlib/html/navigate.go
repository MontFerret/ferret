package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// NAVIGATE navigates a given page to a new resource.
// The operation blocks the execution until the page gets loaded.
// Which means there is no need in WAIT_NAVIGATION function.
// @param {HTMLPage} page - Target page.
// @param {String} url - Target url to navigate.
// @param {Int} [timeout=5000] - Navigation timeout.
func Navigate(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(drivers.DefaultWaitTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.True, page.Navigate(ctx, args[1].(values.String))
}
