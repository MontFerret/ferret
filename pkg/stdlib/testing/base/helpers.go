package base

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TypeAssertion(expected runtime.Type) Assertion {
	var assert runtime.TypeAssertion

	// We do this in the constructor to avoid unnecessary type assertions on every assertion call.
	// This way we only do it once when the assertion is created.
	// Also, this allows us to panic early if an unsupported type is used, rather than at runtime when the assertion is called.
	switch expected {
	case runtime.TypeNone:
		assert = runtime.AssertNone
	case runtime.TypeString:
		assert = runtime.AssertString
	case runtime.TypeInt:
		assert = runtime.AssertInt
	case runtime.TypeFloat:
		assert = runtime.AssertFloat
	case runtime.TypeBoolean:
		assert = runtime.AssertBoolean
	case runtime.TypeArray:
		assert = runtime.AssertArray
	case runtime.TypeObject:
		assert = runtime.AssertObject
	case runtime.TypeDateTime:
		assert = runtime.AssertDateTime
	case runtime.TypeBinary:
		assert = runtime.AssertBinary
	default:
		panic(fmt.Sprintf("unsupported type assertion for %s", expected))
	}

	return Assertion{
		DefaultMessage: func(args []runtime.Value) string {
			return fmt.Sprintf("be %s", expected)
		},
		Args: Args{
			Min: 1,
			Max: 2,
		},
		Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {

			err := assert(args[0])

			return err == nil, nil
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
		Fn: func(ctx context.Context, args []runtime.Value) (bool, error) {
			if len(args) != 2 {
				return false, fmt.Errorf("expected 2 arguments, got %d", len(args))
			}

			return op.Compare(args)
		},
	}
}
