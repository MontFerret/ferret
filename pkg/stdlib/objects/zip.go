package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

/*
 * Returns an object assembled from the separate parameters keys and values.
 * Keys and values must be arrays and have the same length.
 * @params keys (Array of Objects) - result object keys.
 * @params values (Array of Objects) - result object values.
 * @returns (Object) - Object assembled from the separate parameters keys and values.
 */
func Zip(_ context.Context, args ...core.Value) (core.Value, error) {
	return nil, nil
}
