package runtime

import "context"

type conversionErrorIterable struct {
	err error
}

func newConversionErrorIterable(err error) *conversionErrorIterable {
	return &conversionErrorIterable{err: err}
}

func (i *conversionErrorIterable) String() string {
	return "conversionErrorIterable"
}

func (i *conversionErrorIterable) Hash() uint64 {
	return 0
}

func (i *conversionErrorIterable) Copy() Value {
	return &conversionErrorIterable{err: i.err}
}

func (i *conversionErrorIterable) Iterate(context.Context) (Iterator, error) {
	return nil, i.err
}
