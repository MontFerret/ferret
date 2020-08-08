package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// IS_HTML_DOCUMENT checks whether value is a HTMLDocument value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is HTMLDocument, otherwise false.
func IsHTMLDocument(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], drivers.HTMLDocumentType), nil
}
