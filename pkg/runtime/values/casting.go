package values

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func CastBoolean(input core.Value) (Boolean, error) {
	boolean, ok := input.(Boolean)

	if ok {
		return boolean, nil
	}

	return False, core.TypeError(input, types.Boolean)
}

func SafeCastBoolean(input core.Value, fallback Boolean) Boolean {
	boolean, ok := input.(Boolean)

	if ok {
		return boolean
	}

	return fallback
}

func CastInt(input core.Value) (Int, error) {
	integer, ok := input.(Int)

	if ok {
		return integer, nil
	}

	return ZeroInt, core.TypeError(input, types.Int)
}

func SafeCastInt(input core.Value, fallback Int) Int {
	integer, ok := input.(Int)

	if ok {
		return integer
	}

	return fallback
}

func CastFloat(input core.Value) (Float, error) {
	float, ok := input.(Float)

	if ok {
		return float, nil
	}

	return ZeroFloat, core.TypeError(input, types.Float)
}

func SafeCastFloat(input core.Value, fallback Float) Float {
	float, ok := input.(Float)

	if ok {
		return float
	}

	return fallback
}

func CastString(input core.Value) (String, error) {
	str, ok := input.(String)

	if ok {
		return str, nil
	}

	return EmptyString, core.TypeError(input, types.String)
}

func SafeCastString(input core.Value, fallback String) String {
	str, ok := input.(String)

	if ok {
		return str
	}

	return fallback
}

func CastDateTime(input core.Value) (DateTime, error) {
	dt, ok := input.(DateTime)

	if ok {
		return dt, nil
	}

	return ZeroDateTime, core.TypeError(input, types.DateTime)
}

func SafeCastDateTime(input core.Value, fallback DateTime) DateTime {
	dt, ok := input.(DateTime)

	if ok {
		return dt
	}

	return fallback
}

func CastArray(input core.Value) (*Array, error) {
	arr, ok := input.(*Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(input, types.Array)
}

func SafeCastArray(input core.Value, fallback *Array) *Array {
	arr, ok := input.(*Array)

	if ok {
		return arr
	}

	return fallback
}

func CastObject(input core.Value) (*Object, error) {
	obj, ok := input.(*Object)

	if ok {
		return obj, nil
	}

	return nil, core.TypeError(input, types.Object)
}

func SafeCastObject(input core.Value, fallback *Object) *Object {
	obj, ok := input.(*Object)

	if ok {
		return obj
	}
	return fallback

}

func CastBinary(input core.Value) (Binary, error) {
	bin, ok := input.(Binary)

	if ok {
		return bin, nil
	}

	return nil, core.TypeError(input, types.Binary)
}

func SafeCastBinary(input core.Value, fallback Binary) Binary {
	bin, ok := input.(Binary)

	if ok {
		return bin
	}

	return fallback
}
