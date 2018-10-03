package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func Outersection(_ context.Context, args ...core.Value) (core.Value, error) {
	return sections(args, 1)
}
