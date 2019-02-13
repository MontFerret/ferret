package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Length(_ context.Context, inputs ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(inputs, 1, 1)

	if err != nil {
		return values.None, err
	}

	value := inputs[0]

	c, ok := value.(collections.Collection)

	if !ok {
		return values.None, core.TypeError(value.Type(),
			types.String,
			types.Array,
			types.Object,
			types.HTMLElement,
			types.HTMLDocument,
			types.Binary,
			core.NewType("Collection"),
		)
	}

	return c.Length(), nil
}
