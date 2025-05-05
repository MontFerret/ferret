package runtime

import (
	"context"
	"io"
)

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

	listIterator struct {
		items List
		pos   Int
	}
)

func ForEachOf(ctx context.Context, input Iterable, predicate Predicate) error {
	iter, err := input.Iterate(ctx)

	if err != nil {
		return err
	}

	err = ForEach(ctx, iter, predicate)
	closable, ok := iter.(io.Closer)

	if ok {
		if err := closable.Close(); err != nil {
			return err
		}
	}

	return err
}

func ForEach(ctx context.Context, iter Iterator, predicate Predicate) error {
	for {
		hasNext, err := iter.HasNext(ctx)

		if err != nil {
			return err
		}

		if !hasNext {
			return nil
		}

		val, key, err := iter.Next(ctx)

		if err != nil {
			return err
		}

		res, err := predicate(ctx, val, key)

		if err != nil {
			return err
		}

		if !res {
			return nil
		}
	}
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
