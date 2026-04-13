package base

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TypeAssertion(expected runtime.Type) Assertion {
	if expected == nil {
		panic("unsupported type assertion for <nil>")
	}

	return Assertion{
		DefaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("be %s", expected)
		},
		Args: Args{
			Min: 1,
			Max: 2,
		},
		Fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return expected.Is(args[0]), nil
		},
	}
}

func EqualityAssertion(op CompareOperator) Assertion {
	return Assertion{
		DefaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("be %s %s", op, FormatValue(args[1]))
		},
		Args: Args{
			Min: 2,
			Max: 3,
		},
		Fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return op.Compare(args)
		},
	}
}
