package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// X returns QuerySelector of XPath kind.
// @param {String} expression - XPath expression.
// @return {Any} - Returns QuerySelector of XPath kind.
func XPathSelector(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	return drivers.NewXPathSelector(values.ToString(args[0])), nil
}
