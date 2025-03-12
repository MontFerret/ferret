package core

import "context"

type (
	// Iterable represents an interface of a value that can be iterated by using an iterator.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of an iterator.
	Iterator interface {
		HasNext(ctx context.Context) (bool, error)
		Next(ctx context.Context) (value Value, key Value, err error)
	}
)

type listIterator struct {
	items List
	pos   Int
}

func NewListIterator(list List) Iterator {
	return &listIterator{items: list}
}

func (l *listIterator) HasNext(ctx context.Context) (bool, error) {
	length, err := l.items.Length(ctx)

	if err != nil {
		return false, err
	}

	return l.pos < length, nil
}

func (l *listIterator) Next(ctx context.Context) (value Value, key Value, err error) {
	idx := l.pos
	l.pos++

	value, err = l.items.Get(ctx, idx)

	if err != nil {
		return None, None, err
	}

	return value, idx, nil
}

type mapIterator struct {
	items Map
}
