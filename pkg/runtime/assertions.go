package runtime

import "context"

type TypeAssertion func(input Value) error

func AssertString(input Value) error {
	_, ok := input.(String)

	if !ok {
		return TypeError(
			input,
			TypeString,
		)
	}

	return nil
}

func AssertInt(input Value) error {
	_, ok := input.(Int)

	if !ok {
		return TypeError(
			input,
			TypeInt,
		)
	}

	return nil
}

func AssertFloat(input Value) error {
	_, ok := input.(Float)

	if !ok {
		return TypeError(
			input,
			TypeFloat,
		)
	}

	return nil
}

func AssertNumber(input Value) error {
	switch input.(type) {
	case Int, Float:
		return nil
	default:
		return TypeError(input, TypeInt, TypeFloat)
	}
}

func AssertBoolean(input Value) error {
	_, ok := input.(Boolean)

	if !ok {
		return TypeError(
			input,
			TypeBoolean,
		)
	}

	return nil
}

func AssertList(input Value) error {
	_, ok := input.(List)

	if !ok {
		return TypeError(
			input,
			"items",
		)
	}

	return nil
}

func AssertItemsOf(ctx context.Context, input Iterable, assertion TypeAssertion) error {
	return ForEachOf(ctx, input, func(ctx context.Context, value, _ Value) (Boolean, error) {
		if err := assertion(value); err != nil {
			return false, err
		}

		return true, nil
	})
}

func AssertMap(input Value) error {
	_, ok := input.(Map)

	if !ok {
		return TypeError(
			input,
			TypeMap,
		)
	}

	return nil
}

func AssertBinary(input Value) error {
	_, ok := input.(*Binary)

	if !ok {
		return TypeError(
			input,
			TypeBinary,
		)
	}

	return nil
}

func AssertDateTime(input Value) error {
	_, ok := input.(DateTime)

	if !ok {
		return TypeError(
			input,
			TypeDateTime,
		)
	}

	return nil
}
