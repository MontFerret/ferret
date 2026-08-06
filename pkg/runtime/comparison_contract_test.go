package runtime_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestStrictNativeComparisonContract(t *testing.T) {
	instant := time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC)

	equalityTests := []struct {
		left       runtime.Value
		right      runtime.Value
		name       string
		expected   runtime.Boolean
		comparison runtime.Ordering
	}{
		{name: "none", left: runtime.None, right: runtime.None, expected: runtime.True, comparison: runtime.Equal},
		{name: "boolean", left: runtime.False, right: runtime.True, expected: runtime.False, comparison: runtime.Less},
		{name: "numeric cross representation", left: runtime.NewInt(4), right: runtime.NewFloat(4), expected: runtime.True, comparison: runtime.Equal},
		{name: "numeric order", left: runtime.NewFloat(4.5), right: runtime.NewInt(5), expected: runtime.False, comparison: runtime.Less},
		{name: "duration", left: runtime.NewDuration(time.Second), right: runtime.NewDuration(2 * time.Second), expected: runtime.False, comparison: runtime.Less},
		{name: "string", left: runtime.NewString("a"), right: runtime.NewString("b"), expected: runtime.False, comparison: runtime.Less},
		{name: "datetime", left: runtime.NewDateTime(instant), right: runtime.NewDateTime(instant), expected: runtime.True, comparison: runtime.Equal},
		{name: "binary content", left: runtime.NewBinary([]byte{1, 2}), right: runtime.NewBinary([]byte{1, 3}), expected: runtime.False, comparison: runtime.Less},
		{name: "binary length first", left: runtime.NewBinary([]byte{9}), right: runtime.NewBinary([]byte{0, 0}), expected: runtime.False, comparison: runtime.Less},
		{name: "array structure", left: runtime.NewArrayWith(runtime.NewInt(1)), right: runtime.NewArrayWith(runtime.NewInt(1)), expected: runtime.True, comparison: runtime.Equal},
		{name: "array length first", left: runtime.NewArrayWith(runtime.NewInt(9)), right: runtime.NewArrayWith(runtime.NewInt(0), runtime.NewInt(0)), expected: runtime.False, comparison: runtime.Less},
		{name: "object structure", left: runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewInt(1)}), right: runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewInt(1)}), expected: runtime.True, comparison: runtime.Equal},
		{name: "strict datetime domain", left: runtime.NewDateTime(instant), right: runtime.NewString(instant.Format(time.RFC3339)), expected: runtime.False, comparison: runtime.Greater},
	}

	for _, test := range equalityTests {
		t.Run(test.name, func(t *testing.T) {
			equal, err := runtime.EqualValues(t.Context(), test.left, test.right)
			if err != nil {
				t.Fatalf("EqualValues() error = %v", err)
			}
			if equal != test.expected {
				t.Fatalf("EqualValues() = %v, want %v", equal, test.expected)
			}

			comparison, err := runtime.CompareValues(t.Context(), test.left, test.right)
			if err != nil {
				t.Fatalf("CompareValues() error = %v", err)
			}
			if comparison != test.comparison {
				t.Fatalf("CompareValues() = %v, want %v", comparison, test.comparison)
			}
		})
	}
}

func TestDurationComparisonIsStrict(t *testing.T) {
	duration := runtime.NewDuration(time.Second)
	instant := runtime.NewDateTime(time.Date(2026, time.August, 3, 12, 0, 0, 0, time.UTC))

	for _, other := range []runtime.Value{
		runtime.NewString("1s"),
		runtime.NewInt(1000),
		runtime.NewFloat(1000),
		runtime.True,
		runtime.None,
		runtime.NewArrayWith(runtime.NewString("1s")),
		runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewString("1s")}),
		instant,
		runtime.NewBinary([]byte("1s")),
	} {
		for _, operands := range [][2]runtime.Value{{duration, other}, {other, duration}} {
			equal, err := runtime.EqualValues(t.Context(), operands[0], operands[1])
			if err != nil || equal {
				t.Fatalf("EqualValues(%T, %T) = %v, %v, want false, nil", operands[0], operands[1], equal, err)
			}

			if _, err := runtime.CompareValues(t.Context(), operands[0], operands[1]); !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("CompareValues(%T, %T) error = %v, want ErrInvalidOperation", operands[0], operands[1], err)
			}
		}
	}

	converted, err := runtime.ToDuration(t.Context(), runtime.NewString("1s"))
	if err != nil {
		t.Fatal(err)
	}
	equal, err := runtime.EqualValues(t.Context(), duration, converted)
	if err != nil || !equal {
		t.Fatalf("explicitly converted equality = %v, %v, want true, nil", equal, err)
	}
	comparison, err := runtime.CompareValues(t.Context(), duration, converted)
	if err != nil || comparison != runtime.Equal {
		t.Fatalf("explicitly converted ordering = %v, %v, want Equal, nil", comparison, err)
	}
}

