package testing

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func hasFailureMessage(ctx context.Context, descriptor assertion, args []runtime.Value, positive bool) string {
	message := descriptor.defaultFailureMessage(ctx, args, positive)
	if !positive || len(args) == descriptor.args.max {
		return message
	}

	keys, ok := args[1].(runtime.List)
	if !ok {
		return message
	}

	target := args[0].(runtime.Map)
	missing := runtime.NewArray(0)
	err := keys.ForEach(ctx, func(ctx context.Context, key runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		contains, err := target.ContainsKey(ctx, key)
		if err != nil {
			return false, err
		}

		if !contains {
			if err := missing.Append(ctx, key); err != nil {
				return false, err
			}
		}

		return true, nil
	})
	if err != nil {
		return message
	}

	length, err := missing.Length(ctx)
	if err != nil || length == 0 {
		return message
	}

	return message + "\nmissing: " + formatValue(ctx, missing)
}
