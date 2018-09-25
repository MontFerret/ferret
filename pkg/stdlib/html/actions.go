package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser"
)

/*
 *
 */
func Click(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.False, err
	}

	// CLICK(el)
	if len(args) == 1 {
		arg1 := args[0]

		err := core.ValidateType(arg1, core.HtmlElementType)

		if err != nil {
			return values.False, err
		}

		el, ok := arg1.(*browser.HtmlElement)

		if !ok {
			return values.False, core.Error(core.ErrInvalidType, "expected dynamic element")
		}

		return el.Click()
	} else {
		// CLICK(doc, selector)
		arg1 := args[0]
		selector := args[1].String()

		err = core.ValidateType(arg1, core.HtmlDocumentType)

		if err != nil {
			return values.None, err
		}

		doc, ok := arg1.(*browser.HtmlDocument)

		if !ok {
			return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
		}

		return doc.ClickBySelector(values.NewString(selector))
	}
}
