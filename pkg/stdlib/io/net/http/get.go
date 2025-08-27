package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// GET makes a GET request.
// @param {hashMap | String} urlOrParam - Target url or parameters.
// @param {String} [param.url] - Target url or parameters.
// @param {hashMap} [param.headers] - HTTP headers
// @return {Binary} - Response in binary format
func GET(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arg := args[0]

	switch v := arg.(type) {
	case runtime.String:
		return makeRequest(ctx, Params{
			Method:  "GET",
			URL:     v,
			Headers: nil,
			Body:    nil,
		})
	case *runtime.Object:
		return execMethod(ctx, h.MethodGet, args)
	default:
		return runtime.None, runtime.TypeError(runtime.Reflect(arg), runtime.TypeString, runtime.TypeObject)
	}
}
