package runtime

import "math"

type nativeNumber struct {
	integer  Int
	floating Float
	isFloat  bool
}

func classifyNativeNumber(value Value) (nativeNumber, bool) {
	switch value := value.(type) {
	case Int:
		return nativeNumber{integer: value}, true
	case Float:
		return nativeNumber{floating: value, isFloat: true}, true
	default:
		return nativeNumber{}, false
	}
}

func nativeFloat(value nativeNumber) Float {
	if value.isFloat {
		return value.floating
	}

	return Float(value.integer)
}

func concatValues(left, right Value) Value {
	return String(left.String() + right.String())
}

func addNumbers(left, right nativeNumber) (Value, error) {
	if !left.isFloat && !right.isFloat {
		return addInts(left.integer, right.integer)
	}

	return addFloats(nativeFloat(left), nativeFloat(right), "addition")
}

func subtractNumbers(left, right nativeNumber) (Value, error) {
	if !left.isFloat && !right.isFloat {
		return subtractInts(left.integer, right.integer)
	}

	return checkedFloat(nativeFloat(left)-nativeFloat(right), "subtraction")
}

func multiplyNumbers(left, right nativeNumber) (Value, error) {
	if !left.isFloat && !right.isFloat {
		return multiplyInts(left.integer, right.integer)
	}

	return checkedFloat(nativeFloat(left)*nativeFloat(right), "multiplication")
}

func divideNumbers(left, right nativeNumber) (Value, error) {
	if !left.isFloat && !right.isFloat {
		if right.integer == 0 {
			return None, divisionByZeroError()
		}

		if left.integer == Int(math.MinInt64) && right.integer == -1 {
			return None, integerRangeError("division")
		}

		if left.integer%right.integer == 0 {
			return left.integer / right.integer, nil
		}
	}

	rightFloat := nativeFloat(right)
	if rightFloat == 0 {
		return None, divisionByZeroError()
	}

	return checkedFloat(nativeFloat(left)/rightFloat, "division")
}

func moduloNumbers(left, right nativeNumber) (Value, error) {
	if !left.isFloat && !right.isFloat {
		if right.integer == 0 {
			return None, moduloByZeroError()
		}

		return left.integer % right.integer, nil
	}

	rightFloat := nativeFloat(right)
	if rightFloat == 0 {
		return None, moduloByZeroError()
	}

	result := Float(math.Mod(float64(nativeFloat(left)), float64(rightFloat)))

	return checkedFloat(result, "modulo")
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
