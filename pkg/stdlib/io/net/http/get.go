package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// GET makes a HTTP GET request.
// @params url or  (String) - target url or parameters.
//    * url (String) - Target url
//    * headers (Object) optional - HTTP headers
func GET(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	arg := args[0]

	if err := core.ValidateType(arg, types.String, types.Object); err != nil {
		return values.None, err
	}

	if arg.Type() == types.String {
		return makeRequest(ctx, Params{
			Method:  "GET",
			URL:     values.ToString(arg),
			Headers: nil,
			Body:    nil,
		})
	}

	return execMethod(ctx, h.MethodGet, args)
}
