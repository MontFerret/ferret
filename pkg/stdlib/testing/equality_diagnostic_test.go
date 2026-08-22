package testing

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestScalarEqualityDiagnostics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		actual   runtime.Value
		expected runtime.Value
		message  string
	}{
		{
			name:     "float",
			actual:   runtime.NewFloat(1.25),
			expected: runtime.NewFloat(1.5),
			message:  "assertion error: values are not equal\nexpected: Float '1.5'\nactual:   Float '1.25'",
		},
		{
			name:     "type mismatch",
			actual:   runtime.NewString("1"),
			expected: runtime.NewInt(1),
			message:  "assertion error: values are not equal\nexpected: Int '1'\nactual:   String '1'",
		},
		{
			name:     "none",
			actual:   runtime.None,
			expected: runtime.NewInt(1),
			message:  "assertion error: values are not equal\nexpected: Int '1'\nactual:   None 'none'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			requireAssertionFailure(t, equalAssertion, true, test.message, test.actual, test.expected)
		})
	}

	requireAssertionSuccess(t, equalAssertion, true, runtime.None, runtime.None)
}

func TestEqualityDiagnosticsPreserveCustomContext(t *testing.T) {
	t.Parallel()

	requireAssertionFailure(
		t,
		equalAssertion,
		true,
		"assertion error: unexpected API response\nvalues are not equal\nexpected: Int '2'\nactual:   Int '1'",
		runtime.NewInt(1),
		runtime.NewInt(2),
		runtime.NewString("unexpected API response"),
	)
	requireAssertionFailure(
		t,
		equalAssertion,
		false,
		"assertion error: response should differ\nexpected values to differ\nboth: Int '1'",
		runtime.NewInt(1),
		runtime.NewInt(1),
		runtime.NewString("response should differ"),
	)
}

func TestObjectEqualityDiagnostics(t *testing.T) {
	t.Parallel()

	expected := runtime.NewObjectWith(map[string]runtime.Value{
		"metadata": runtime.NewObjectWith(map[string]runtime.Value{
			"none":    runtime.None,
			"version": runtime.NewInt(2),
		}),
		"users": runtime.NewArrayWith(
			runtime.NewObjectWith(map[string]runtime.Value{
				"name": runtime.NewString("Alice"),
				"tags": runtime.NewArrayWith(runtime.NewString("admin"), runtime.NewString("active")),
			}),
		),
		"version": runtime.NewInt(2),
	})
	actual := runtime.NewObjectWith(map[string]runtime.Value{
		"metadata": runtime.NewObjectWith(map[string]runtime.Value{
			"debug":   runtime.True,
			"version": runtime.NewInt(3),
		}),
		"users": runtime.NewArrayWith(
			runtime.NewObjectWith(map[string]runtime.Value{
				"name": runtime.NewString("Bob"),
				"tags": runtime.NewArrayWith(runtime.NewString("admin"), runtime.NewString("disabled")),
			}),
		),
		"version": runtime.NewInt(3),
	})

	requireAssertionFailure(
		t,
		equalAssertion,
		true,
		"assertion error: values are not equal\n"+
			"$.metadata.debug\n  expected: <missing>\n  actual:   Boolean 'true'\n"+
			"$.metadata.none\n  expected: None 'none'\n  actual:   <missing>\n"+
			"$.metadata.version\n  expected: Int '2'\n  actual:   Int '3'\n"+
			"$.users[0].name\n  expected: String 'Alice'\n  actual:   String 'Bob'\n"+
			"$.users[0].tags[1]\n  expected: String 'active'\n  actual:   String 'disabled'\n"+
			"$.version\n  expected: Int '2'\n  actual:   Int '3'",
		actual,
		expected,
	)
}

func TestObjectEqualityDiagnosticOrderingAndEscapedPaths(t *testing.T) {
	t.Parallel()

	expected := runtime.NewObjectWith(map[string]runtime.Value{
		`hello "world"`: runtime.NewInt(1),
		"content-type":  runtime.NewInt(2),
		"ordinary_key":  runtime.NewInt(3),
	})
	actual := runtime.NewObjectWith(map[string]runtime.Value{
		`hello "world"`: runtime.NewInt(4),
		"content-type":  runtime.NewInt(5),
		"ordinary_key":  runtime.NewInt(6),
	})

	want := "values are not equal\n" +
		`$["content-type"]` + "\n  expected: Int '2'\n  actual:   Int '5'\n" +
		`$["hello \"world\""]` + "\n  expected: Int '1'\n  actual:   Int '4'\n" +
		"$.ordinary_key\n  expected: Int '3'\n  actual:   Int '6'"

	for run := 0; run < 25; run++ {
		if got := unequalValuesMessage(t.Context(), actual, expected); got != want {
			t.Fatalf("run %d diagnostic =\n%s\nwant:\n%s", run, got, want)
		}
	}
}

