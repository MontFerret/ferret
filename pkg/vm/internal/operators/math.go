package operators

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal/data"
)

func Add(_ context.Context, inputL, inputR runtime.Value) runtime.Value {
	left := ToNumberOrString(inputL)

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
	right := ToNumberOrString(input)

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
	right := ToNumberOrString(input)

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
	left := data.ToNumberOnly(ctx, inputL)

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
	right := data.ToNumberOnly(ctx, input)

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
	right := data.ToNumberOnly(ctx, input)

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
	left := data.ToNumberOnly(ctx, inputL)

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
	right := data.ToNumberOnly(ctx, input)

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
	right := data.ToNumberOnly(ctx, input)

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
	left := data.ToNumberOnly(ctx, inputL)

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
	right := data.ToNumberOnly(ctx, input)

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
	right := data.ToNumberOnly(ctx, input)

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
	left := data.ToNumberOnly(ctx, input)

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
	left := data.ToNumberOnly(ctx, input)

	switch value := left.(type) {
	case runtime.Int:
		return value - 1
	case runtime.Float:
		return value - 1
	default:
		return runtime.None
	}
}
