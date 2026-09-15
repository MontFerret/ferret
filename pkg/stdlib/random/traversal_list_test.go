package random

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Benchmarks use direct traversal so iterator key boxing does not obscure
	// the collection algorithm's allocation cost.
	traversalValues struct {
		runtime.List
		values []runtime.Value
	}

	// Only ForEach is supported. Length, Iterate, Copy, New, indexing, sorting,
	// and mutation fail through the nil embedded interface if called.
	traversalList struct {
		runtime.List
		visit        func(context.Context, int) error
		closeErr     error
		traversalErr error
		values       []runtime.Value
		indices      []runtime.Int
		calls        int
		closes       int
		cursorCloses int
	}

	traversalCursor struct {
		list *traversalList
		next int
	}

	// An opaque borrowed resource: even copying, hashing, or formatting it fails.
	borrowedValue struct {
		runtime.Value
		closes int
	}
)

func (list *traversalValues) ForEach(ctx context.Context, predicate runtime.IndexReadablePredicate) error {
	for i, value := range list.values {
		more, err := predicate(ctx, value, runtime.Int(i))
		if err != nil {
			return err
		}

		if !more {
			return nil
		}
	}

	return nil
}

func (list *traversalList) ForEach(ctx context.Context, predicate runtime.IndexReadablePredicate) error {
	list.calls++
	cursor := &traversalCursor{list: list}
	list.traversalErr = runtime.ForEach(ctx, cursor, func(c context.Context, value, index runtime.Value) (runtime.Boolean, error) {
		return predicate(c, value, index.(runtime.Int))
	})

	return list.traversalErr
}

func (list *traversalList) Close() error {
	list.closes++

	return nil
}

func (cursor *traversalCursor) Iterate(context.Context) (runtime.Iterator, error) {
	return cursor, nil
}

func (cursor *traversalCursor) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if cursor.list.visit != nil {
		if err := cursor.list.visit(ctx, cursor.next); err != nil {
			return nil, nil, err
		}
	}

	if cursor.next == len(cursor.list.values) {
		return nil, nil, io.EOF
	}

	index := runtime.Int(cursor.next)
	if cursor.list.indices != nil {
		index = cursor.list.indices[cursor.next]
	}

	value := cursor.list.values[cursor.next]
	cursor.next++

	return value, index, nil
}

func (cursor *traversalCursor) Close() error {
	cursor.list.cursorCloses++

	return cursor.list.closeErr
}

func (value *borrowedValue) Close() error {
	value.closes++

	return nil
}
