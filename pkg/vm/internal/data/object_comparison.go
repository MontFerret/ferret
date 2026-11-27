package data

import (
	"context"
	"slices"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type objectComparisonSnapshot struct {
	fast   *FastObject
	values map[runtime.String]runtime.Value
	keys   []runtime.String
}

func equalObjectLike(ctx context.Context, left, right runtime.ObjectLike) (bool, error) {
	leftSize, err := left.Length(ctx)
	if err != nil {
		return false, err
	}

	rightSize, err := right.Length(ctx)
	if err != nil {
		return false, err
	}

	if leftSize != rightSize {
		return false, nil
	}

	leftSnapshot, err := snapshotObjectLike(ctx, left, leftSize)
	if err != nil {
		return false, err
	}

	rightSnapshot, err := snapshotObjectLike(ctx, right, rightSize)
	if err != nil {
		return false, err
	}

	for idx, leftKey := range leftSnapshot.keys {
		rightKey := rightSnapshot.keys[idx]
		if leftKey != rightKey {
			return false, nil
		}

		equal, err := runtime.EqualValues(
			ctx,
			objectSnapshotValue(leftSnapshot, idx),
			objectSnapshotValue(rightSnapshot, idx),
		)
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func compareObjectLike(ctx context.Context, left, right runtime.ObjectLike) (runtime.Ordering, error) {
	leftSize, err := left.Length(ctx)
	if err != nil {
		return runtime.Equal, err
	}

	rightSize, err := right.Length(ctx)
	if err != nil {
		return runtime.Equal, err
	}

	switch {
	case leftSize < rightSize:
		return runtime.Less, nil
	case leftSize > rightSize:
		return runtime.Greater, nil
	}

	leftSnapshot, err := snapshotObjectLike(ctx, left, leftSize)
	if err != nil {
		return runtime.Equal, err
	}

	rightSnapshot, err := snapshotObjectLike(ctx, right, rightSize)
	if err != nil {
		return runtime.Equal, err
	}

	for idx, leftKey := range leftSnapshot.keys {
		rightKey := rightSnapshot.keys[idx]
		if leftKey != rightKey {
			// Object ordering historically reverses lexical key order.
			if leftKey < rightKey {
				return runtime.Greater, nil
			}

			return runtime.Less, nil
		}

		ordering, err := runtime.CompareValues(
			ctx,
			objectSnapshotValue(leftSnapshot, idx),
			objectSnapshotValue(rightSnapshot, idx),
		)
		if err != nil {
			return runtime.Equal, err
		}

		if ordering != runtime.Equal {
			return ordering, nil
		}
	}

	return runtime.Equal, nil
}

func snapshotObjectLike(ctx context.Context, object runtime.ObjectLike, size runtime.Int) (objectComparisonSnapshot, error) {
	if fastObject, ok := object.(*FastObject); ok {
		keys := fastObject.keys()
		slices.Sort(keys)

		return objectComparisonSnapshot{fast: fastObject, keys: keys}, nil
	}

	capacity := runtime.CapacityHint(size, 1, 0)
	keys := make([]runtime.String, 0, capacity)
	values := make(map[runtime.String]runtime.Value, capacity)
	err := object.ForEach(ctx, func(ctx context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		name, ok := key.(runtime.String)
		if !ok {
			return runtime.False, runtime.TypeErrorOf(key, runtime.TypeString)
		}

		keys = append(keys, name)
		values[name] = value

		return runtime.True, nil
	})
	if err != nil {
		return objectComparisonSnapshot{}, err
	}

	slices.Sort(keys)

	return objectComparisonSnapshot{keys: keys, values: values}, nil
}

func objectSnapshotValue(snapshot objectComparisonSnapshot, idx int) runtime.Value {
	if snapshot.fast != nil {
		return snapshot.fast.getByKey(&snapshot.keys[idx])
	}

	return snapshot.values[snapshot.keys[idx]]
}
