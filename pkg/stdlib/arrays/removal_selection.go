package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func removeLimited(ctx context.Context, list runtime.List, target runtime.Value, limit runtime.Int) (runtime.Value, error) {
	var keep runtime.IndexReadablePredicate
	if limit < 0 {
		keep = removalPredicate(target)
	} else {
		keep = limitedRemovalPredicate(target, limit)
	}

	return list.Filter(ctx, keep)
}

// Unlimited removal does not allocate or update legacy limit bookkeeping.
func removalPredicate(target runtime.Value) runtime.IndexReadablePredicate {
	return func(ctx context.Context, item runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		equal, err := runtime.EqualValues(ctx, item, target)
		if err != nil {
			return false, err
		}

		return !equal, nil
	}
}

func limitedRemovalPredicate(target runtime.Value, limit runtime.Int) runtime.IndexReadablePredicate {
	var counter runtime.Int

	return func(ctx context.Context, item runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		// Even zero or exhausted limits must preserve fallible host comparisons.
		equal, err := runtime.EqualValues(ctx, item, target)
		if err != nil {
			return false, err
		}

		if equal {
			counter++

			if limit == 0 {
				return true, nil
			}

			if counter <= limit {
				return false, nil
			}
		}

		return true, nil
	}
}

// Compact survivors before trimming the tail to avoid quadratic middle removals.
// A failed comparison or write can leave earlier writes applied; the list is borrowed.
func removeMatchingInPlace(ctx context.Context, list runtime.List, target runtime.Value) error {
	size, err := mutableListLength(ctx, list)
	if err != nil {
		return err
	}

	keep := removalPredicate(target)
	write := runtime.ZeroInt
	for read := runtime.ZeroInt; read < size; read++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		value, err := list.At(ctx, read)
		if err != nil {
			return err
		}

		retained, err := keep(ctx, value, read)
		if err != nil {
			return err
		}

		if !retained {
			continue
		}

		if write != read {
			if err := list.SetAt(ctx, write, value); err != nil {
				return err
			}
		}

		write++
	}

	for end := size; end > write; end-- {
		if err := ctx.Err(); err != nil {
			return err
		}

		if _, err := list.RemoveAt(ctx, end-1); err != nil {
			return err
		}
	}

	return nil
}
