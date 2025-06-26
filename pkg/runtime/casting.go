package runtime

func CastBoolean(input Value) (Boolean, error) {
	boolean, ok := input.(Boolean)

	if ok {
		return boolean, nil
	}

	return False, TypeErrorOf(input, TypeBoolean)
}

func SafeCastBoolean(input Value, fallback Boolean) Boolean {
	boolean, ok := input.(Boolean)

	if ok {
		return boolean
	}

	return fallback
}

func CastInt(input Value) (Int, error) {
	integer, ok := input.(Int)

	if ok {
		return integer, nil
	}

	return ZeroInt, TypeErrorOf(input, TypeInt)
}

func SafeCastInt(input Value, fallback Int) Int {
	integer, ok := input.(Int)

	if ok {
		return integer
	}

	return fallback
}

func CastFloat(input Value) (Float, error) {
	float, ok := input.(Float)

	if ok {
		return float, nil
	}

	return ZeroFloat, TypeErrorOf(input, TypeFloat)
}

func SafeCastFloat(input Value, fallback Float) Float {
	float, ok := input.(Float)

	if ok {
		return float
	}

	return fallback
}

func CastString(input Value) (String, error) {
	str, ok := input.(String)

	if ok {
		return str, nil
	}

	return EmptyString, TypeErrorOf(input, TypeString)
}

func SafeCastString(input Value, fallback String) String {
	str, ok := input.(String)

	if ok {
		return str
	}

	return fallback
}

func CastDateTime(input Value) (DateTime, error) {
	dt, ok := input.(DateTime)

	if ok {
		return dt, nil
	}

	return ZeroDateTime, TypeErrorOf(input, TypeDateTime)
}

func SafeCastDateTime(input Value, fallback DateTime) DateTime {
	dt, ok := input.(DateTime)

	if ok {
		return dt
	}

	return fallback
}

func CastList(input Value) (List, error) {
	arr, ok := input.(List)

	if ok {
		return arr, nil
	}

	return nil, TypeErrorOf(input, TypeList)
}

func SafeCastList(input Value, fallback List) List {
	arr, ok := input.(List)

	if ok {
		return arr
	}

	return fallback
}

func CastMap(input Value) (Map, error) {
	obj, ok := input.(Map)

	if ok {
		return obj, nil
	}

	return nil, TypeErrorOf(input, TypeMap)
}

func SafeCastMap(input Value, fallback Map) Map {
	obj, ok := input.(Map)

	if ok {
		return obj
	}
	return fallback

}

func CastBinary(input Value) (Binary, error) {
	bin, ok := input.(Binary)

	if ok {
		return bin, nil
	}

	return nil, TypeErrorOf(input, TypeBinary)
}

func SafeCastBinary(input Value, fallback Binary) Binary {
	bin, ok := input.(Binary)

	if ok {
		return bin
	}

	return fallback
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
