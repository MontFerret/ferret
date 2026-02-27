package runtime

// Cast attempts to cast the input value to the specified type T. If successful, it returns the value and nil error.
// If the cast fails, it returns the zero value of type T and a TypeError indicating the mismatch.
func Cast[T Value](input Value) (T, error) {
	val, ok := input.(T)

	if ok {
		return val, nil
	}

	var zero T

	return zero, TypeErrorOf(input, TypeOf(zero))
}

// CastOr attempts to cast the input value to the specified type T.
// If the cast is successful, it returns the casted value.
// Otherwise, it returns the provided fallback value.
func CastOr[T Value](input Value, fallback T) T {
	val, ok := input.(T)

	if ok {
		return val
	}

	return fallback
}

func CastBoolean(input Value) (Boolean, error) {
	boolean, ok := input.(Boolean)

	if ok {
		return boolean, nil
	}

	return False, TypeErrorOf(input, TypeBoolean)
}

func CastInt(input Value) (Int, error) {
	integer, ok := input.(Int)

	if ok {
		return integer, nil
	}

	return ZeroInt, TypeErrorOf(input, TypeInt)
}

func CastFloat(input Value) (Float, error) {
	float, ok := input.(Float)

	if ok {
		return float, nil
	}

	return ZeroFloat, TypeErrorOf(input, TypeFloat)
}

func CastString(input Value) (String, error) {
	str, ok := input.(String)

	if ok {
		return str, nil
	}

	return EmptyString, TypeErrorOf(input, TypeString)
}

func CastDateTime(input Value) (DateTime, error) {
	dt, ok := input.(DateTime)

	if ok {
		return dt, nil
	}

	return ZeroDateTime, TypeErrorOf(input, TypeDateTime)
}

func CastCollection(input Value) (Collection, error) {
	arr, ok := input.(Collection)

	if ok {
		return arr, nil
	}

	return nil, TypeErrorOf(input, TypeCollection)
}

func CastArray(input Value) (*Array, error) {
	arr, ok := input.(*Array)

	if ok {
		return arr, nil
	}

	return nil, TypeErrorOf(input, TypeArray)
}

func CastList(input Value) (List, error) {
	arr, ok := input.(List)

	if ok {
		return arr, nil
	}

	return nil, TypeErrorOf(input, TypeList)
}

func CastObject(input Value) (*Object, error) {
	obj, ok := input.(*Object)

	if ok {
		return obj, nil
	}

	return nil, TypeErrorOf(input, TypeObject)
}

func CastMap(input Value) (Map, error) {
	obj, ok := input.(Map)

	if ok {
		return obj, nil
	}

	return nil, TypeErrorOf(input, TypeMap)
}

func CastBinary(input Value) (Binary, error) {
	bin, ok := input.(Binary)

	if ok {
		return bin, nil
	}

	return nil, TypeErrorOf(input, TypeBinary)
}

func CastComparable(input Value) (Comparable, error) {
	comp, ok := input.(Comparable)

	if ok {
		return comp, nil
	}

	return nil, TypeErrorOf(input, TypeComparable)
}

func CastCloneable(input Value) (Cloneable, error) {
	cloneable, ok := input.(Cloneable)

	if ok {
		return cloneable, nil
	}

	return nil, TypeErrorOf(input, TypeCloneable)
}
