package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// PUT makes a PUT HTTP request.
// @param params (Object) - Request parameters.
//    * url (String) - Target url.
//    * body (Binary) - POST data.
//    * headers (Object) optional - HTTP headers.
// @return (Binary) - Response in binary format
func PUT(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod(ctx, h.MethodPut, args)
}
