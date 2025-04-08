package internal

import (
	"context"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func Contains(ctx context.Context, input core.Value, value core.Value) core.Boolean {
	switch val := input.(type) {
	case core.List:
		contains, err := val.Contains(ctx, value)
		if err != nil {
			return core.False
		}

		return contains
	case core.Map:
		containsValue, err := val.ContainsValue(ctx, value)

		if err != nil {
			return core.False
		}

		return containsValue
	case core.String:
		return core.Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}

func Add(ctx context.Context, inputL, inputR core.Value) core.Value {
	left := core.ToNumberOrString(inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return addLeftInt(leftVal, inputR)
	case core.Float:
		return addLeftFloat(leftVal, inputR)
	case core.String:
		return addLeftString(leftVal, inputR)
	default:
		return core.String(leftVal.String() + inputR.String())
	}
}

func addLeftInt(integer core.Int, input core.Value) core.Value {
	right := core.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer + rightVal
	case core.Float:
		return core.Float(integer) + rightVal
	default:
		return core.String(integer.String() + rightVal.String())
	}
}

func addLeftFloat(float core.Float, input core.Value) core.Value {
	right := core.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case core.Int:
		return float + core.Float(rightVal)
	case core.Float:
		return float + rightVal
	default:
		return core.String(float.String() + rightVal.String())
	}
}

func addLeftString(str core.String, input core.Value) core.Value {
	return core.String(str.String() + input.String())
}

func Subtract(ctx context.Context, inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return subtractLeftInt(ctx, leftVal, inputR)
	case core.Float:
		return subtractLeftFloat(ctx, leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func subtractLeftInt(ctx context.Context, integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer - rightVal
	case core.Float:
		return core.Float(integer) - rightVal
	default:
		return core.ZeroInt
	}
}

func subtractLeftFloat(ctx context.Context, float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return float - core.Float(rightVal)
	case core.Float:
		return float - rightVal
	default:
		return core.ZeroInt
	}
}

func Multiply(ctx context.Context, inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return multiplyLeftInt(ctx, leftVal, inputR)
	case core.Float:
		return multiplyLeftFloat(ctx, leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func multiplyLeftInt(ctx context.Context, integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer * rightVal
	case core.Float:
		return core.Float(integer) * rightVal
	default:
		return core.ZeroInt
	}
}

func multiplyLeftFloat(ctx context.Context, float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return float * core.Float(rightVal)
	case core.Float:
		return float * rightVal
	default:
		return core.ZeroInt
	}
}

func Divide(ctx context.Context, inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return divideLeftInt(ctx, leftVal, inputR)
	case core.Float:
		return divideLeftFloat(ctx, leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func divideLeftInt(ctx context.Context, integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer / rightVal
	case core.Float:
		return core.Float(integer) / rightVal
	default:
		return core.ZeroInt
	}
}

func divideLeftFloat(ctx context.Context, float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case core.Int:
		return float / core.Float(rightVal)
	case core.Float:
		return float / rightVal
	default:
		return core.ZeroInt
	}
}

func Modulus(ctx context.Context, inputL, inputR core.Value) core.Value {
	left, _ := core.ToInt(ctx, inputL)
	right, _ := core.ToInt(ctx, inputR)

	return left % right
}

func Increment(ctx context.Context, input core.Value) core.Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case core.Int:
		return value + 1
	case core.Float:
		return value + 1
	default:
		return core.None
	}
}

func Decrement(ctx context.Context, input core.Value) core.Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case core.Int:
		return value - 1
	case core.Float:
		return value - 1
	default:
		return core.None
	}
}

func ToRange(ctx context.Context, left, right core.Value) (core.Value, error) {
	start, err := core.ToInt(ctx, left)

	if err != nil {
		return core.ZeroInt, err
	}

	end, err := core.ToInt(ctx, right)

	if err != nil {
		return core.ZeroInt, err
	}

	return NewRange(int64(start), int64(end)), nil
}

func Like(left, right core.Value) (core.Boolean, error) {
	if err := core.AssertString(left); err != nil {
		// TODO: Return the error? AQL just returns false
		return core.False, nil
	}

	if err := core.AssertString(right); err != nil {
		// TODO: Return the error? AQL just returns false
		return core.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return core.False, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return core.NewBoolean(result), nil
}
