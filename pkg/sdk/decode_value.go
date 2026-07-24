package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DecodeValue binds a Ferret runtime value into T and returns the result.
func DecodeValue[T any](ctx context.Context, src runtime.Value, options ...DecodeOption) (T, error) {
	var output T

	if err := Decode(ctx, src, &output, options...); err != nil {
		return output, err
	}

	return output, nil
}
