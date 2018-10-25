package runner

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Assertions() map[string]core.Function {
	return map[string]core.Function{
		"EXPECT": expect,
	}
}

func expect(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	if args[0].Compare(args[1]) == 0 {
		return values.EmptyString, nil
	}

	return values.NewString(fmt.Sprintf(`expected "%s", but got "%s"`, args[0], args[1])), nil
}
