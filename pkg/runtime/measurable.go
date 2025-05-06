package runtime

import "context"

type Measurable interface {
	Length(ctx context.Context) (Int, error)
}

func Length(ctx context.Context, value Value) (Int, error) {
	c, ok := value.(Measurable)

	if !ok {
		return 0, nil
	}

	return c.Length(ctx)
}

func IsEmpty(ctx context.Context, value Value) (bool, error) {
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
