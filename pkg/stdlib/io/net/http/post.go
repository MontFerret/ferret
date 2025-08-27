package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// POST makes a POST request.
// @param {hashMap} params - Request parameters.
// @param {String} params.url - Target url
// @param {Any} params.body - Request data
// @param {hashMap} [params.headers] - HTTP headers
// @return {Binary} - Response in binary format
func POST(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodPost, args)
}
