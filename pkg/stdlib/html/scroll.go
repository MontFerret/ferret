package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ScrollTop Scrolls the document's window to its top.
// @param doc (Document) - Target document.
func ScrollTop(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, doc.ScrollTop()
}

// ScrollTop Scrolls the document's window to its bottom.
// @param doc (Document) - Target document.
func ScrollBottom(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HTMLDocument)

	if !ok {
		return values.None, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return values.None, doc.ScrollBottom()
}
