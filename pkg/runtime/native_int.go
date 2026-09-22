package runtime

import "math"

// ToNativeInt returns the exact native integer when value fits the host's int
// range, or zero and false otherwise. Negative values are allowed; callers
// remain responsible for domain restrictions and error reporting.
func ToNativeInt(value Int) (int, bool) {
	if value < math.MinInt || value > math.MaxInt {
		return 0, false
	}

	return int(value), true
}

// CapacityHint returns length*multiplier+extra when it fits a native int.
// It returns zero if length or extra is negative, multiplier is nonpositive,
// or the conversion or arithmetic would overflow. Zero omits the optional hint;
// it does not validate a mandatory size or guarantee allocation success.
// CapacityHint performs no allocation.
func CapacityHint(length Int, multiplier, extra int) int {
	size, ok := ToNativeInt(length)
	if !ok || size < 0 || multiplier <= 0 || extra < 0 {
		return 0
	}

	if size > (math.MaxInt-extra)/multiplier {
		return 0
	}

	return size*multiplier + extra
}
