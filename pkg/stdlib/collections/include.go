package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Includes(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	var result values.Boolean
	haystack := args[0]

	switch v := haystack.(type) {
	case values.String:
		result = v.Contains(values.NewString(args[1].String()))

		break
	case *values.Array:
		v.Find()
		break
	case *values.Object:
		break
	case values.Binary:
		break
	case core.Iterable:
	default:
		return values.None, core.TypeError(haystack.Type(),
			types.String,
			types.Array,
			types.Object,
			types.Binary,
			core.NewType("Iterable"),
		)
	}

	return result, nil
}
