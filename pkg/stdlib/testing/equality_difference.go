package testing

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	differenceKind uint8

	difference struct {
		expected runtime.Value
		actual   runtime.Value
		path     string
		kind     differenceKind
	}

	differenceResult struct {
		items []difference
		total int
	}

	diagnosticObject struct {
		values map[string]runtime.Value
		keys   []string
	}
)

const (
	differenceValues differenceKind = iota
	differenceExpectedMissing
	differenceActualMissing
)

const (
	// maxReportedDifferences keeps one assertion failure readable while the
	// complete traversal still provides an exact omitted count.
	maxReportedDifferences = 10
	// maxDifferenceDepth prevents pathological nesting from exhausting the Go
	// stack; the unequal subtree at the boundary is still reported generically.
	maxDifferenceDepth = 64
)

func discoverStructuralDifferences(
	ctx context.Context,
	actual runtime.Value,
	expected runtime.Value,
) (differenceResult, bool) {
	expectedObject, expectedObjectOK := diagnosticObjectOf(ctx, expected)
	actualObject, actualObjectOK := diagnosticObjectOf(ctx, actual)
	if expectedObjectOK && actualObjectOK {
		result := differenceResult{items: make([]difference, 0, maxReportedDifferences)}
		if !diffObjects(ctx, &result, "$", expectedObject, actualObject, 0) || result.total == 0 {
			return differenceResult{}, false
		}

		return result, true
	}

	expectedList, expectedListOK := diagnosticListOf(ctx, expected)
	actualList, actualListOK := diagnosticListOf(ctx, actual)
	if expectedListOK && actualListOK {
		result := differenceResult{items: make([]difference, 0, maxReportedDifferences)}
		if !diffLists(ctx, &result, "$", expectedList, actualList, 0) || result.total == 0 {
			return differenceResult{}, false
		}

		return result, true
	}

	return differenceResult{}, false
}

func walkDifference(
	ctx context.Context,
	result *differenceResult,
	path string,
	expected runtime.Value,
	actual runtime.Value,
	depth int,
) bool {
	equal, err := runtime.EqualValues(ctx, expected, actual)
	if err != nil {
		return false
	}

	if equal {
		return true
	}

	if depth >= maxDifferenceDepth {
		addDifference(result, difference{path: path, expected: expected, actual: actual})

		return true
	}

	expectedObject, expectedObjectOK := diagnosticObjectOf(ctx, expected)
	actualObject, actualObjectOK := diagnosticObjectOf(ctx, actual)
	if expectedObjectOK && actualObjectOK {
		return diffObjects(ctx, result, path, expectedObject, actualObject, depth)
	}

	expectedList, expectedListOK := diagnosticListOf(ctx, expected)
	actualList, actualListOK := diagnosticListOf(ctx, actual)
	if expectedListOK && actualListOK {
		return diffLists(ctx, result, path, expectedList, actualList, depth)
	}

	addDifference(result, difference{path: path, expected: expected, actual: actual})

	return true
}

func diffObjects(
	ctx context.Context,
	result *differenceResult,
	path string,
	expected diagnosticObject,
	actual diagnosticObject,
	depth int,
) bool {
	expectedIndex := 0
	actualIndex := 0

	for expectedIndex < len(expected.keys) || actualIndex < len(actual.keys) {
		switch {
		case expectedIndex >= len(expected.keys):
			key := actual.keys[actualIndex]
			addDifference(result, difference{
				path:   appendObjectPath(path, key),
				actual: actual.values[key],
				kind:   differenceExpectedMissing,
			})
			actualIndex++
		case actualIndex >= len(actual.keys):
			key := expected.keys[expectedIndex]
			addDifference(result, difference{
				path:     appendObjectPath(path, key),
				expected: expected.values[key],
				kind:     differenceActualMissing,
			})
			expectedIndex++
		default:
			expectedKey := expected.keys[expectedIndex]
			actualKey := actual.keys[actualIndex]

			switch {
			case expectedKey < actualKey:
				addDifference(result, difference{
					path:     appendObjectPath(path, expectedKey),
					expected: expected.values[expectedKey],
					kind:     differenceActualMissing,
				})
				expectedIndex++
			case expectedKey > actualKey:
				addDifference(result, difference{
					path:   appendObjectPath(path, actualKey),
					actual: actual.values[actualKey],
					kind:   differenceExpectedMissing,
				})
				actualIndex++
			default:
				if !walkDifference(
					ctx,
					result,
					appendObjectPath(path, expectedKey),
					expected.values[expectedKey],
					actual.values[actualKey],
					depth+1,
				) {
					return false
				}

				expectedIndex++
				actualIndex++
			}
		}
	}

	return true
}

