package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func AssertString(input core.Value) error {
	_, ok := input.(String)

	if !ok {
		return core.TypeError(
			input,
			types.String,
		)
	}

	return nil
}

func AssertInt(input core.Value) error {
	_, ok := input.(Int)

	if !ok {
		return core.TypeError(
			input,
			types.Int,
		)
	}

	return nil
}

func AssertFloat(input core.Value) error {
	_, ok := input.(Float)

	if !ok {
		return core.TypeError(
			input,
			types.Float,
		)
	}

	return nil
}

func AssertNumber(input core.Value) error {
	switch input.(type) {
	case Int, Float:
		return nil
	default:
		return core.TypeError(input, types.Int, types.Float)
	}
}

func AssertBoolean(input core.Value) error {
	_, ok := input.(Boolean)

	if !ok {
		return core.TypeError(
			input,
			types.Boolean,
		)
	}

	return nil
}

func AssertArray(input core.Value) error {
	_, ok := input.(*Array)

	if !ok {
		return core.TypeError(
			input,
			types.Array,
		)
	}

	return nil
}

func AssertObject(input core.Value) error {
	_, ok := input.(*Object)

	if !ok {
		return core.TypeError(
			input,
			types.Object,
		)
	}

	return nil
}

func AssertBinary(input core.Value) error {
	_, ok := input.(*Binary)

	if !ok {
		return core.TypeError(
			input,
			types.Binary,
		)
	}

	return nil
}

func AssertDateTime(input core.Value) error {
	_, ok := input.(DateTime)

	if !ok {
		return core.TypeError(
			input,
			types.DateTime,
		)
	}

	return nil
}
