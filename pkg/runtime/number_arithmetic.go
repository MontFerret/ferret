package runtime

import (
	"context"
	"errors"
	"io"
	"math"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

func addNonTemporal(_ context.Context, inputL, inputR Value) (Value, error) {
	left := ToNumberOrString(inputL)

	switch leftVal := left.(type) {
	case Int:
		return addLeftInt(leftVal, inputR)
	case Float:
		return addLeftFloat(leftVal, inputR)
	case String:
		return addLeftString(leftVal, inputR), nil
	default:
		return String(leftVal.String() + inputR.String()), nil
	}
}

func addLeftInt(integer Int, input Value) (Value, error) {
	right := ToNumberOrString(input)

	switch rightVal := right.(type) {
	case Int:
		return addInts(integer, rightVal)
	case Float:
		return addFloats(Float(integer), rightVal, "addition")
	default:
		return String(integer.String() + rightVal.String()), nil
	}
}

func addLeftFloat(float Float, input Value) (Value, error) {
	right := ToNumberOrString(input)

	switch rightVal := right.(type) {
	case Int:
		return addFloats(float, Float(rightVal), "addition")
	case Float:
		return addFloats(float, rightVal, "addition")
	default:
		return String(float.String() + rightVal.String()), nil
	}
}

func addLeftString(str String, input Value) Value {
	return String(str.String() + input.String())
}

func subtractNonTemporal(ctx context.Context, inputL, inputR Value) (Value, error) {
	left, err := arithmeticNumber(ctx, inputL)
	if err != nil {
		return None, err
	}

	right, err := arithmeticNumber(ctx, inputR)
	if err != nil {
		return None, err
	}

	return subtractNumbers(left, right)
}

func multiplyNonTemporal(ctx context.Context, inputL, inputR Value) (Value, error) {
	left, err := arithmeticNumber(ctx, inputL)
	if err != nil {
		return None, err
	}

	right, err := arithmeticNumber(ctx, inputR)
	if err != nil {
		return None, err
	}

	return multiplyNumbers(left, right)
}

func divideNonTemporal(ctx context.Context, inputL, inputR Value) (Value, error) {
	left, err := arithmeticNumber(ctx, inputL)
	if err != nil {
		return None, err
	}

	right, err := arithmeticNumber(ctx, inputR)
	if err != nil {
		return None, err
	}

	return divideNumbers(left, right)
}

func moduloNonTemporal(ctx context.Context, inputL, inputR Value) (Value, error) {
	left, err := arithmeticInt(ctx, inputL)
	if err != nil {
		return None, err
	}

	right, err := arithmeticInt(ctx, inputR)
	if err != nil {
		return None, err
	}

	if right == 0 {
		return None, Error(ErrInvalidOperation, "modulo by zero")
	}

	return left % right, nil
}

func incrementNonTemporal(ctx context.Context, input Value) (Value, error) {
	number, err := arithmeticNumber(ctx, input)
	if err != nil {
		return None, err
	}

	switch value := number.(type) {
	case Int:
		return addInts(value, 1)
	case Float:
		return addFloats(value, 1, "increment")
	default:
		return None, unaryOperatorTypeError(operator.Increment, input)
	}
}

func decrementNonTemporal(ctx context.Context, input Value) (Value, error) {
	number, err := arithmeticNumber(ctx, input)
	if err != nil {
		return None, err
	}

	switch value := number.(type) {
	case Int:
		return subtractInts(value, 1)
	case Float:
		return addFloats(value, -1, "decrement")
	default:
		return None, unaryOperatorTypeError(operator.Decrement, input)
	}
}

func arithmeticNumber(ctx context.Context, input Value) (Value, error) {
	switch value := input.(type) {
	case Int, Float:
		return value, nil
	case String:
		if strings.Contains(value.String(), ".") {
			converted, err := ToFloat(ctx, value)
			if err != nil {
				return None, err
			}

			return converted, nil
		}

		converted, err := ToInt(ctx, value)
		if err != nil {
			return None, err
		}

		return converted, nil
	case Iterable:
		return iterableArithmeticNumber(ctx, value)
	default:
		converted, err := ToFloat(ctx, input)
		if err != nil {
			return None, err
		}

		return converted, nil
	}
}

func iterableArithmeticNumber(ctx context.Context, input Iterable) (Value, error) {
	iterator, err := input.Iterate(ctx)
	if err != nil {
		return None, err
	}

	integerTotal := ZeroInt
	floatTotal := ZeroFloat

	for {
		value, _, err := iterator.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return None, err
		}

		number, err := arithmeticNumber(ctx, value)
		if err != nil {
			return None, err
		}

		switch number := number.(type) {
		case Int:
			integerTotal, err = addInts(integerTotal, number)
		case Float:
			floatTotal, err = addFloats(floatTotal, number, "addition")
		}
		if err != nil {
			return None, err
		}
	}

	if floatTotal == 0 {
		return integerTotal, nil
	}

	return addFloats(Float(integerTotal), floatTotal, "addition")
}

