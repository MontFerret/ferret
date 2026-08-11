package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DELETE makes a DELETE request.
// @param params {Map} Request parameters containing url, body, and optional headers fields.
// @return {Binary} Response in binary format
func DELETE(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodDelete, arg)
}
