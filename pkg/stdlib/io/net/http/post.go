package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// POST makes a POST request.
// @param {Map} params - Request parameters.
// @param {String} params.url - Target url
// @param {Any} params.body - Request data
// @param {Map} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func POST(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodPost, arg)
}
