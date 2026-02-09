package runtime

import (
	"context"
)

func Add(_ context.Context, inputL, inputR Value) Value {
	left := ToNumberOrString(inputL)

	switch leftVal := left.(type) {
	case Int:
		return AddLeftInt(leftVal, inputR)
	case Float:
		return AddLeftFloat(leftVal, inputR)
	case String:
		return addLeftString(leftVal, inputR)
	default:
		return String(leftVal.String() + inputR.String())
	}
}

func AddLeftInt(integer Int, input Value) Value {
	right := ToNumberOrString(input)

	switch rightVal := right.(type) {
	case Int:
		return integer + rightVal
	case Float:
		return Float(integer) + rightVal
	default:
		return String(integer.String() + rightVal.String())
	}
}

func AddLeftFloat(float Float, input Value) Value {
	right := ToNumberOrString(input)

	switch rightVal := right.(type) {
	case Int:
		return float + Float(rightVal)
	case Float:
		return float + rightVal
	default:
		return String(float.String() + rightVal.String())
	}
}

func addLeftString(str String, input Value) Value {
	return String(str.String() + input.String())
}

func Subtract(ctx context.Context, inputL, inputR Value) Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case Int:
		return subtractLeftInt(ctx, leftVal, inputR)
	case Float:
		return subtractLeftFloat(ctx, leftVal, inputR)
	default:
		return ZeroInt
	}
}

func subtractLeftInt(ctx context.Context, integer Int, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		return integer - rightVal
	case Float:
		return Float(integer) - rightVal
	default:
		return ZeroInt
	}
}

func subtractLeftFloat(ctx context.Context, float Float, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		return float - Float(rightVal)
	case Float:
		return float - rightVal
	default:
		return ZeroInt
	}
}

func Multiply(ctx context.Context, inputL, inputR Value) Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case Int:
		return multiplyLeftInt(ctx, leftVal, inputR)
	case Float:
		return multiplyLeftFloat(ctx, leftVal, inputR)
	default:
		return ZeroInt
	}
}

func multiplyLeftInt(ctx context.Context, integer Int, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		return integer * rightVal
	case Float:
		return Float(integer) * rightVal
	default:
		return ZeroInt
	}
}

func multiplyLeftFloat(ctx context.Context, float Float, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		return float * Float(rightVal)
	case Float:
		return float * rightVal
	default:
		return ZeroInt
	}
}

func Divide(ctx context.Context, inputL, inputR Value) Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case Int:
		return divideLeftInt(ctx, leftVal, inputR)
	case Float:
		return divideLeftFloat(ctx, leftVal, inputR)
	default:
		return ZeroInt
	}
}

func divideLeftInt(ctx context.Context, integer Int, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		if rightVal != 0 && integer%rightVal != 0 {
			return Float(integer) / Float(rightVal)
		}
		return integer / rightVal
	case Float:
		return Float(integer) / rightVal
	default:
		return ZeroInt
	}
}

func divideLeftFloat(ctx context.Context, float Float, input Value) Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case Int:
		return float / Float(rightVal)
	case Float:
		return float / rightVal
	default:
		return ZeroInt
	}
}

func Modulus(ctx context.Context, inputL, inputR Value) Value {
	left, _ := ToInt(ctx, inputL)
	right, _ := ToInt(ctx, inputR)

	return left % right
}

func Increment(ctx context.Context, input Value) Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case Int:
		return value + 1
	case Float:
		return value + 1
	default:
		return None
	}
}

func Decrement(ctx context.Context, input Value) Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case Int:
		return value - 1
	case Float:
		return value - 1
	default:
		return None
	}
}
