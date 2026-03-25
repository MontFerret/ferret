package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DELETE makes a DELETE request.
// @param {Map} params - Request parameters.
// @param {String} params.url - Target url
// @param {binary} params.body - Request data
// @param {Map} [params.headers] - HTTP headers
// @return {binary} - Response in binary format
func DELETE(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodDelete, arg)
}
