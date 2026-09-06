package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func copyListValues(ctx context.Context, values runtime.List) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	length, err := values.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if length < 0 {
		return runtime.None, runtime.Error(runtime.ErrInvalidArgument, "list length must not be negative")
	}

	result := runtime.NewArray64(length)
	err = values.ForEach(ctx, func(ctx context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
		if err := ctx.Err(); err != nil {
			return false, err
		}

		copied, err := runtime.CloneOrCopy(ctx, value)
		if err != nil {
			return false, fmt.Errorf("item %d: %w", index, err)
		}

		if err := result.Append(ctx, copied); err != nil {
			return false, fmt.Errorf("item %d: %w", index, err)
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	return result, nil
}
