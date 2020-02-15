package http

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// POST makes a POST request to the specified URL.
// @params url or  (String) - path to file to write into.
func POST(ctx context.Context, args ...core.Value) (core.Value, error) {
	return execMethod("POST", ctx, args)
}