func TestArrayEqualityDiagnostics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		actual   *runtime.Array
		expected *runtime.Array
		message  string
	}{
		{
			name:     "nested arrays and objects",
			actual:   runtime.NewArrayWith(runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{"name": runtime.NewString("Bob")}))),
			expected: runtime.NewArrayWith(runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{"name": runtime.NewString("Alice")}))),
			message:  "values are not equal\n$[0][0].name\n  expected: String 'Alice'\n  actual:   String 'Bob'",
		},
		{
			name:     "shorter actual",
			actual:   runtime.NewArrayWith(runtime.NewInt(1)),
			expected: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.None),
			message: "values are not equal\n$[1]\n  expected: Int '2'\n  actual:   <missing>\n" +
				"$[2]\n  expected: None 'none'\n  actual:   <missing>",
		},
		{
			name:     "longer actual",
			actual:   runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
			expected: runtime.NewArrayWith(runtime.NewInt(1)),
			message:  "values are not equal\n$[1]\n  expected: <missing>\n  actual:   Int '2'",
		},
		{
			name:     "multiple indexes",
			actual:   runtime.NewArrayWith(runtime.NewInt(3), runtime.NewInt(2), runtime.NewInt(1)),
			expected: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
			message: "values are not equal\n$[0]\n  expected: Int '1'\n  actual:   Int '3'\n" +
				"$[2]\n  expected: Int '3'\n  actual:   Int '1'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := unequalValuesMessage(t.Context(), test.actual, test.expected); got != test.message {
				t.Fatalf("diagnostic =\n%s\nwant:\n%s", got, test.message)
			}
		})
	}
}

func TestEqualityDifferenceLimit(t *testing.T) {
	t.Parallel()

	expectedValues := make(map[string]runtime.Value, 12)
	actualValues := make(map[string]runtime.Value, 12)
	for index := 0; index < 12; index++ {
		key := fmt.Sprintf("key%02d", index)
		expectedValues[key] = runtime.NewInt(index)
		actualValues[key] = runtime.NewInt(index + 100)
	}

	message := unequalValuesMessage(
		t.Context(),
		runtime.NewObjectWith(actualValues),
		runtime.NewObjectWith(expectedValues),
	)
	if strings.Count(message, "\n$.key") != maxReportedDifferences {
		t.Fatalf("reported difference count = %d, want %d\n%s", strings.Count(message, "\n$.key"), maxReportedDifferences, message)
	}
	if !strings.HasSuffix(message, "... additional differences may be omitted") {
		t.Fatalf("diagnostic does not report the traversal limit:\n%s", message)
	}
	if strings.Contains(message, "$.key10") || strings.Contains(message, "$.key11") {
		t.Fatalf("diagnostic contains differences past the report limit:\n%s", message)
	}
}

