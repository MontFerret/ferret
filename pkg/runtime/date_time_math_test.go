package runtime_test

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToDateTime(t *testing.T) {
	t.Parallel()

	expected := runtime.NewDateTime(time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC))

	actual, err := runtime.ToDateTime(t.Context(), expected)
	if err != nil || !actual.Equal(expected.Time) {
		t.Fatalf("ToDateTime(DateTime) = %v, %v", actual, err)
	}

	actual, err = runtime.ToDateTime(t.Context(), runtime.NewString("2026-08-01T12:00:00Z"))
	if err != nil || !actual.Equal(expected.Time) {
		t.Fatalf("ToDateTime(String) = %v, %v", actual, err)
	}

	for _, input := range []runtime.Value{runtime.NewInt(0), runtime.NewFloat(0)} {
		if _, err := runtime.ToDateTime(t.Context(), input); !errors.Is(err, runtime.ErrInvalidArgument) ||
			!strings.Contains(err.Error(), "numeric DateTime conversion requires an explicit epoch unit") {
			t.Fatalf("ToDateTime(%s) error = %v", runtime.TypeOf(input), err)
		}
	}
}

func TestDateTimeComparisonRemainsStrict(t *testing.T) {
	t.Parallel()

	dateTime := runtime.NewDateTime(time.Date(2026, time.August, 2, 12, 0, 0, 0, time.UTC))
	values := []runtime.Value{
		runtime.NewString("2026-08-02T12:00:00Z"),
		runtime.NewString("2026-08-02T13:00:00Z"),
		runtime.NewString("not-a-date"),
	}

	for _, value := range values {
		value := value
		t.Run(value.String(), func(t *testing.T) {
			equal, err := runtime.EqualValues(t.Context(), dateTime, value)
			if err != nil || equal {
				t.Fatalf("EqualValues(DateTime, %q) = %v, %v; want false, nil", value, equal, err)
			}

			equal, err = runtime.EqualValues(t.Context(), value, dateTime)
			if err != nil || equal {
				t.Fatalf("EqualValues(%q, DateTime) = %v, %v; want false, nil", value, equal, err)
			}

			actual, err := runtime.CompareValues(t.Context(), dateTime, value)
			expected := compareValues(dateTime, value)
			if err != nil || actual != expected {
				t.Fatalf("CompareValues(DateTime, %q) = %d, %v; want %d, nil", value, actual, err, expected)
			}

			actual, err = runtime.CompareValues(t.Context(), value, dateTime)
			expected = compareValues(value, dateTime)
			if err != nil || actual != expected {
				t.Fatalf("CompareValues(%q, DateTime) = %d, %v; want %d, nil", value, actual, err, expected)
			}
		})
	}
}

