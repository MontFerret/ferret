package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// NavigateBack navigates a document back within its navigation history.
// The operation blocks the execution until the page gets loaded.
// If the history is empty, the function returns FALSE.
// @param doc (Document) - Target document.
// @param entry (Int, optional) - Optional value indicating how many pages to skip. Default 1.
// @param timeout (Int, optional) - Optional timeout. Default is 5000.
// @returns (Boolean) - Returns TRUE if history exists and the operation succeeded, otherwise FALSE.
func NavigateBack(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 3)

	if err != nil {
		return values.False, err
	}

	doc, err := toDocument(args[0])

	if err != nil {
		return values.None, err
	}

	skip := values.NewInt(1)
	timeout := values.NewInt(defaultTimeout)

	if len(args) > 1 {
		err = core.ValidateType(args[1], types.Int)

		if err != nil {
			return values.None, err
		}

		skip = args[1].(values.Int)
	}

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	ctx, fn := context.WithTimeout(ctx, time.Duration(timeout))
	defer fn()

	return doc.NavigateBack(ctx, skip)
}
