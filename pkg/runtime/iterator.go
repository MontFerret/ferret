package runtime

import (
	"context"
	"errors"
	"io"
)

type (
	// Iterable represents an interface of a value that can be iterated by using an iterator.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of an iterator.
	// Next must return io.EOF when the iterator is exhausted.
	Iterator interface {
		Next(ctx context.Context) (value Value, key Value, err error)
	}
)

func ForEach(ctx context.Context, input Iterable, predicate Predicate) error {
	iter, err := input.Iterate(ctx)

	if err != nil {
		return err
	}

	err = ForEachIter(ctx, iter, predicate)
	closable, ok := iter.(io.Closer)

	if ok {
		if err := closable.Close(); err != nil {
			return err
		}
	}

	return err
}

func ForEachIter(ctx context.Context, iter Iterator, predicate Predicate) error {
	for {
		val, key, err := iter.Next(ctx)

		if errors.Is(err, io.EOF) {
			return nil
		}

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