func TestNestedCollectionsUseStrictDurationComparison(t *testing.T) {
	tests := []struct {
		left  runtime.Value
		right runtime.Value
		name  string
	}{
		{
			name:  "array",
			left:  runtime.NewArrayWith(runtime.NewDuration(time.Second)),
			right: runtime.NewArrayWith(runtime.NewString("1s")),
		},
		{
			name:  "object",
			left:  runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewDuration(time.Second)}),
			right: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewString("1s")}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			equal, err := runtime.EqualValues(t.Context(), test.left, test.right)
			if err != nil || equal {
				t.Fatalf("EqualValues() = %v, %v, want false, nil", equal, err)
			}
			if _, err := runtime.CompareValues(t.Context(), test.left, test.right); !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("CompareValues() error = %v, want ErrInvalidOperation", err)
			}
		})
	}
}

func TestUnsupportedHostValuesAreNotEqualOrOrderableWithoutInspection(t *testing.T) {
	left := &opaqueHostValue{hash: 1}
	right := &opaqueHostValue{hash: 1}

	equal, err := runtime.EqualValues(t.Context(), left, right)
	if err != nil {
		t.Fatalf("EqualValues() error = %v", err)
	}
	if equal {
		t.Fatal("EqualValues() = true for unsupported host values")
	}

	if _, err := runtime.CompareValues(t.Context(), left, right); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("CompareValues() error = %v, want ErrInvalidOperation", err)
	}
}

func TestEqualityAndOrderingCapabilitiesAreIndependent(t *testing.T) {
	equatable := equatableOnlyValue(7)
	equal, err := runtime.EqualValues(t.Context(), equatable, equatableOnlyValue(7))
	if err != nil || !equal {
		t.Fatalf("equality-only EqualValues() = %v, %v, want true", equal, err)
	}
	if _, err := runtime.CompareValues(t.Context(), equatable, equatableOnlyValue(7)); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("equality-only CompareValues() error = %v, want ErrInvalidOperation", err)
	}

	comparable := comparableOnlyValue(7)
	equal, err = runtime.EqualValues(t.Context(), comparable, comparableOnlyValue(7))
	if err != nil || equal {
		t.Fatalf("ordering-only EqualValues() = %v, %v, want false", equal, err)
	}
	ordering, err := runtime.CompareValues(t.Context(), comparable, comparableOnlyValue(8))
	if err != nil || ordering != runtime.Less {
		t.Fatalf("ordering-only CompareValues() = %v, %v, want Less", ordering, err)
	}
}

