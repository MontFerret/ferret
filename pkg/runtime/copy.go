package runtime

// Copy returns a shallow copy of src preserving its static Value type T, which
// may be a concrete type or a capability interface such as List. It returns an
// error wrapping ErrInvalidType if src.Copy() produces a value incompatible with T.
func Copy[T Value](src T) (T, error) {
	copied := src.Copy()
	result, err := Cast[T](copied)
	if err != nil {
		var zero T

		return zero, Errorf(err, "copy of %T returned %T", src, copied)
	}

	return result, nil
}