func diffLists(
	ctx context.Context,
	result *differenceResult,
	path string,
	expected *runtime.Array,
	actual *runtime.Array,
	depth int,
) bool {
	expectedLength, err := expected.Length(ctx)
	if err != nil {
		return false
	}

	actualLength, err := actual.Length(ctx)
	if err != nil {
		return false
	}

	commonLength := expectedLength
	if actualLength < commonLength {
		commonLength = actualLength
	}

	for index := runtime.ZeroInt; index < commonLength; index++ {
		expectedValue, expectedFound, err := expected.LookupAt(ctx, index)
		if err != nil || !expectedFound {
			return false
		}

		actualValue, actualFound, err := actual.LookupAt(ctx, index)
		if err != nil || !actualFound {
			return false
		}

		if !walkDifference(
			ctx,
			result,
			appendIndexPath(path, index),
			expectedValue,
			actualValue,
			depth+1,
		) {
			return false
		}
	}

	for index := commonLength; index < expectedLength; index++ {
		value, found, err := expected.LookupAt(ctx, index)
		if err != nil || !found {
			return false
		}

		addDifference(result, difference{
			path:     appendIndexPath(path, index),
			expected: value,
			kind:     differenceActualMissing,
		})
	}

	for index := commonLength; index < actualLength; index++ {
		value, found, err := actual.LookupAt(ctx, index)
		if err != nil || !found {
			return false
		}

		addDifference(result, difference{
			path:   appendIndexPath(path, index),
			actual: value,
			kind:   differenceExpectedMissing,
		})
	}

	return true
}

func diagnosticObjectOf(ctx context.Context, value runtime.Value) (diagnosticObject, bool) {
	// ObjectLike identifies Ferret's materialized object implementations. Other
	// maps must explicitly opt into a native snapshot so diagnostics do not walk
	// arbitrary remote or lazy collections.
	if object, ok := value.(runtime.ObjectLike); ok {
		return snapshotDiagnosticObject(ctx, object)
	}

	if _, ok := value.(runtime.Map); !ok {
		return diagnosticObject{}, false
	}

	snapshotter, ok := value.(runtime.MapSnapshotter)
	if !ok {
		return diagnosticObject{}, false
	}

	snapshot, err := snapshotter.Snapshot(ctx)
	if err != nil || snapshot == nil {
		return diagnosticObject{}, false
	}

	return snapshotDiagnosticObject(ctx, snapshot)
}

func snapshotDiagnosticObject(ctx context.Context, object runtime.ObjectLike) (diagnosticObject, bool) {
	length, err := object.Length(ctx)
	if err != nil {
		return diagnosticObject{}, false
	}

	result := diagnosticObject{
		values: make(map[string]runtime.Value, length),
		keys:   make([]string, 0, length),
	}
	err = object.ForEach(ctx, func(_ context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		stringKey, ok := key.(runtime.String)
		if !ok {
			return runtime.False, runtime.TypeErrorOf(key, runtime.TypeString)
		}

		keyText := string(stringKey)
		if _, exists := result.values[keyText]; !exists {
			result.keys = append(result.keys, keyText)
		}
		result.values[keyText] = value

		return runtime.True, nil
	})
	if err != nil {
		return diagnosticObject{}, false
	}

	sort.Strings(result.keys)

	return result, true
}

func diagnosticListOf(ctx context.Context, value runtime.Value) (*runtime.Array, bool) {
	// Native arrays are already materialized. Other lists require the explicit
	// snapshot capability for the same side-effect boundary as maps above.
	if array, ok := value.(*runtime.Array); ok {
		return array, true
	}

	if _, ok := value.(runtime.List); !ok {
		return nil, false
	}

	snapshotter, ok := value.(runtime.ListSnapshotter)
	if !ok {
		return nil, false
	}

	snapshot, err := snapshotter.Snapshot(ctx)
	if err != nil || snapshot == nil {
		return nil, false
	}

	return snapshot, true
}

func addDifference(result *differenceResult, item difference) {
	result.total++
	if len(result.items) < maxReportedDifferences {
		result.items = append(result.items, item)
	}
}
