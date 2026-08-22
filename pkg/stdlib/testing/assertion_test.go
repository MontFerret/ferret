package testing

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAssertionWrapperPreservesSharedBehavior(t *testing.T) {
	t.Parallel()

	descriptor := assertion{
		defaultMessage: func(_ context.Context, _ []runtime.Value) string {
			return "be ready"
		},
		args: assertionArgs{min: 1, max: 2},
		fn: func(_ context.Context, args []runtime.Value) (bool, error) {
			return args[0] == runtime.True, nil
		},
	}

	requireAssertionSuccess(t, descriptor, true, runtime.True)
	requireAssertionSuccess(t, descriptor, false, runtime.False)
	requireAssertionFailure(t, descriptor, true, "assertion error: expected Boolean 'false' to be ready", runtime.False)
	requireAssertionFailure(t, descriptor, false, "assertion error: expected Boolean 'true' not to be ready", runtime.True)
	requireAssertionFailure(t, descriptor, true, "assertion error: caller message", runtime.False, runtime.NewString("caller message"))
	requireAssertionFailure(t, descriptor, false, "assertion error: 42", runtime.True, runtime.NewInt(42))
}

func TestAssertionWrapperPreservesArityAndPredicateErrors(t *testing.T) {
	t.Parallel()

	predicateErr := errors.New("predicate failed")
	descriptor := assertion{
		defaultMessage: func(_ context.Context, _ []runtime.Value) string {
			return "succeed"
		},
		args: assertionArgs{min: 1, max: 2},
		fn: func(_ context.Context, _ []runtime.Value) (bool, error) {
			return false, predicateErr
		},
	}

	fn := descriptor.positive()
	for _, args := range [][]runtime.Value{
		nil,
		{runtime.True, runtime.NewString("message"), runtime.None},
	} {
		out, err := fn(context.Background(), args...)
		if out != runtime.None {
			t.Fatalf("invalid arity output = %v, want None", out)
		}
		if !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
			t.Fatalf("invalid arity error = %v, want ErrInvalidArgumentNumber", err)
		}
	}

	out, err := fn(context.Background(), runtime.True)
	if out != runtime.None {
		t.Fatalf("predicate error output = %v, want None", out)
	}
	if !errors.Is(err, predicateErr) {
		t.Fatalf("predicate error = %v, want predicate identity", err)
	}
	if errors.Is(err, errAssertion) {
		t.Fatalf("predicate error %v was replaced with an assertion error", err)
	}
}

func TestAssertionDefaultFallbacksRemainStable(t *testing.T) {
	t.Parallel()

	descriptor := assertion{
		args: assertionArgs{min: 1, max: 3},
		fn: func(_ context.Context, _ []runtime.Value) (bool, error) {
			return false, nil
		},
	}

	requireAssertionFailure(
		t,
		descriptor,
		true,
		"assertion error: expected String 'actual' to be expected",
		runtime.NewString("actual"),
		runtime.NewString("expected"),
	)
}

func TestFailAssertionPreservesMessages(t *testing.T) {
	t.Parallel()

	requireAssertionFailure(t, failAssertion, true, "assertion error: expected to not fail")
	requireAssertionFailure(t, failAssertion, true, "assertion error: stop now", runtime.NewString("stop now"))
}

func TestFormatValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    runtime.Value
		expected string
	}{
		{name: "plain string", value: runtime.NewString("test"), expected: "String 'test'"},
		{name: "none", value: runtime.None, expected: "None 'none'"},
		{name: "apostrophe and backslash", value: runtime.NewString(`can't\stop`), expected: `String 'can\'t\\stop'`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if actual := formatValue(t.Context(), test.value); actual != test.expected {
				t.Fatalf("formatValue() = %q, want %q", actual, test.expected)
			}
		})
	}
}

func TestFormatValueLimitsLargeValues(t *testing.T) {
	t.Parallel()

	largeArray := runtime.NewArray(1_000)
	for index := 0; index < 1_000; index++ {
		if err := largeArray.Append(t.Context(), runtime.NewInt(index)); err != nil {
			t.Fatal(err)
		}
	}

	largeObjectValues := make(map[string]runtime.Value, maxRenderedObjectProperties+1)
	for index := 0; index <= maxRenderedObjectProperties; index++ {
		largeObjectValues[fmt.Sprintf("key%02d", index)] = runtime.NewInt(index)
	}

	tests := []struct {
		value runtime.Value
		name  string
	}{
		{name: "long Unicode string", value: runtime.NewString(strings.Repeat("界", 1_000))},
		{name: "large array", value: largeArray},
		{name: "large object", value: runtime.NewObjectWith(largeObjectValues)},
		{name: "large binary", value: runtime.NewBinary(make([]byte, 1_000))},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			first := formatValue(t.Context(), test.value)
			second := formatValue(t.Context(), test.value)
			if first != second {
				t.Fatalf("formatValue() is not deterministic:\nfirst:  %q\nsecond: %q", first, second)
			}
			if size := utf8.RuneCountInString(first); size > maxFormattedValueRunes {
				t.Fatalf("formatted value has %d runes, want at most %d: %q", size, maxFormattedValueRunes, first)
			}
			if !strings.Contains(first, truncatedValueMarker) {
				t.Fatalf("formatted value does not identify truncation: %q", first)
			}
		})
	}
}

func requireAssertionSuccess(t *testing.T, descriptor assertion, positive bool, args ...runtime.Value) {
	t.Helper()

	fn := descriptor.negative()
	if positive {
		fn = descriptor.positive()
	}

	out, err := fn(context.Background(), args...)
	if err != nil {
		t.Fatalf("assertion returned unexpected error: %v", err)
	}
	if out != runtime.None {
		t.Fatalf("assertion output = %v, want None", out)
	}
}

func requireAssertionFailure(t *testing.T, descriptor assertion, positive bool, expected string, args ...runtime.Value) {
	t.Helper()

	fn := descriptor.negative()
	if positive {
		fn = descriptor.positive()
	}

	out, err := fn(context.Background(), args...)
	if out != runtime.None {
		t.Fatalf("assertion output = %v, want None", out)
	}
	if err == nil || err.Error() != expected {
		t.Fatalf("assertion error = %v, want %q", err, expected)
	}
	if !errors.Is(err, errAssertion) {
		t.Fatalf("assertion error %v does not preserve ErrAssertion identity", err)
	}
}
