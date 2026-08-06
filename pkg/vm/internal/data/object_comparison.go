package data

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type objectComparisonSnapshot struct {
	fast   *FastObject
	values map[string]runtime.Value
	keys   []string
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

	leftSnapshot, err := snapshotObjectLike(ctx, left, int(leftSize))
	if err != nil {
		return false, err
	}

	rightSnapshot, err := snapshotObjectLike(ctx, right, int(rightSize))
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
			objectSnapshotValue(leftSnapshot, leftKey),
			objectSnapshotValue(rightSnapshot, rightKey),
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

	leftSnapshot, err := snapshotObjectLike(ctx, left, int(leftSize))
	if err != nil {
		return runtime.Equal, err
	}

	rightSnapshot, err := snapshotObjectLike(ctx, right, int(rightSize))
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
			objectSnapshotValue(leftSnapshot, leftKey),
			objectSnapshotValue(rightSnapshot, rightKey),
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

func snapshotObjectLike(ctx context.Context, object runtime.ObjectLike, size int) (objectComparisonSnapshot, error) {
	if fastObject, ok := object.(*FastObject); ok {
		keys := fastObject.keys()
		sort.Strings(keys)

		return objectComparisonSnapshot{fast: fastObject, keys: keys}, nil
	}

	keys := make([]string, 0, size)
	values := make(map[string]runtime.Value, size)
	err := object.ForEach(ctx, func(ctx context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		stringKey, ok := key.(runtime.String)
		if !ok {
			return runtime.False, runtime.TypeErrorOf(key, runtime.TypeString)
		}

		keyText := string(stringKey)
		keys = append(keys, keyText)
		values[keyText] = value

		return runtime.True, nil
	})
	if err != nil {
		return objectComparisonSnapshot{}, err
	}

	sort.Strings(keys)

	return objectComparisonSnapshot{keys: keys, values: values}, nil
}

func objectSnapshotValue(snapshot objectComparisonSnapshot, key string) runtime.Value {
	if snapshot.fast != nil {
		return snapshot.fast.getByKey(key)
	}

	return snapshot.values[key]
}
