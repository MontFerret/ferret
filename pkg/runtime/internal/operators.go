package internal

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func Add(inputL, inputR core.Value) core.Value {
	left := ToNumberOrString(inputL)

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
	right := ToNumberOrString(input)

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
	right := ToNumberOrString(input)

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

func Subtract(inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return subtractLeftInt(leftVal, inputR)
	case core.Float:
		return subtractLeftFloat(leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func subtractLeftInt(integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer - rightVal
	case core.Float:
		return core.Float(integer) - rightVal
	default:
		return core.ZeroInt
	}
}

func subtractLeftFloat(float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return float - core.Float(rightVal)
	case core.Float:
		return float - rightVal
	default:
		return core.ZeroInt
	}
}

func Multiply(inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return multiplyLeftInt(leftVal, inputR)
	case core.Float:
		return multiplyLeftFloat(leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func multiplyLeftInt(integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer * rightVal
	case core.Float:
		return core.Float(integer) * rightVal
	default:
		return core.ZeroInt
	}
}

func multiplyLeftFloat(float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return float * core.Float(rightVal)
	case core.Float:
		return float * rightVal
	default:
		return core.ZeroInt
	}
}

func Divide(inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case core.Int:
		return divideLeftInt(leftVal, inputR)
	case core.Float:
		return divideLeftFloat(leftVal, inputR)
	default:
		return core.ZeroInt
	}
}

func divideLeftInt(integer core.Int, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return integer / rightVal
	case core.Float:
		return core.Float(integer) / rightVal
	default:
		return core.ZeroInt
	}
}

func divideLeftFloat(float core.Float, input core.Value) core.Value {
	right := ToNumberOnly(input)

	switch rightVal := right.(type) {
	case core.Int:
		return float / core.Float(rightVal)
	case core.Float:
		return float / rightVal
	default:
		return core.ZeroInt
	}
}

func Modulus(inputL, inputR core.Value) core.Value {
	left := ToInt(inputL)
	right := ToInt(inputR)

	return left % right
}

func Increment(input core.Value) core.Value {
	left := ToNumberOnly(input)

	switch value := left.(type) {
	case core.Int:
		return value + 1
	case core.Float:
		return value + 1
	default:
		return core.None
	}
}

func Decrement(input core.Value) core.Value {
	left := ToNumberOnly(input)

	switch value := left.(type) {
	case core.Int:
		return value - 1
	case core.Float:
		return value - 1
	default:
		return core.None
	}
}

func Range(left, right core.Value) (core.Value, error) {
	start := ToInt(left)
	end := ToInt(right)

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
