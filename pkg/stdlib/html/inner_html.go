package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns inner HTML string of a matched element
 * @param doc (Document|Element) - Parent document or element.
 * @param selector (String) - String of CSS selector.
 * @returns (String) - Inner HTML string if an element found, otherwise empty string.
 */
func InnerHTML(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType, core.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	node := args[0].(values.HTMLNode)
	selector := args[1].(values.String)

	return node.InnerHTMLBySelector(selector), nil
}
