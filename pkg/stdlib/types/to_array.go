package types

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// toArray takes an input value of any type and convert it into an array value.
// @param (Value) - Input value of arbitrary type.
// @returns (Array)
// None is converted to an empty array
// Boolean values, numbers and strings are converted to an array containing the original value as its single element
// Arrays keep their original value
// Objects / HTML nodes are converted to an array containing their attribute values as array elements.
func ToArray(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	input := args[0]

	switch input.Type() {
	case core.BooleanType,
		core.IntType,
		core.FloatType,
		core.StringType,
		core.DateTimeType:

		return values.NewArrayWith(input), nil
	case core.ArrayType:
		return input.Copy(), nil
	case core.ObjectType:
		obj, ok := input.(*values.Object)

		if !ok {
			return values.NewArray(0), nil
		}

		arr := values.NewArray(int(obj.Length()))

		obj.ForEach(func(value core.Value, key string) bool {
			arr.Push(value)

			return true
		})

		return obj, nil
	default:
		iterable, ok := input.(collections.IterableCollection)

		if !ok {
			return values.NewArray(0), nil
		}

		iterator, err := iterable.Iterate(ctx)

		if err != nil {
			return values.None, err
		}

		arr := values.NewArray(20)

		for {
			val, _, err := iterator.Next(ctx)

			if err != nil {
				return nil, err
			}

			// end of iteration
			if val == values.None {
				break
			}

			arr.Push(val)
		}

		return arr, nil
	}
}
