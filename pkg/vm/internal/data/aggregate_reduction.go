package data

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ReduceAggregate finalizes a collected list with the same semantics as native
// aggregate collectors. The list and its elements remain borrowed. The caller
// validates the aggregate kind while preparing the program.
func ReduceAggregate(ctx context.Context, values runtime.List, kind bytecode.AggregateKind) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	var state aggregateState

	err := values.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		updateAggregateState(&state, kind, value)

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return aggregateValueFor(state, kind), nil
}
