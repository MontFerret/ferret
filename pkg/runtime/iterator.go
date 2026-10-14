package runtime

import (
	"context"
	"errors"
	"io"
)

type (
	// Iterable represents an interface of a value that can be iterated by using an
	// iterator. Implementations that may block must observe ctx while creating it.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of an iterator.
	// Next must return io.EOF when the iterator is exhausted and observe ctx while
	// waiting for the next value.
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
		if closeErr := closable.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
	}

	return err
}

func ForEachIter(ctx context.Context, iter Iterator, predicate Predicate) error {
	for {
		val, key, err := iter.Next(ctx)
		if err != nil {
			if iteratorEOF(err) {
				return nil
			}

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

// iteratorEOF accepts wrapped exhaustion only when every underlying cause is EOF.
// errors.Is alone would discard a failure joined with normal exhaustion.
func iteratorEOF(err error) bool {
	if err == io.EOF {
		return true
	}
	switch cause := err.(type) {
	case interface{ Unwrap() []error }:
		children := cause.Unwrap()
		if len(children) == 0 {
			return false
		}
		for _, child := range children {
			if !iteratorEOF(child) {
				return false
			}
		}
		return true
	case interface{ Unwrap() error }:
		return iteratorEOF(cause.Unwrap())
	default:
		return errors.Is(err, io.EOF)
	}
}
