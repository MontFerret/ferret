package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
)

/*
 * Dispatches click event on a given element
 * @param source (Document | Element) - Event source.
 * @param selector (String, optional) - Optional selector. Only used when a document instance is passed.
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

		el, ok := arg1.(*dynamic.HtmlElement)

		if !ok {
			return values.False, core.Error(core.ErrInvalidType, "expected dynamic element")
		}

		return el.Click()
	}

	// CLICK(doc, selector)
	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.ClickBySelector(values.NewString(selector))
}

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

	err = core.ValidateType(args[0], core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HtmlDocument)

	if !ok {
		return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return values.None, doc.Navigate(args[1].(values.String))
}

/*
 * Sends a value to an underlying input element.
 * @param source (Document | Element) - Event target.
 * @param valueOrSelector (String) - Selector or a value.
 * @param value (String) - Target value.
 * @returns (Boolean) - Returns true if an element was found.
 */
func Input(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	// TYPE(el, "foobar")
	if len(args) == 2 {
		arg1 := args[0]

		err := core.ValidateType(arg1, core.HtmlElementType)

		if err != nil {
			return values.False, err
		}

		el, ok := arg1.(*dynamic.HtmlElement)

		if !ok {
			return values.False, core.Error(core.ErrInvalidType, "expected dynamic element")
		}

		err = el.Input(args[1])

		if err != nil {
			return values.False, err
		}

		return values.True, nil
	}

	arg1 := args[0]

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.False, err
	}

	arg2 := args[1]

	err = core.ValidateType(arg2, core.StringType)

	if err != nil {
		return values.False, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.InputBySelector(arg2.(values.String), args[2])
}
