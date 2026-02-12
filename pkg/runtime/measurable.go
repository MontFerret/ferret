package runtime

// Measurable is an interface for values that have a length or size, such as strings, arrays, or objects.
type Measurable interface {
	// Length returns the length or size of the value.
	// For strings, it returns the number of characters;
	// for arrays, it returns the number of elements;
	// for objects, it returns the number of properties.
	Length(ctx Context) (Int, error)
}

// Length returns the length or size of the given value if it implements the Measurable interface.
func Length(ctx Context, value Value) (Int, error) {
	c, ok := value.(Measurable)

	if !ok {
		return 0, nil
	}

	return c.Length(ctx)
}

// IsEmpty checks if the given value is empty (i.e., has a length of zero) if it implements the Measurable interface.
func IsEmpty(ctx Context, value Value) (bool, error) {
	size, err := Length(ctx, value)

	if err != nil {
		return false, err
	}

	intVal, err := ToInt(ctx, size)

	if err != nil {
		return false, err
	}

	return intVal == 0, nil
}
