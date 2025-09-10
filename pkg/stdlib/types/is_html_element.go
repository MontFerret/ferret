package types

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// IS_HTML_ELEMENT checks whether value is a HTMLElement value.
// @param {Any} value - Input value of arbitrary type.
// @return {Boolean} - Returns true if value is HTMLElement, otherwise false.
func IsHTMLElement(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	//err := runtime.ValidateArgs(args, 1, 1)
	//
	//if err != nil {
	//	return runtime.None, err
	//}
	//
	//return isTypeof(args[0], drivers.HTMLElementType), nil

	panic("implement me")
}
