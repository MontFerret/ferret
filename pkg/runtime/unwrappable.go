package runtime

// Unwrappable represents an interface that can be unwrapped to get the underlying value.
type Unwrappable interface {
	Unwrap() any
}

// Unwrap attempts to unwrap the given value if it implements the Unwrappable interface.
// If the value does not implement Unwrappable, it returns nil.
func Unwrap(val any) any {
	if unwrappable, ok := val.(Unwrappable); ok {
		return unwrappable.Unwrap()
	}

	return nil
}

// UnwrapAs attempts to unwrap the given value if it implements the Unwrappable interface and checks if the unwrapped value is of type T.
// If the value does not implement Unwrappable or the unwrapped value is not of type T, it returns the zero value of T and false.
func UnwrapAs[T any](val any) (T, bool) {
	if unwrappable, ok := val.(Unwrappable); ok {
		unwrapped := unwrappable.Unwrap()

		if result, ok := unwrapped.(T); ok {
			return result, true
		}
	}

	var zero T

	return zero, false
}
