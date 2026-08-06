package runtime

import "math"

const (
	minInt64Float          = Float(-1 << 63)
	maxInt64ExclusiveFloat = Float(1 << 63)
)

func compareFloatValues(left, right Float) Ordering {
	leftNaN := math.IsNaN(float64(left))
	rightNaN := math.IsNaN(float64(right))

	if leftNaN {
		if rightNaN {
			return Equal
		}

		return Greater
	}

	if rightNaN {
		return Less
	}

	return compareOrdered(left, right)
}

// compareIntFloatValues avoids converting the integer to Float, which would
// lose information outside float64's exact integer range.
func compareIntFloatValues(integer Int, floating Float) Ordering {
	if math.IsNaN(float64(floating)) {
		return Less
	}

	if floating >= maxInt64ExclusiveFloat {
		return Less
	}

	if floating < minInt64Float {
		return Greater
	}

	truncated := Int(floating)
	if result := compareOrdered(integer, truncated); result != Equal {
		return result
	}

	return compareOrdered(Float(truncated), floating)
}

func exactIntFromFloat(value Float) (Int, bool) {
	if math.IsNaN(float64(value)) || value < minInt64Float || value >= maxInt64ExclusiveFloat {
		return ZeroInt, false
	}

	integer := Int(value)

	return integer, Float(integer) == value
}
