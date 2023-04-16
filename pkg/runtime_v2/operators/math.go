package operators

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Add(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOrString(inputL)
	right := values.ToNumberOrString(inputR)

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

	return values.String(left.String() + right.String())
}

func Subtract(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)
	right := values.ToNumberOnly(inputR)

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

	return values.ZeroInt
}

func Multiply(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)
	right := values.ToNumberOnly(inputR)

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

func Divide(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)
	right := values.ToNumberOnly(inputR)

	if left.Type() == types.Int {
		if right.Type() == types.Int {
			l := values.Float(left.(values.Int))
			r := values.Float(right.(values.Int))

			if r == 0.0 {
				panic("divide by zero")
			}

			return l / r
		}

		if right.Type() == types.Float {
			l := values.Float(left.(values.Int))
			r := right.(values.Float)

			if r == 0.0 {
				panic("divide by zero")
			}

			return l / r
		}
	}

	if left.Type() == types.Float {
		if right.Type() == types.Float {
			l := left.(values.Float)
			r := right.(values.Float)

			if r == 0.0 {
				panic("divide by zero")
			}

			return l / r
		}

		if right.Type() == types.Int {
			l := left.(values.Float)
			r := values.Float(right.(values.Int))

			if r == 0.0 {
				panic("divide by zero")
			}

			return l / r
		}
	}

	return values.ZeroInt
}

func Modulus(inputL, inputR core.Value) core.Value {
	left := values.ToNumberOnly(inputL)
	right := values.ToNumberOnly(inputR)

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

func Increment(input core.Value) core.Value {
	left := values.ToNumberOnly(input)

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

func Decrement(input core.Value) core.Value {
	left := values.ToNumberOnly(input)

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