func TestEqualityDifferenceTraversalLimit(t *testing.T) {
	t.Parallel()

	assertLimited := func(t *testing.T, result differenceResult, ok bool) {
		t.Helper()

		if !ok || !result.limitReached || len(result.items) != maxReportedDifferences {
			t.Fatalf("limited difference result = %#v, %t", result, ok)
		}
	}

	t.Run("array", func(t *testing.T) {
		t.Parallel()

		visited := 0
		visitedAfterLimit := 0
		size := maxReportedDifferences + 50
		expectedValues := make([]runtime.Value, size)
		actualValues := make([]runtime.Value, size)
		for index := 0; index < size; index++ {
			calls := &visited
			if index >= maxReportedDifferences {
				calls = &visitedAfterLimit
			}

			expectedValues[index] = &diagnosticEqualityProbe{calls: calls}
			actualValues[index] = &diagnosticEqualityProbe{}
		}
		expectedValues[maxReportedDifferences] = &diagnosticEqualityProbe{
			calls: &visitedAfterLimit,
			err:   errors.New("difference past traversal limit was visited"),
		}

		result, ok := discoverStructuralDifferences(
			t.Context(),
			runtime.NewArrayWith(actualValues...),
			runtime.NewArrayWith(expectedValues...),
		)
		assertLimited(t, result, ok)
		if visited != maxReportedDifferences {
			t.Fatalf("visited equality probes = %d, want %d", visited, maxReportedDifferences)
		}
		if visitedAfterLimit != 0 {
			t.Fatalf("visited %d equality probes after the reporting limit", visitedAfterLimit)
		}
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()

		visited := 0
		visitedAfterLimit := 0
		size := maxReportedDifferences + 50
		expectedValues := make(map[string]runtime.Value, size)
		actualValues := make(map[string]runtime.Value, size)
		for index := 0; index < size; index++ {
			calls := &visited
			if index >= maxReportedDifferences {
				calls = &visitedAfterLimit
			}

			key := fmt.Sprintf("key%02d", index)
			expectedValues[key] = &diagnosticEqualityProbe{calls: calls}
			actualValues[key] = &diagnosticEqualityProbe{}
		}
		expectedValues[fmt.Sprintf("key%02d", maxReportedDifferences)] = &diagnosticEqualityProbe{
			calls: &visitedAfterLimit,
			err:   errors.New("difference past traversal limit was visited"),
		}

		result, ok := discoverStructuralDifferences(
			t.Context(),
			runtime.NewObjectWith(actualValues),
			runtime.NewObjectWith(expectedValues),
		)
		assertLimited(t, result, ok)
		if visited != maxReportedDifferences {
			t.Fatalf("visited equality probes = %d, want %d", visited, maxReportedDifferences)
		}
		if visitedAfterLimit != 0 {
			t.Fatalf("visited %d equality probes after the reporting limit", visitedAfterLimit)
		}
	})

	t.Run("nested", func(t *testing.T) {
		t.Parallel()

		visited := 0
		visitedAfterLimit := 0
		laterSiblingVisits := 0
		size := maxReportedDifferences + 10
		expectedValues := make([]runtime.Value, size)
		actualValues := make([]runtime.Value, size-1)
		for index := 0; index < size; index++ {
			calls := &visited
			if index >= maxReportedDifferences {
				calls = &visitedAfterLimit
			}

			expectedValues[index] = &diagnosticEqualityProbe{calls: calls}
			if index < len(actualValues) {
				actualValues[index] = &diagnosticEqualityProbe{}
			}
		}

		expected := runtime.NewObjectWith(map[string]runtime.Value{
			"first": runtime.NewArrayWith(expectedValues...),
			"later": &diagnosticEqualityProbe{
				calls: &laterSiblingVisits,
				err:   errors.New("later sibling was visited"),
			},
		})
		actual := runtime.NewObjectWith(map[string]runtime.Value{
			"first": runtime.NewArrayWith(actualValues...),
			"later": &diagnosticEqualityProbe{},
		})

		result, ok := discoverStructuralDifferences(t.Context(), actual, expected)
		assertLimited(t, result, ok)
		if visited != maxReportedDifferences {
			t.Fatalf("nested equality probes = %d, want %d", visited, maxReportedDifferences)
		}
		if visitedAfterLimit != 0 {
			t.Fatalf("visited %d nested probes after the reporting limit", visitedAfterLimit)
		}
		if laterSiblingVisits != 0 {
			t.Fatalf("visited later sibling %d times after the reporting limit", laterSiblingVisits)
		}
	})
}

func TestEqualityDifferenceDepthLimit(t *testing.T) {
	t.Parallel()

	var expected runtime.Value = runtime.NewInt(1)
	var actual runtime.Value = runtime.NewInt(2)
	for depth := 0; depth < maxDifferenceDepth+1; depth++ {
		expected = runtime.NewArrayWith(expected)
		actual = runtime.NewArrayWith(actual)
	}

	result, ok := discoverStructuralDifferences(t.Context(), actual, expected)
	if !ok || result.limitReached || len(result.items) != 1 {
		t.Fatalf("depth-limited result = %#v, %t", result, ok)
	}
	if got := strings.Count(result.items[0].path, "[0]"); got != maxDifferenceDepth {
		t.Fatalf("depth-limited path has %d indexes, want %d", got, maxDifferenceDepth)
	}
}

