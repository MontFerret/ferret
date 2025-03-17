package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DELETE makes a DELETE request.
// @param {hashMap} params - Request parameters.
// @param {String} params.url - Target url
// @param {Binary} params.body - Request data
// @param {hashMap} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func DELETE(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod(ctx, h.MethodDelete, args)
}
