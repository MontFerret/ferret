package types

import (
	"context"

<<<<<<< HEAD
	"github.com/MontFerret/ferret/pkg/drivers"
=======
>>>>>>> 9f24172... rewrite comments
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// IsHTMLDocument checks whether value is a HTMLDocument value.
// @param value (Value) - Input value of arbitrary type.
// @returns (Boolean) - Returns true if value is HTMLDocument, otherwise false.
func IsHTMLDocument(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	return isTypeof(args[0], drivers.HTMLDocumentType), nil
}