func TestHostComparisonUsesOneSidedCapabilitySymmetrically(t *testing.T) {
	equalityCalls := 0
	compareCalls := 0
	capable := &contractHostValue{
		equal:         true,
		ordering:      runtime.Ordering(-7),
		equalityCalls: &equalityCalls,
		compareCalls:  &compareCalls,
	}
	opaque := &opaqueHostValue{}

	for _, operands := range [][2]runtime.Value{{capable, opaque}, {opaque, capable}} {
		equal, err := runtime.EqualValues(t.Context(), operands[0], operands[1])
		if err != nil {
			t.Fatalf("EqualValues() error = %v", err)
		}
		if !equal {
			t.Fatal("EqualValues() = false, want one-sided result")
		}
	}
	if equalityCalls != 2 {
		t.Fatalf("EqualTo calls = %d, want exactly one per dispatch", equalityCalls)
	}

	comparison, err := runtime.CompareValues(t.Context(), capable, opaque)
	if err != nil {
		t.Fatalf("CompareValues(capable, opaque) error = %v", err)
	}
	if comparison != runtime.Less {
		t.Fatalf("CompareValues(capable, opaque) = %v, want Less", comparison)
	}

	comparison, err = runtime.CompareValues(t.Context(), opaque, capable)
	if err != nil {
		t.Fatalf("CompareValues(opaque, capable) error = %v", err)
	}
	if comparison != runtime.Greater {
		t.Fatalf("CompareValues(opaque, capable) = %v, want Greater", comparison)
	}
	if compareCalls != 2 {
		t.Fatalf("Compare calls = %d, want exactly one per dispatch", compareCalls)
	}
}

func TestHostComparisonNormalizesPositiveAndNegativeResults(t *testing.T) {
	capable := &contractHostValue{ordering: runtime.Ordering(-42)}
	opaque := &opaqueHostValue{}

	comparison, err := runtime.CompareValues(t.Context(), capable, opaque)
	if err != nil || comparison != runtime.Less {
		t.Fatalf("negative Compare result = %v, %v, want Less", comparison, err)
	}

	capable.ordering = runtime.Ordering(42)
	comparison, err = runtime.CompareValues(t.Context(), capable, opaque)
	if err != nil || comparison != runtime.Greater {
		t.Fatalf("positive Compare result = %v, %v, want Greater", comparison, err)
	}
}

func TestHostComparisonChoosesCanonicalReceiverOnce(t *testing.T) {
	aEqualityCalls := 0
	aCompareCalls := 0
	zEqualityCalls := 0
	zCompareCalls := 0
	a := &canonicalAValue{&contractHostValue{
		equal:         true,
		ordering:      runtime.Ordering(-9),
		equalityCalls: &aEqualityCalls,
		compareCalls:  &aCompareCalls,
	}}
	z := &canonicalZValue{&contractHostValue{
		equal:         false,
		ordering:      runtime.Ordering(11),
		equalityCalls: &zEqualityCalls,
		compareCalls:  &zCompareCalls,
	}}

	for _, operands := range [][2]runtime.Value{{a, z}, {z, a}} {
		equal, err := runtime.EqualValues(t.Context(), operands[0], operands[1])
		if err != nil {
			t.Fatalf("EqualValues() error = %v", err)
		}
		if !equal {
			t.Fatal("EqualValues() did not use canonical A receiver")
		}
	}
	if aEqualityCalls != 2 || zEqualityCalls != 0 {
		t.Fatalf("EqualTo calls = A:%d Z:%d, want A:2 Z:0", aEqualityCalls, zEqualityCalls)
	}

	comparison, err := runtime.CompareValues(t.Context(), a, z)
	if err != nil {
		t.Fatalf("CompareValues(A, Z) error = %v", err)
	}
	if comparison != runtime.Less {
		t.Fatalf("CompareValues(A, Z) = %v, want Less", comparison)
	}
	comparison, err = runtime.CompareValues(t.Context(), z, a)
	if err != nil {
		t.Fatalf("CompareValues(Z, A) error = %v", err)
	}
	if comparison != runtime.Greater {
		t.Fatalf("CompareValues(Z, A) = %v, want Greater", comparison)
	}
	if aCompareCalls != 2 || zCompareCalls != 0 {
		t.Fatalf("Compare calls = A:%d Z:%d, want A:2 Z:0", aCompareCalls, zCompareCalls)
	}
}

func TestHostComparisonRejectsIncompatibleTypedDomainsWithoutDispatch(t *testing.T) {
	equalityCalls := 0
	compareCalls := 0
	left := &contractHostValue{equalityCalls: &equalityCalls, compareCalls: &compareCalls}
	right := &contractHostValue{
		typ:           comparisonOtherType,
		equalityCalls: &equalityCalls,
		compareCalls:  &compareCalls,
	}

	equal, err := runtime.EqualValues(t.Context(), left, right)
	if err != nil {
		t.Fatalf("EqualValues() error = %v", err)
	}
	if equal {
		t.Fatal("EqualValues() = true for incompatible domains")
	}
	if _, err := runtime.CompareValues(t.Context(), left, right); !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("CompareValues() error = %v, want ErrInvalidOperation", err)
	}
	if equalityCalls != 0 || compareCalls != 0 {
		t.Fatalf("capability calls = equality:%d compare:%d, want zero", equalityCalls, compareCalls)
	}
}

