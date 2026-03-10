package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func ToRange(ctx context.Context, left, right runtime.Value) (runtime.Value, error) {
	start, err := runtime.ToInt(ctx, left)

	if err != nil {
		return runtime.ZeroInt, err
	}

	end, err := runtime.ToInt(ctx, right)

	if err != nil {
		return runtime.ZeroInt, err
	}

	return data.NewRange(int64(start), int64(end)), nil
}
