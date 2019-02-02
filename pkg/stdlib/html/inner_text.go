package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// InnerText returns inner text string of a given or matched by CSS selector element
// @param doc (HTMLDocument|HTMLNode) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @returns (String) - Inner text if an element found, otherwise empty string.
func InnerText(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLNodeType)

	if err != nil {
		return values.None, err
	}

	node := args[0].(drivers.HTMLNode)

	if len(args) == 1 {
		return node.InnerText(), nil
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return node.InnerTextBySelector(selector), nil
}