func TestHostComparisonDomainSelectionIsSymmetricWithOneUnstableType(t *testing.T) {
	equalityCalls := 0
	compareCalls := 0
	stable := &contractHostValue{
		equal:         true,
		equalityCalls: &equalityCalls,
		compareCalls:  &compareCalls,
	}
	unstable := &contractHostValue{
		unknownType:   true,
		equal:         true,
		equalityCalls: &equalityCalls,
		compareCalls:  &compareCalls,
	}

	for _, operands := range [][2]runtime.Value{{stable, unstable}, {unstable, stable}} {
		equal, err := runtime.EqualValues(t.Context(), operands[0], operands[1])
		if err != nil || equal {
			t.Fatalf("EqualValues() = %v, %v, want incompatible", equal, err)
		}
		if _, err := runtime.CompareValues(t.Context(), operands[0], operands[1]); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("CompareValues() error = %v, want ErrInvalidOperation", err)
		}
	}
	if equalityCalls != 0 || compareCalls != 0 {
		t.Fatalf("capability calls = equality:%d compare:%d, want zero", equalityCalls, compareCalls)
	}
}

func TestHostComparisonPropagatesErrorsAndCapabilityCancellation(t *testing.T) {
	operationalErr := errors.New("remote comparison failed")
	left := &contractHostValue{equalityErr: operationalErr, comparisonErr: operationalErr}
	right := &contractHostValue{}

	if _, err := runtime.EqualValues(t.Context(), left, right); !errors.Is(err, operationalErr) {
		t.Fatalf("EqualValues() error = %v, want operational error", err)
	}
	if _, err := runtime.CompareValues(t.Context(), left, right); !errors.Is(err, operationalErr) {
		t.Fatalf("CompareValues() error = %v, want operational error", err)
	}

	equalityCalls := 0
	compareCalls := 0
	cancelled := &contractHostValue{equalityCalls: &equalityCalls, compareCalls: &compareCalls}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := runtime.EqualValues(ctx, cancelled, right); !errors.Is(err, context.Canceled) {
		t.Fatalf("EqualValues() error = %v, want context.Canceled", err)
	}
	if _, err := runtime.CompareValues(ctx, cancelled, right); !errors.Is(err, context.Canceled) {
		t.Fatalf("CompareValues() error = %v, want context.Canceled", err)
	}
	if equalityCalls != 1 || compareCalls != 1 {
		t.Fatalf("cancelled capability calls = equality:%d compare:%d, want one each", equalityCalls, compareCalls)
	}
}

