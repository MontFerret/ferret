package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// put makes a PUT HTTP request.
// @param params {Map} Request parameters containing url, body, and optional headers fields.
// @return {Binary} Response in binary format
func PUT(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	return execMethod(ctx, h.MethodPut, arg)
}
