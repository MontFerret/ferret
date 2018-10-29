package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func NOOP(_ context.Context, args ...core.Value) (core.Value, error) {
	return values.None, nil
}
