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
