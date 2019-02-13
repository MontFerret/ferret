package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ScrollTop scrolls the document's window to its bottom.
// @param doc (HTMLDocument) - Target document.
func ScrollBottom(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.HTMLDocument)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(values.DHTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, doc.ScrollBottom()
}
