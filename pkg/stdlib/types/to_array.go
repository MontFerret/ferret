package types

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Takes an input value of any type and convert it into an array value.
 * @param (Value) - Input value of arbitrary type.
 * @returns (Array)
 * None is converted to an empty array
 * Boolean values, numbers and strings are converted to an array containing the original value as its single element
 * Arrays keep their original value
 * Objects / HTML nodes are converted to an array containing their attribute values as array elements.
 */
func ToArray(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	arg := args[0]

	switch arg.Type() {
	case core.BooleanType,
		core.IntType,
		core.FloatType,
		core.StringType,
		core.DateTimeType:
		return values.NewArrayWith(arg), nil
	case core.HTMLElementType,
		core.HTMLDocumentType:
		val := arg.(values.HTMLNode)
		attrs := val.GetAttributes()

		obj, ok := attrs.(*values.Object)

		if !ok {
			return values.NewArray(0), nil
		}

		return collections.ToArray(collections.NewObjectIterator(obj))
	case core.ArrayType:
		return arg, nil
	case core.ObjectType:
		return collections.ToArray(collections.NewObjectIterator(arg.(*values.Object)))
	default:
		return values.NewArray(0), nil
	}
}
