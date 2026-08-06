package runtime

import "context"

type Number interface {
	Int | Float
}

func ToNumber(ctx context.Context, input Value) (Value, error) {
	switch value := input.(type) {
	case Int, Float:
		return value, nil
	default:
		return ToFloat(ctx, input)
	}
}
