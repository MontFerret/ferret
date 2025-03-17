package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// PUT makes a PUT HTTP request.
// @param {hashMap} params - Request parameters.
// @param {String} params.url - Target url
// @param {Any} params.body - Request data
// @param {hashMap} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func PUT(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod(ctx, h.MethodPut, args)
}
