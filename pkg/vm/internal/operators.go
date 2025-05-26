package internal

import (
	"context"
	"strings"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func Contains(ctx context.Context, input runtime.Value, value runtime.Value) runtime.Boolean {
	switch val := input.(type) {
	case runtime.List:
		idx, err := val.IndexOf(ctx, value)

		if err != nil {
			return runtime.False
		}

		return idx > -1
	case runtime.Map:
		containsValue, err := val.ContainsValue(ctx, value)

		if err != nil {
			return runtime.False
		}

		return containsValue
	case runtime.String:
		return runtime.Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}

func Add(_ context.Context, inputL, inputR runtime.Value) runtime.Value {
	left := runtime.ToNumberOrString(inputL)

	switch leftVal := left.(type) {
	case runtime.Int:
		return addLeftInt(leftVal, inputR)
	case runtime.Float:
		return addLeftFloat(leftVal, inputR)
	case runtime.String:
		return addLeftString(leftVal, inputR)
	default:
		return runtime.String(leftVal.String() + inputR.String())
	}
}

func addLeftInt(integer runtime.Int, input runtime.Value) runtime.Value {
	right := runtime.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return integer + rightVal
	case runtime.Float:
		return runtime.Float(integer) + rightVal
	default:
		return runtime.String(integer.String() + rightVal.String())
	}
}

func addLeftFloat(float runtime.Float, input runtime.Value) runtime.Value {
	right := runtime.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return float + runtime.Float(rightVal)
	case runtime.Float:
		return float + rightVal
	default:
		return runtime.String(float.String() + rightVal.String())
	}
}

func addLeftString(str runtime.String, input runtime.Value) runtime.Value {
	return runtime.String(str.String() + input.String())
}

func Subtract(ctx context.Context, inputL, inputR runtime.Value) runtime.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case runtime.Int:
		return subtractLeftInt(ctx, leftVal, inputR)
	case runtime.Float:
		return subtractLeftFloat(ctx, leftVal, inputR)
	default:
		return runtime.ZeroInt
	}
}

func subtractLeftInt(ctx context.Context, integer runtime.Int, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return integer - rightVal
	case runtime.Float:
		return runtime.Float(integer) - rightVal
	default:
		return runtime.ZeroInt
	}
}

func subtractLeftFloat(ctx context.Context, float runtime.Float, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return float - runtime.Float(rightVal)
	case runtime.Float:
		return float - rightVal
	default:
		return runtime.ZeroInt
	}
}

func Multiply(ctx context.Context, inputL, inputR runtime.Value) runtime.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case runtime.Int:
		return multiplyLeftInt(ctx, leftVal, inputR)
	case runtime.Float:
		return multiplyLeftFloat(ctx, leftVal, inputR)
	default:
		return runtime.ZeroInt
	}
}

func multiplyLeftInt(ctx context.Context, integer runtime.Int, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return integer * rightVal
	case runtime.Float:
		return runtime.Float(integer) * rightVal
	default:
		return runtime.ZeroInt
	}
}

func multiplyLeftFloat(ctx context.Context, float runtime.Float, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return float * runtime.Float(rightVal)
	case runtime.Float:
		return float * rightVal
	default:
		return runtime.ZeroInt
	}
}

func Divide(ctx context.Context, inputL, inputR runtime.Value) runtime.Value {
	left := ToNumberOnly(ctx, inputL)

	switch leftVal := left.(type) {
	case runtime.Int:
		return divideLeftInt(ctx, leftVal, inputR)
	case runtime.Float:
		return divideLeftFloat(ctx, leftVal, inputR)
	default:
		return runtime.ZeroInt
	}
}

func divideLeftInt(ctx context.Context, integer runtime.Int, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return integer / rightVal
	case runtime.Float:
		return runtime.Float(integer) / rightVal
	default:
		return runtime.ZeroInt
	}
}

func divideLeftFloat(ctx context.Context, float runtime.Float, input runtime.Value) runtime.Value {
	right := ToNumberOnly(ctx, input)

	switch rightVal := right.(type) {
	case runtime.Int:
		return float / runtime.Float(rightVal)
	case runtime.Float:
		return float / rightVal
	default:
		return runtime.ZeroInt
	}
}

func Modulus(ctx context.Context, inputL, inputR runtime.Value) runtime.Value {
	left, _ := runtime.ToInt(ctx, inputL)
	right, _ := runtime.ToInt(ctx, inputR)

	return left % right
}

func Increment(ctx context.Context, input runtime.Value) runtime.Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case runtime.Int:
		return value + 1
	case runtime.Float:
		return value + 1
	default:
		return runtime.None
	}
}

func Decrement(ctx context.Context, input runtime.Value) runtime.Value {
	left := ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case runtime.Int:
		return value - 1
	case runtime.Float:
		return value - 1
	default:
		return runtime.None
	}
}

func ToRange(ctx context.Context, left, right runtime.Value) (runtime.Value, error) {
	start, err := runtime.ToInt(ctx, left)

	if err != nil {
		return runtime.ZeroInt, err
	}

	end, err := runtime.ToInt(ctx, right)

	if err != nil {
		return runtime.ZeroInt, err
	}

	return NewRange(int64(start), int64(end)), nil
}

func Like(left, right runtime.Value) (runtime.Boolean, error) {
	if err := runtime.AssertString(left); err != nil {
		// TODO: Return the error? AQL just returns false
		return runtime.False, nil
	}

	if err := runtime.AssertString(right); err != nil {
		// TODO: Return the error? AQL just returns false
		return runtime.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return runtime.False, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return runtime.NewBoolean(result), nil
}
