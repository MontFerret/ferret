package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

func (operator *baseOperator) Eval(_ context.Context, left, right core.Value) (core.Value, error) {
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

// Adds numbers
// Concats strings
func Add(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.(values.Int)
			r := right.(values.Int)

			return l + r
		}

		if right.Type() == core.FloatType {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) + r
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.(values.Float)
			r := right.(values.Float)

			return l + r
		}

		if right.Type() == core.IntType {
			l := left.(values.Float)
			r := right.(values.Int)

			return l + values.Float(r)
		}
	}

	return values.NewString(left.String() + right.String())
}

func Subtract(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.(values.Int)
			r := right.(values.Int)

			return l - r
		}

		if right.Type() == core.FloatType {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) - r
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.(values.Float)
			r := right.(values.Float)

			return l - r
		}

		if right.Type() == core.IntType {
			l := left.(values.Float)
			r := right.(values.Int)

			return l - values.Float(r)
		}
	}

	return values.ZeroInt
}

func Multiply(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.(values.Int)
			r := right.(values.Int)

			return l * r
		}

		if right.Type() == core.FloatType {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) * r
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.(values.Float)
			r := right.(values.Float)

			return l * r
		}

		if right.Type() == core.IntType {
			l := left.(values.Float)
			r := right.(values.Int)

			return l * values.Float(r)
		}
	}

	return values.ZeroInt
}

func Divide(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.(values.Int)
			r := right.(values.Int)

			return l / r
		}

		if right.Type() == core.FloatType {
			l := left.(values.Int)
			r := right.(values.Float)

			return values.Float(l) / r
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.(values.Float)
			r := right.(values.Float)

			return l / r
		}

		if right.Type() == core.IntType {
			l := left.(values.Float)
			r := right.(values.Int)

			return l / values.Float(r)
		}
	}

	return values.ZeroInt
}

func Modulus(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.(values.Int)
			r := right.(values.Int)

			return l % r
		}

		if right.Type() == core.FloatType {
			l := left.(values.Int)
			r := right.(values.Float)

			return l % values.Int(r)
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.(values.Float)
			r := right.(values.Float)

			return values.Int(l) % values.Int(r)
		}

		if right.Type() == core.IntType {
			l := left.(values.Float)
			r := right.(values.Int)

			return values.Int(l) % r
		}
	}

	return values.ZeroInt
}

func Increment(left, _ core.Value) core.Value {
	if left.Type() == core.IntType {
		l := left.(values.Int)

		return l + 1
	}

	if left.Type() == core.FloatType {
		l := left.(values.Float)

		return l + 1
	}

	return values.None
}

func Decrement(left, _ core.Value) core.Value {
	if left.Type() == core.IntType {
		l := left.(values.Int)

		return l - 1
	}

	if left.Type() == core.FloatType {
		l := left.(values.Float)

		return l - 1
	}

	return values.None
}
