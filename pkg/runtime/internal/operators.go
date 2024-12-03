package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"
)

func Add(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOrString(inputL)

	switch leftVal := left.(type) {
	case values.Int:
		return addLeftInt(leftVal, inputR)
	case values.Float:
		return addLeftFloat(leftVal, inputR)
	case values.String:
		return addLeftString(leftVal, inputR)
	default:
		return values.String(leftVal.String() + inputR.String())
	}
}

func addLeftInt(integer values.Int, input core.Value) core.Value {
	right := values.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case values.Int:
		return integer + rightVal
	case values.Float:
		return values.Float(integer) + rightVal
	default:
		return values.String(integer.String() + rightVal.String())
	}
}

func addLeftFloat(float values.Float, input core.Value) core.Value {
	right := values.ToNumberOrString(input)

	switch rightVal := right.(type) {
	case values.Int:
		return float + values.Float(rightVal)
	case values.Float:
		return float + rightVal
	default:
		return values.String(float.String() + rightVal.String())
	}
}

func addLeftString(str values.String, input core.Value) core.Value {
	return values.String(str.String() + input.String())
}

func Subtract(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case values.Int:
		return subtractLeftInt(leftVal, inputR)
	case values.Float:
		return subtractLeftFloat(leftVal, inputR)
	default:
		return values.ZeroInt
	}
}

func subtractLeftInt(integer values.Int, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return integer - rightVal
	case values.Float:
		return values.Float(integer) - rightVal
	default:
		return values.ZeroInt
	}
}

func subtractLeftFloat(float values.Float, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return float - values.Float(rightVal)
	case values.Float:
		return float - rightVal
	default:
		return values.ZeroInt
	}
}

func Multiply(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case values.Int:
		return multiplyLeftInt(leftVal, inputR)
	case values.Float:
		return multiplyLeftFloat(leftVal, inputR)
	default:
		return values.ZeroInt
	}
}

func multiplyLeftInt(integer values.Int, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return integer * rightVal
	case values.Float:
		return values.Float(integer) * rightVal
	default:
		return values.ZeroInt
	}
}

func multiplyLeftFloat(float values.Float, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return float * values.Float(rightVal)
	case values.Float:
		return float * rightVal
	default:
		return values.ZeroInt
	}
}

func Divide(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)

	switch leftVal := left.(type) {
	case values.Int:
		return divideLeftInt(leftVal, inputR)
	case values.Float:
		return divideLeftFloat(leftVal, inputR)
	default:
		return values.ZeroInt
	}
}

func divideLeftInt(integer values.Int, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return integer / rightVal
	case values.Float:
		return values.Float(integer) / rightVal
	default:
		return values.ZeroInt
	}
}

func divideLeftFloat(float values.Float, input core.Value) core.Value {
	right := values.ToNumberOnly(input)

	switch rightVal := right.(type) {
	case values.Int:
		return float / values.Float(rightVal)
	case values.Float:
		return float / rightVal
	default:
		return values.ZeroInt
	}
}

func Modulus(inputL, inputR core.Value) core.Value {
	left := values.ToInt(inputL)
	right := values.ToInt(inputR)

	return left % right
}

func Increment(input core.Value) core.Value {
	left := values.ToNumberOnly(input)

	switch value := left.(type) {
	case values.Int:
		return value + 1
	case values.Float:
		return value + 1
	default:
		return values.None
	}
}

func Decrement(input core.Value) core.Value {
	left := values.ToNumberOnly(input)

	switch value := left.(type) {
	case values.Int:
		return value - 1
	case values.Float:
		return value - 1
	default:
		return values.None
	}
}

func Range(left, right core.Value) (core.Value, error) {
	start := values.ToInt(left)
	end := values.ToInt(right)

	return values.NewRange(int64(start), int64(end)), nil
}

func Like(left, right core.Value) (values.Boolean, error) {
	if err := values.AssertString(left); err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	if err := values.AssertString(right); err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return values.False, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return values.NewBoolean(result), nil
}
