package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// GET makes a GET request.
// @param {hashMap | String} urlOrParam - Target url or parameters.
// @param {String} [param.url] - Target url or parameters.
// @param {Map} [param.headers] - HTTP headers
// @return {Binary} - Response in binary format
func GET(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	switch v := arg.(type) {
	case runtime.String:
		return makeRequest(ctx, Params{
			Method:  "GET",
			URL:     v,
			Headers: nil,
			Body:    nil,
		})
	case runtime.Map:
		return execMethod(ctx, h.MethodGet, arg)
	default:
		return runtime.None, runtime.TypeError(runtime.TypeOf(arg), runtime.TypeString, runtime.TypeObject)
	}
}
