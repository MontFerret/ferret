package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
)

/*
 * Returns inner html of a matched element
 * @param doc (Document) - Document
 * @param selector (String) - Selector
 * @returns str (String) - String value of inner html.
 */
func InnerHTML(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.EmptyString, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.InnerHTMLBySelector(values.NewString(selector))
}

/*
 * Returns inner html of all matched elements.
 * @param doc (Document) - Document
 * @param selector (String) - Selector
 * @returns array (Array) - Array of string values.
 */
func InnerHTMLAll(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.EmptyString, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.InnerHTMLBySelectorAll(values.NewString(selector))
}

/*
 * Returns inner text of a matched element
 * @param doc (Document) - Document
 * @param selector (String) - Selector
 * @returns str (String) - String value of inner text.
 */
func InnerText(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.EmptyString, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.InnerHTMLBySelector(values.NewString(selector))
}

/*
 * Returns inner text of all matched elements.
 * @param doc (Document) - Document
 * @param selector (String) - Selector
 * @returns array (Array) - Array of string values.
 */
func InnerTextAll(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	arg1 := args[0]
	selector := args[1].String()

	err = core.ValidateType(arg1, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg1.(*dynamic.HtmlDocument)

	if !ok {
		return values.EmptyString, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return doc.InnerHTMLBySelectorAll(values.NewString(selector))
}