func TestEqualityDiagnosticsUseCollectionSnapshots(t *testing.T) {
	t.Parallel()

	list := &snapshotList{
		List:     runtime.NewArrayWith(runtime.NewInt(2)),
		snapshot: runtime.NewArrayWith(runtime.NewInt(2)),
	}
	listMessage := unequalValuesMessage(t.Context(), runtime.NewArrayWith(runtime.NewInt(1)), list)
	if !strings.Contains(listMessage, "$[0]") || list.snapshotCalls != 1 {
		t.Fatalf("list diagnostic = %q, snapshot calls = %d", listMessage, list.snapshotCalls)
	}

	mapValue := &snapshotMap{
		Map: runtime.NewObjectWith(map[string]runtime.Value{
			"value": runtime.NewInt(2),
		}),
		snapshot: runtime.NewObjectWith(map[string]runtime.Value{
			"value": runtime.NewInt(2),
		}),
	}
	mapMessage := unequalValuesMessage(
		t.Context(),
		runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}),
		mapValue,
	)
	if !strings.Contains(mapMessage, "$.value") || mapValue.snapshotCalls != 1 {
		t.Fatalf("map diagnostic = %q, snapshot calls = %d", mapMessage, mapValue.snapshotCalls)
	}
}

func TestEqualityDiagnosticsFallBackForUnsafeCollections(t *testing.T) {
	t.Parallel()

	list := &unsafeList{List: runtime.NewArrayWith(runtime.NewInt(2))}
	requireAssertionFailure(
		t,
		equalAssertion,
		true,
		"assertion error: values are not equal\nexpected: List '[2]'\nactual:   Array '[1]'",
		runtime.NewArrayWith(runtime.NewInt(1)),
		list,
	)
	if list.atCalls != 1 {
		t.Fatalf("unsafe list At calls = %d, want only the authoritative equality lookup", list.atCalls)
	}

	mapValue := &unsafeMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(2)})}
	requireAssertionFailure(
		t,
		equalAssertion,
		true,
		"assertion error: values are not equal\nexpected: Map '{\"value\":2}'\nactual:   Object '{\"value\":1}'",
		runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}),
		mapValue,
	)
	if mapValue.forEachCalls != 0 {
		t.Fatalf("unsafe map ForEach calls = %d, want no diagnostic traversal", mapValue.forEachCalls)
	}
}

func TestEqualityDiagnosticsFallBackWhenSnapshotFails(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		err      error
		snapshot *runtime.Array
		name     string
	}{
		{name: "error", err: errors.New("snapshot failed")},
		{name: "nil"},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			list := &snapshotList{
				List:        runtime.NewArrayWith(runtime.NewInt(2)),
				snapshot:    test.snapshot,
				snapshotErr: test.err,
			}
			message := unequalValuesMessage(t.Context(), runtime.NewArrayWith(runtime.NewInt(1)), list)
			if strings.Contains(message, "$[0]") || !strings.Contains(message, "expected: List '[2]'") {
				t.Fatalf("fallback diagnostic = %q", message)
			}
			if list.snapshotCalls != 1 {
				t.Fatalf("snapshot calls = %d, want 1", list.snapshotCalls)
			}
		})
	}
}

func TestEqualityDiagnosticsIgnoreInspectionEqualityErrors(t *testing.T) {
	t.Parallel()

	operationalErr := errors.New("later equality failed")
	actual := runtime.NewArrayWith(runtime.NewInt(1), &equalityErrorValue{err: operationalErr})
	expected := runtime.NewArrayWith(runtime.NewInt(2), &equalityErrorValue{err: operationalErr})

	message := unequalValuesMessage(t.Context(), actual, expected)
	if strings.Contains(message, "$[0]") || !strings.Contains(message, "expected: Array") {
		t.Fatalf("inspection error did not produce generic fallback:\n%s", message)
	}
}

func TestNegatedEqualityDiagnosticsDoNotDiff(t *testing.T) {
	t.Parallel()

	requireAssertionFailure(
		t,
		equalAssertion,
		false,
		"assertion error: expected values to differ\nfirst:  Int '1'\nsecond: Float '1'",
		runtime.NewInt(1),
		runtime.NewFloat(1),
	)

	snapshot := runtime.NewArrayWith(runtime.NewInt(1))
	list := &snapshotList{List: snapshot, snapshot: snapshot}
	requireAssertionFailure(
		t,
		equalAssertion,
		false,
		"assertion error: expected values to differ\nfirst:  Array '[1]'\nsecond: List '[1]'",
		snapshot,
		list,
	)
	if list.snapshotCalls != 0 {
		t.Fatalf("negated equality requested %d structural snapshots", list.snapshotCalls)
	}
}