func TestRangeEqualityOrderingAndLaws(t *testing.T) {
	values := []*runtime.Range{
		runtime.NewRange(1, 2),
		runtime.NewRange(1, 3),
		runtime.NewRange(2, 0),
	}

	for _, value := range values {
		equal, err := runtime.EqualValues(t.Context(), value, value)
		if err != nil || !equal {
			t.Fatalf("reflexive EqualValues(%v) = %v, %v", value, equal, err)
		}
		comparison, err := runtime.CompareValues(t.Context(), value, value)
		if err != nil || comparison != runtime.Equal {
			t.Fatalf("reflexive CompareValues(%v) = %v, %v", value, comparison, err)
		}
	}

	for idx := 0; idx < len(values)-1; idx++ {
		left := values[idx]
		right := values[idx+1]
		forward, err := runtime.CompareValues(t.Context(), left, right)
		if err != nil || forward != runtime.Less {
			t.Fatalf("CompareValues(%v, %v) = %v, %v, want Less", left, right, forward, err)
		}
		reverse, err := runtime.CompareValues(t.Context(), right, left)
		if err != nil || reverse != runtime.Greater {
			t.Fatalf("CompareValues(%v, %v) = %v, %v, want Greater", right, left, reverse, err)
		}
		equal, err := runtime.EqualValues(t.Context(), left, right)
		if err != nil || equal || forward == runtime.Equal {
			t.Fatalf("equality/order agreement for %v and %v = %v, %v, %v", left, right, equal, forward, err)
		}
	}

	comparison, err := runtime.CompareValues(t.Context(), values[0], values[2])
	if err != nil || comparison != runtime.Less {
		t.Fatalf("transitive lexicographic comparison = %v, %v, want Less", comparison, err)
	}
	equalCopy, err := runtime.EqualValues(t.Context(), values[0], runtime.NewRange(1, 2))
	if err != nil || !equalCopy {
		t.Fatalf("equal endpoints = %v, %v, want true", equalCopy, err)
	}
	equalComparison, err := runtime.CompareValues(t.Context(), values[0], runtime.NewRange(1, 2))
	if err != nil || equalComparison != runtime.Equal {
		t.Fatalf("equal endpoint ordering = %v, %v, want Equal", equalComparison, err)
	}
	if values[0].Hash() != runtime.NewRange(1, 2).Hash() {
		t.Fatal("equal ranges have different hashes")
	}
}

func TestNestedCollectionsPropagateComparisonFailures(t *testing.T) {
	operationalErr := errors.New("nested comparison failed")

	tests := []struct {
		left  runtime.Value
		right runtime.Value
		name  string
	}{
		{
			name:  "array",
			left:  runtime.NewArrayWith(&contractHostValue{equalityErr: operationalErr, comparisonErr: operationalErr}),
			right: runtime.NewArrayWith(&contractHostValue{}),
		},
		{
			name:  "object",
			left:  runtime.NewObjectWith(map[string]runtime.Value{"key": &contractHostValue{equalityErr: operationalErr, comparisonErr: operationalErr}}),
			right: runtime.NewObjectWith(map[string]runtime.Value{"key": &contractHostValue{}}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := runtime.EqualValues(t.Context(), test.left, test.right); !errors.Is(err, operationalErr) {
				t.Fatalf("EqualValues() error = %v, want nested error", err)
			}
			if _, err := runtime.CompareValues(t.Context(), test.left, test.right); !errors.Is(err, operationalErr) {
				t.Fatalf("CompareValues() error = %v, want nested error", err)
			}
		})
	}
}

func TestDurationComparisonDoesNotInspectOpaqueHosts(t *testing.T) {
	duration := runtime.NewDuration(time.Second)
	opaque := &opaqueHostValue{}

	for _, operands := range [][2]runtime.Value{{duration, opaque}, {opaque, duration}} {
		equal, err := runtime.EqualValues(t.Context(), operands[0], operands[1])
		if err != nil {
			t.Fatalf("EqualValues() error = %v", err)
		}
		if equal {
			t.Fatal("EqualValues() = true for opaque host and Duration")
		}
		if _, err := runtime.CompareValues(t.Context(), operands[0], operands[1]); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("CompareValues() error = %v, want ErrInvalidOperation", err)
		}
	}
}

func TestCompatibleHostDispatchPrecedesDurationIncompatibility(t *testing.T) {
	equalityCalls := 0
	compareCalls := 0
	host := &contractHostValue{
		typ:           runtime.TypeDuration,
		equal:         true,
		ordering:      runtime.Ordering(-5),
		equalityCalls: &equalityCalls,
		compareCalls:  &compareCalls,
	}
	duration := runtime.NewDuration(time.Second)

	equal, err := runtime.EqualValues(t.Context(), duration, host)
	if err != nil || !equal {
		t.Fatalf("EqualValues(Duration, host) = %v, %v, want host true", equal, err)
	}
	comparison, err := runtime.CompareValues(t.Context(), duration, host)
	if err != nil || comparison != runtime.Greater {
		t.Fatalf("CompareValues(Duration, host) = %v, %v, want reversed Greater", comparison, err)
	}
	if equalityCalls != 1 || compareCalls != 1 {
		t.Fatalf("host calls = equality:%d compare:%d, want one each", equalityCalls, compareCalls)
	}
}
