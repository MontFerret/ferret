package runtime

// Add performs addition of two values.
// It supports addition for integers, floats, and strings.
func Add(_ Context, inputL, inputR Value) Value {
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

// AddLeftInt performs addition when the left operand is an integer.
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

// AddLeftFloat performs addition when the left operand is a float.
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

// Subtract performs subtraction of two values.
func Subtract(ctx Context, inputL, inputR Value) Value {
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

func subtractLeftInt(ctx Context, integer Int, input Value) Value {
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

func subtractLeftFloat(ctx Context, float Float, input Value) Value {
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

// Subtract performs subtraction of two values.
func Multiply(ctx Context, inputL, inputR Value) Value {
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

func multiplyLeftInt(ctx Context, integer Int, input Value) Value {
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

func multiplyLeftFloat(ctx Context, float Float, input Value) Value {
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

// Divide performs division of two values.
func Divide(ctx Context, inputL, inputR Value) Value {
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

func divideLeftInt(ctx Context, integer Int, input Value) Value {
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

func divideLeftFloat(ctx Context, float Float, input Value) Value {
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

// Modulus performs modulus of two values.
func Modulus(ctx Context, inputL, inputR Value) Value {
	left, _ := ToInt(ctx, inputL)
	right, _ := ToInt(ctx, inputR)

	return left % right
}

// Increment performs increment of a value by 1.
func Increment(ctx Context, input Value) Value {
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

// Decrement performs decrement of a value by 1.
func Decrement(ctx Context, input Value) Value {
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
