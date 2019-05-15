package operators

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	OperatorFunc func(left, right core.Value) core.Value
	baseOperator struct {
		src   core.SourceMap
		left  core.Expression
		right core.Expression
	}
)

func (operator *baseOperator) Exec(_ context.Context, _ *core.Scope) (core.Value, error) {
	return values.None, core.ErrInvalidOperation
}

func (operator *baseOperator) Eval(_ context.Context, _, _ core.Value) (core.Value, error) {
	return values.None, core.ErrInvalidOperation
}

// Equality
func Equal(left, right core.Value) core.Value {
	if left.Compare(right) == 0 {
		return values.True
	}

	return values.False
}

func NotEqual(left, right core.Value) core.Value {
	if left.Compare(right) != 0 {
		return values.True
	}

	return values.False
}

func Less(left, right core.Value) core.Value {
	if left.Compare(right) < 0 {
		return values.True
	}

	return values.False
}

func LessOrEqual(left, right core.Value) core.Value {
	out := left.Compare(right)

	if out < 0 || out == 0 {
		return values.True
	}

	return values.False
}

func Greater(left, right core.Value) core.Value {
	if left.Compare(right) > 0 {
		return values.True
	}

	return values.False
}

func GreaterOrEqual(left, right core.Value) core.Value {
	out := left.Compare(right)

	if out > 0 || out == 0 {
		return values.True
	}

	return values.False
}

func Not(left, _ core.Value) core.Value {
	b := values.ToBoolean(left)

	if b == values.True {
		return values.False
	}

	return values.True
}

func ToNumberOrString(input core.Value) core.Value {
	switch input.Type() {
	case types.Int, types.Float, types.String:
		return input
	default:
		return values.ToInt(input)
	}
}

func ToNumberOnly(input core.Value) core.Value {
	switch input.Type() {
	case types.Int, types.Float:
		return input
	case types.String:
		if strings.Contains(input.String(), ".") {
			return values.ToFloat(input)
		}

		return values.ToInt(input)
	case types.Array:
		arr := input.(*values.Array)
		length := arr.Length()

		if length == 0 {
			return values.ZeroInt
		}

		i := values.ZeroInt
		f := values.ZeroFloat

		for y := values.Int(0); y < length; y++  {
			out := ToNumberOnly(arr.Get(y))

			if out.Type() == types.Int {
				i += out.(values.Int)
			} else {
				f += out.(values.Float)
			}
		}

		if f == 0 {
			return i
		}

		return values.Float(i) + f
	default:
		return values.ToInt(input)
	}
}

// Adds numbers
// Concatenates strings
func Add(inputL, inputR core.Value) core.Value {
	left := ToNumberOrString(inputL)
	right := ToNumberOrString(inputR)

	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := left.(values.Int)
			r := right.(values.Int)

			return l + r
		}

		if right.Type() == types.Float {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) + r
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			return l + r
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := right.(values.Int)

			return l + values.Float(r)
		}
	}

	return values.NewString(left.String() + right.String())
}

func Subtract(inputL, inputR core.Value) core.Value {
	left := ToNumberOnly(inputL)
	right := ToNumberOnly(inputR)

	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := left.(values.Int)
			r := right.(values.Int)

			return l - r
		}

		if right.Type() == types.Float {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) - r
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			return l - r
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := right.(values.Int)

			return l - values.Float(r)
		}
	}

	return values.NaN
}

func Multiply(left, right core.Value) core.Value {
	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := left.(values.Int)
			r := right.(values.Int)

			return l * r
		}

		if right.Type() == types.Float {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) * r
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			return l * r
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := right.(values.Int)

			return l * values.Float(r)
		}
	}

	return values.ZeroInt
}

func Divide(left, right core.Value) core.Value {
	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := left.(values.Int)
			r := right.(values.Int)

			return l / r
		}

		if right.Type() == types.Float {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) / r
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			return l / r
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := right.(values.Int)

			return l / values.Float(r)
		}
	}

	return values.ZeroInt
}

func Modulus(left, right core.Value) core.Value {
	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := left.(values.Int)
			r := right.(values.Int)

			return l % r
		}

		if right.Type() == types.Float {
			l := left.(values.Int)
			r := right.(values.Float)

			return l % values.Int(r)
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			return values.Int(l) % values.Int(r)
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := right.(values.Int)

			return values.Int(l) % r
		}
	}

	return values.ZeroInt
}

func Increment(left, _ core.Value) core.Value {
	if left.Type() == types.Int {
		l := left.(values.Int)

		return l + 1
	}

	if left.Type() == types.Float {
		l := left.(values.Float)

		return l + 1
	}

	return values.None
}

func Decrement(left, _ core.Value) core.Value {
	if left.Type() == types.Int {
		l := left.(values.Int)

		return l - 1
	}

	if left.Type() == types.Float {
		l := left.(values.Float)

		return l - 1
	}

	return values.None
}

func Negative(value, _ core.Value) core.Value {
	err := core.ValidateType(value, types.Int, types.Float)

	if err != nil {
		return values.ZeroInt
	}

	if value.Type() == types.Int {
		return -value.(values.Int)
	}

	return -value.(values.Float)
}

func Positive(value, _ core.Value) core.Value {
	err := core.ValidateType(value, types.Int, types.Float)

	if err != nil {
		return values.ZeroInt
	}

	if value.Type() == types.Int {
		return +value.(values.Int)
	}

	return +value.(values.Float)
}

func ToBoolean(value, _ core.Value) core.Value {
	return values.ToBoolean(value)
}