func TestDateTimeArithmetic(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	baseTime := time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)
	base := runtime.NewDateTime(baseTime)

	assertDateTime := func(name string, value runtime.Value, err error, expected time.Time) {
		t.Helper()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		actual, ok := value.(runtime.DateTime)
		if !ok || !actual.Equal(expected) {
			t.Fatalf("%s = %v, want %v", name, value, expected)
		}
	}

	value, err := runtime.Add(ctx, base, runtime.NewDuration(30*time.Second))
	assertDateTime("DateTime + Duration", value, err, baseTime.Add(30*time.Second))
	value, err = runtime.Add(ctx, runtime.NewDuration(30*time.Second), base)
	assertDateTime("Duration + DateTime", value, err, baseTime.Add(30*time.Second))
	value, err = runtime.Subtract(ctx, base, runtime.NewDuration(5*time.Minute))
	assertDateTime("DateTime - Duration", value, err, baseTime.Add(-5*time.Minute))

	right := runtime.NewDateTime(baseTime.Add(-30 * time.Second))
	difference, err := runtime.Subtract(ctx, base, right)
	if err != nil || difference != runtime.NewDuration(30*time.Second) {
		t.Fatalf("DateTime - DateTime = %v, %v", difference, err)
	}
	concatenated, err := runtime.Add(ctx, runtime.NewString("at "), base)
	if err != nil || concatenated != runtime.NewString("at "+base.String()) {
		t.Fatalf("String + DateTime = %v, %v", concatenated, err)
	}

	for name, operation := range map[string]func() (runtime.Value, error){
		"DateTime + DateTime": func() (runtime.Value, error) { return runtime.Add(ctx, base, base) },
		"DateTime + Number": func() (runtime.Value, error) {
			return runtime.Add(ctx, base, runtime.NewInt(1500))
		},
		"DateTime - String": func() (runtime.Value, error) {
			return runtime.Subtract(ctx, base, runtime.NewString("5m"))
		},
		"Duration - DateTime": func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.NewDuration(time.Second), base)
		},
		"DateTime * Number": func() (runtime.Value, error) { return runtime.Multiply(ctx, base, runtime.NewInt(2)) },
		"DateTime / Number": func() (runtime.Value, error) { return runtime.Divide(ctx, base, runtime.NewInt(2)) },
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := operation(); !errors.Is(err, runtime.ErrInvalidOperation) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestDateTimeArithmeticUsesCanonicalInstants(t *testing.T) {
	t.Parallel()

	utc := runtime.NewDateTime(time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC))
	offset := runtime.NewDateTime(time.Date(2026, time.August, 1, 7, 0, 0, 0, time.FixedZone("EST", -5*60*60)))

	difference, err := runtime.Subtract(t.Context(), utc, offset)
	if err != nil || difference != runtime.ZeroDuration {
		t.Fatalf("timezone-equivalent difference = %v, %v", difference, err)
	}

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	before := runtime.NewDateTime(time.Date(2026, time.March, 8, 1, 30, 0, 0, location))
	value, err := runtime.Add(t.Context(), before, runtime.NewDuration(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	after := value.(runtime.DateTime)
	if after.Hour() != 3 || after.Minute() != 30 {
		t.Fatalf("DST addition = %v, want 03:30", after)
	}
}

func TestDateTimeArithmeticRangeErrors(t *testing.T) {
	t.Parallel()

	const unixToInternal = int64(62_135_596_800)
	farFuture := runtime.NewDateTime(time.Unix(math.MaxInt64-unixToInternal, 0))
	if _, err := runtime.Add(t.Context(), farFuture, runtime.NewDuration(time.Second)); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("DateTime addition error = %v", err)
	}

	left := runtime.NewDateTime(time.Date(10_000, time.January, 1, 0, 0, 0, 0, time.UTC))
	right := runtime.NewDateTime(time.Date(-10_000, time.January, 1, 0, 0, 0, 0, time.UTC))
	if _, err := runtime.Subtract(t.Context(), left, right); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("DateTime difference error = %v", err)
	}
}

func TestDateTimeArithmeticHandlesMinimumDuration(t *testing.T) {
	t.Parallel()

	baseTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	base := runtime.NewDateTime(baseTime)
	minimum := runtime.Duration(math.MinInt64)

	value, err := runtime.Add(t.Context(), base, minimum)
	if err != nil {
		t.Fatal(err)
	}
	before := value.(runtime.DateTime)
	if !before.Equal(baseTime.Add(time.Duration(math.MinInt64))) {
		t.Fatalf("DateTime + MinInt64 = %v", before)
	}

	value, err = runtime.Subtract(t.Context(), base, minimum)
	if err != nil {
		t.Fatal(err)
	}
	after := value.(runtime.DateTime)
	expectedAfter := baseTime.Add(time.Duration(math.MaxInt64)).Add(time.Nanosecond)
	if !after.Equal(expectedAfter) {
		t.Fatalf("DateTime - MinInt64 = %v, want %v", after, expectedAfter)
	}

	difference, err := runtime.Subtract(t.Context(), before, base)
	if err != nil || difference != minimum {
		t.Fatalf("minimum DateTime difference = %v, %v", difference, err)
	}
}
