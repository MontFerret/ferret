package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

/*
 * Returns an object assembled from the separate parameters keys and values.
 * Keys and values must be arrays and have the same length.
 * @params keys (Array of Strings) - an array of strings, to be used as key names in the result.
 * @params values (Array of Objects) - an array of core.Value, to be used as key values.
 * @returns (Object) - an object with the keys and values assembled.
 */
func Zip(_ context.Context, args ...core.Value) (core.Value, error) {
	return nil, nil
}
