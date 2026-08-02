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

func TestDateTimeCheckedComparisonRemainsStrict(t *testing.T) {
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
			equal, err := runtime.EqualChecked(t.Context(), dateTime, value)
			if err != nil || equal {
				t.Fatalf("EqualChecked(DateTime, %q) = %v, %v; want false, nil", value, equal, err)
			}

			equal, err = runtime.EqualChecked(t.Context(), value, dateTime)
			if err != nil || equal {
				t.Fatalf("EqualChecked(%q, DateTime) = %v, %v; want false, nil", value, equal, err)
			}

			actual, err := runtime.CompareChecked(t.Context(), dateTime, value)
			expected := runtime.CompareValues(dateTime, value)
			if err != nil || actual != expected {
				t.Fatalf("CompareChecked(DateTime, %q) = %d, %v; want %d, nil", value, actual, err, expected)
			}

			actual, err = runtime.CompareChecked(t.Context(), value, dateTime)
			expected = runtime.CompareValues(value, dateTime)
			if err != nil || actual != expected {
				t.Fatalf("CompareChecked(%q, DateTime) = %d, %v; want %d, nil", value, actual, err, expected)
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

	value, err := runtime.AddChecked(ctx, base, runtime.NewDuration(30*time.Second))
	assertDateTime("DateTime + Duration", value, err, baseTime.Add(30*time.Second))
	value, err = runtime.AddChecked(ctx, runtime.NewString("30s"), base)
	assertDateTime("coercible Duration + DateTime", value, err, baseTime.Add(30*time.Second))
	value, err = runtime.AddChecked(ctx, base, runtime.NewInt(1500))
	assertDateTime("DateTime + milliseconds", value, err, baseTime.Add(1500*time.Millisecond))
	value, err = runtime.SubtractChecked(ctx, base, runtime.NewString("5m"))
	assertDateTime("DateTime - duration string", value, err, baseTime.Add(-5*time.Minute))

	difference, err := runtime.SubtractChecked(ctx, base, runtime.NewString("2026-08-01T11:59:30Z"))
	if err != nil || difference != runtime.NewDuration(30*time.Second) {
		t.Fatalf("DateTime - DateTime string = %v, %v", difference, err)
	}

	for name, operation := range map[string]func() (runtime.Value, error){
		"DateTime + DateTime": func() (runtime.Value, error) { return runtime.AddChecked(ctx, base, base) },
		"Duration - DateTime": func() (runtime.Value, error) {
			return runtime.SubtractChecked(ctx, runtime.NewDuration(time.Second), base)
		},
		"DateTime * Number": func() (runtime.Value, error) { return runtime.MultiplyChecked(ctx, base, runtime.NewInt(2)) },
		"DateTime / Number": func() (runtime.Value, error) { return runtime.DivideChecked(ctx, base, runtime.NewInt(2)) },
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

	difference, err := runtime.SubtractChecked(t.Context(), utc, offset)
	if err != nil || difference != runtime.ZeroDuration {
		t.Fatalf("timezone-equivalent difference = %v, %v", difference, err)
	}

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	before := runtime.NewDateTime(time.Date(2026, time.March, 8, 1, 30, 0, 0, location))
	value, err := runtime.AddChecked(t.Context(), before, runtime.NewDuration(time.Hour))
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
	if _, err := runtime.AddChecked(t.Context(), farFuture, runtime.NewDuration(time.Second)); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("DateTime addition error = %v", err)
	}

	left := runtime.NewDateTime(time.Date(10_000, time.January, 1, 0, 0, 0, 0, time.UTC))
	right := runtime.NewDateTime(time.Date(-10_000, time.January, 1, 0, 0, 0, 0, time.UTC))
	if _, err := runtime.SubtractChecked(t.Context(), left, right); !errors.Is(err, runtime.ErrRange) {
		t.Fatalf("DateTime difference error = %v", err)
	}
}

func TestDateTimeArithmeticHandlesMinimumDuration(t *testing.T) {
	t.Parallel()

	baseTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	base := runtime.NewDateTime(baseTime)
	minimum := runtime.Duration(math.MinInt64)

	value, err := runtime.AddChecked(t.Context(), base, minimum)
	if err != nil {
		t.Fatal(err)
	}
	before := value.(runtime.DateTime)
	if !before.Equal(baseTime.Add(time.Duration(math.MinInt64))) {
		t.Fatalf("DateTime + MinInt64 = %v", before)
	}

	value, err = runtime.SubtractChecked(t.Context(), base, minimum)
	if err != nil {
		t.Fatal(err)
	}
	after := value.(runtime.DateTime)
	expectedAfter := baseTime.Add(time.Duration(math.MaxInt64)).Add(time.Nanosecond)
	if !after.Equal(expectedAfter) {
		t.Fatalf("DateTime - MinInt64 = %v, want %v", after, expectedAfter)
	}

	difference, err := runtime.SubtractChecked(t.Context(), before, base)
	if err != nil || difference != minimum {
		t.Fatalf("minimum DateTime difference = %v, %v", difference, err)
	}
}
