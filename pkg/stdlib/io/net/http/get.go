package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// GET makes a GET request.
// @param {Object | String} urlOrParam - Target url or parameters.
// @param {String} [param.url] - Target url or parameters.
// @param {Object} [param.headers] - HTTP headers
// @return {Binary} - Response in binary format
func GET(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	arg := args[0]

	switch v := arg.(type) {
	case core.String:
		return makeRequest(ctx, Params{
			Method:  "GET",
			URL:     v,
			Headers: nil,
			Body:    nil,
		})
	case *internal.Object:
		return execMethod(ctx, h.MethodGet, args)
	default:
		return core.None, core.TypeError(arg, types.String, types.Object)
	}
}
