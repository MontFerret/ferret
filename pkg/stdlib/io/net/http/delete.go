package http

import (
	"context"
	h "net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// DELETE makes a HTTP DELETE request.
// @param params (Object) - Request parameters.
//    * url (String) - Target url
//    * body (Binary) - POST data
//    * headers (Object) optional - HTTP headers
// @return (Binary) - Response in binary format
func DELETE(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod(ctx, h.MethodDelete, args)
}
