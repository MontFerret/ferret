package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// PUT makes a PUT request to the specified URL.
// @params url or  (String) - path to file to write into.
func PUT(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod("PUT", ctx, args)
}