func arithmeticInt(ctx context.Context, input Value) (Int, error) {
	switch value := input.(type) {
	case Int:
		return value, nil
	case Float:
		floating := float64(value)
		if math.IsNaN(floating) || math.IsInf(floating, 0) || floating >= math.Exp2(63) || floating < -math.Exp2(63) {
			return ZeroInt, integerRangeError("conversion")
		}

		return Int(value), nil
	case List:
		iterator, err := value.Iterate(ctx)
		if err != nil {
			return ZeroInt, err
		}

		result := ZeroInt
		for {
			item, _, err := iterator.Next(ctx)
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			if err != nil {
				return ZeroInt, err
			}

			integer, err := arithmeticInt(ctx, item)
			if err != nil {
				return ZeroInt, err
			}

			result, err = addInts(result, integer)
			if err != nil {
				return ZeroInt, err
			}
		}
	default:
		return ToInt(ctx, input)
	}
}

func addInts(left, right Int) (Int, error) {
	result := left + right
	if ((left ^ result) & (right ^ result)) < 0 {
		return ZeroInt, integerRangeError("addition")
	}

	return result, nil
}

func subtractInts(left, right Int) (Int, error) {
	result := left - right
	if ((left ^ right) & (left ^ result)) < 0 {
		return ZeroInt, integerRangeError("subtraction")
	}

	return result, nil
}

func multiplyInts(left, right Int) (Int, error) {
	if left == 0 || right == 0 {
		return ZeroInt, nil
	}
	if (left == Int(math.MinInt64) && right == -1) || (right == Int(math.MinInt64) && left == -1) {
		return ZeroInt, integerRangeError("multiplication")
	}

	result := left * right
	if result/right != left {
		return ZeroInt, integerRangeError("multiplication")
	}

	return result, nil
}

func negateInt(value Int) (Value, error) {
	if value == Int(math.MinInt64) {
		return None, integerRangeError("negation")
	}

	return -value, nil
}

func negateFloat(value Float) (Value, error) {
	return checkedFloat(-value, "negation")
}

func addFloats(left, right Float, operation string) (Float, error) {
	return checkedFloat(left+right, operation)
}

func subtractNumbers(left, right Value) (Value, error) {
	switch left := left.(type) {
	case Int:
		switch right := right.(type) {
		case Int:
			return subtractInts(left, right)
		case Float:
			return checkedFloat(Float(left)-right, "subtraction")
		}
	case Float:
		switch right := right.(type) {
		case Int:
			return checkedFloat(left-Float(right), "subtraction")
		case Float:
			return checkedFloat(left-right, "subtraction")
		}
	}

	return None, Error(ErrInvalidOperation, "invalid numeric subtraction")
}

func multiplyNumbers(left, right Value) (Value, error) {
	switch left := left.(type) {
	case Int:
		switch right := right.(type) {
		case Int:
			return multiplyInts(left, right)
		case Float:
			return checkedFloat(Float(left)*right, "multiplication")
		}
	case Float:
		switch right := right.(type) {
		case Int:
			return checkedFloat(left*Float(right), "multiplication")
		case Float:
			return checkedFloat(left*right, "multiplication")
		}
	}

	return None, Error(ErrInvalidOperation, "invalid numeric multiplication")
}

func divideNumbers(left, right Value) (Value, error) {
	switch left := left.(type) {
	case Int:
		switch right := right.(type) {
		case Int:
			if right == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}
			if left == Int(math.MinInt64) && right == -1 {
				return None, integerRangeError("division")
			}
			if left%right != 0 {
				return checkedFloat(Float(left)/Float(right), "division")
			}

			return left / right, nil
		case Float:
			if right == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}

			return checkedFloat(Float(left)/right, "division")
		}
	case Float:
		switch right := right.(type) {
		case Int:
			if right == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}

			return checkedFloat(left/Float(right), "division")
		case Float:
			if right == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}

			return checkedFloat(left/right, "division")
		}
	}

	return None, Error(ErrInvalidOperation, "invalid numeric division")
}

func checkedFloat(value Float, operation string) (Float, error) {
	result := float64(value)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return ZeroFloat, Errorf(ErrRange, "numeric %s produced a non-finite result", operation)
	}

	return value, nil
}

func integerRangeError(operation string) error {
	return Errorf(ErrRange, "integer %s exceeds the supported range", operation)
}
