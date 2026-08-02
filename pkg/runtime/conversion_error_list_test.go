package runtime

import "context"

type conversionErrorList struct {
	*Array
	lengthErr  error
	iterateErr error
}

func newConversionErrorList(lengthErr, iterateErr error, values ...Value) *conversionErrorList {
	return &conversionErrorList{
		Array:      NewArrayWith(values...),
		lengthErr:  lengthErr,
		iterateErr: iterateErr,
	}
}

func (l *conversionErrorList) Length(ctx context.Context) (Int, error) {
	if l.lengthErr != nil {
		return 0, l.lengthErr
	}

	return l.Array.Length(ctx)
}

func (l *conversionErrorList) Iterate(ctx context.Context) (Iterator, error) {
	if l.iterateErr != nil {
		return nil, l.iterateErr
	}

	return l.Array.Iterate(ctx)
}
