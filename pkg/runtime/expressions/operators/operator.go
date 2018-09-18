package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Operator func(left, right core.Value) core.Value

type baseOperator struct {
	src   core.SourceMap
	left  core.Expression
	right core.Expression
}

func (operator *baseOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
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
	if left == values.True {
		return values.False
	} else if left == values.False {
		return values.True
	}

	return values.False
}

// Adds numbers
// Concats strings
func Add(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(int)

			return values.NewInt(l + r)
		}

		if right.Type() == core.FloatType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(float64)

			return values.Float(float64(l) + r)
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(float64)

			return values.Float(l + r)
		}

		if right.Type() == core.IntType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(int)

			return values.Float(l + float64(r))
		}
	}

	return values.NewString(left.String() + right.String())
}

func Subtract(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(int)

			return values.NewInt(l - r)
		}

		if right.Type() == core.FloatType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(float64)

			return values.Float(float64(l) - r)
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(float64)

			return values.Float(l - r)
		}

		if right.Type() == core.IntType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(int)

			return values.Float(l - float64(r))
		}
	}

	return values.ZeroInt
}

func Multiply(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(int)

			return values.NewInt(l * r)
		}

		if right.Type() == core.FloatType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(float64)

			return values.Float(float64(l) * r)
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(float64)

			return values.Float(l * r)
		}

		if right.Type() == core.IntType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(int)

			return values.Float(l * float64(r))
		}
	}

	return values.ZeroInt
}

func Divide(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(int)

			return values.NewInt(l / r)
		}

		if right.Type() == core.FloatType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(float64)

			return values.Float(float64(l) / r)
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(float64)

			return values.Float(l / r)
		}

		if right.Type() == core.IntType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(int)

			return values.Float(l / float64(r))
		}
	}

	return values.ZeroInt
}

func Modulus(left, right core.Value) core.Value {
	if left.Type() == core.IntType {
		if right.Type() == core.IntType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(int)

			return values.NewInt(l % r)
		}

		if right.Type() == core.FloatType {
			l := left.Unwrap().(int)
			r := right.Unwrap().(float64)

			return values.Float(l % int(r))
		}
	}

	if left.Type() == core.FloatType {
		if right.Type() == core.FloatType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(float64)

			return values.Float(int(l) % int(r))
		}

		if right.Type() == core.IntType {
			l := left.Unwrap().(float64)
			r := right.Unwrap().(int)

			return values.Float(int(l) % r)
		}
	}

	return values.ZeroInt
}

func Increment(left, _ core.Value) core.Value {
	if left.Type() == core.IntType {
		l := left.Unwrap().(int)

		return values.NewInt(l + 1)
	}

	if left.Type() == core.FloatType {
		l := left.Unwrap().(float64)

		return values.Float(l + 1)
	}

	return values.None
}

func Decrement(left, _ core.Value) core.Value {
	if left.Type() == core.IntType {
		l := left.Unwrap().(int)

		return values.NewInt(l - 1)
	}

	if left.Type() == core.FloatType {
		l := left.Unwrap().(float64)

		return values.Float(l - 1)
	}

	return values.None
}
