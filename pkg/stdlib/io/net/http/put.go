package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// PUT makes a PUT HTTP request.
// @param {Map} params - Request parameters.
// @param {String} params.url - Target url
// @param {Any} params.body - Request data
// @param {Map} [params.headers] - HTTP headers
// @return {binary} - Response in binary format
func PUT(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodPut, arg)
}
