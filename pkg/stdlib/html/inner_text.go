package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns inner text of a matched element
 * @param doc (Document) - Parent document.
 * @param selector (String) - String of CSS selector.
 * @returns (String) - Inner text if an element found, otherwise NONE.
 */
func InnerText(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], core.HTMLDocumentType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	doc := args[0].(values.HTMLDocument)
	selector := args[1].(values.String)

	return doc.InnerTextBySelector(selector), nil
}
