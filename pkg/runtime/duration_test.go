package runtime_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestParseDuration(t *testing.T) {
	t.Parallel()

	tests := map[string]runtime.Duration{
		"250ms":                    runtime.NewDuration(250 * time.Millisecond),
		"1.5s":                     runtime.NewDuration(1500 * time.Millisecond),
		"1e3ms":                    runtime.NewDuration(time.Second),
		"2M":                       runtime.NewDuration(2 * time.Minute),
		"1d":                       runtime.NewDuration(24 * time.Hour),
		"0.0000005ms":              runtime.ZeroDuration,
		"-0.0000005ms":             runtime.ZeroDuration,
		"1h30m":                    runtime.NewDuration(90 * time.Minute),
		"1D2H30M":                  runtime.NewDuration(26*time.Hour + 30*time.Minute),
		"0.5ns0.5ns":               runtime.NewDuration(time.Nanosecond),
		"2562047h47m16.854775807s": runtime.Duration(math.MaxInt64),
		"9223372036.854775807s":    runtime.Duration(math.MaxInt64),
		"-9223372036.854775808s":   runtime.Duration(math.MinInt64),
	}

	for input, expected := range tests {
		input, expected := input, expected
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			actual, err := runtime.ParseDuration(input)
			if err != nil {
				t.Fatalf("ParseDuration(%q): %v", input, err)
			}
			if actual != expected {
				t.Fatalf("ParseDuration(%q) = %s, want %s", input, actual, expected)
			}
		})
	}

	for _, input := range []string{
		"9223372036.854775808s",
		"-9223372036.854775809s",
		"2562047h47m16.854775808s",
	} {
		if _, err := runtime.ParseDuration(input); !errors.Is(err, runtime.ErrRange) {
			t.Fatalf("ParseDuration(%q) error = %v, want range error", input, err)
		}
	}
}

func TestDurationRuntimeContract(t *testing.T) {
	t.Parallel()

	duration := runtime.NewDuration(1500 * time.Millisecond)
	if duration.Type() != runtime.TypeDuration || runtime.TypeOf(duration) != runtime.TypeDuration {
		t.Fatalf("duration type was not preserved")
	}
	if got := duration.String(); got != "1.5s" {
		t.Fatalf("String() = %q, want 1.5s", got)
	}
	if got := duration.Unwrap(); got != 1500*time.Millisecond {
		t.Fatalf("Unwrap() = %v", got)
	}
	if duration.Copy() != duration {
		t.Fatalf("Copy() did not preserve the immutable duration")
	}
	if duration.Hash() != runtime.NewDuration(1500*time.Millisecond).Hash() {
		t.Fatalf("equivalent durations must have equal hashes")
	}
	if !runtime.IsScalar(duration) || !runtime.ToBoolean(duration) || runtime.ToBoolean(runtime.ZeroDuration) {
		t.Fatalf("duration scalar/truthiness contract was not preserved")
	}
	list, err := runtime.ToList(t.Context(), duration)
	if err != nil {
		t.Fatalf("ToList(duration) = %v, %v", list, err)
	}
	length, err := list.Length(t.Context())
	if err != nil || length != 1 {
		t.Fatalf("duration list length = %d, %v", length, err)
	}
	if compareValues(duration, runtime.NewDuration(time.Second)) <= 0 {
		t.Fatalf("duration ordering is incorrect")
	}
	for _, other := range []runtime.Value{runtime.NewInt(1), runtime.ZeroFloat, runtime.EmptyString} {
		if _, err := runtime.CompareValues(t.Context(), duration, other); !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("Duration ordering with %T error = %v, want ErrInvalidOperation", other, err)
		}
	}

	value, err := runtime.ValueOf(1500 * time.Millisecond)
	if err != nil || value != duration {
		t.Fatalf("ValueOf(time.Duration) = %v, %v", value, err)
	}
	if _, err := runtime.CastDuration(runtime.NewInt(1)); err == nil {
		t.Fatal("CastDuration accepted an Int")
	}
	if _, err := runtime.CastDuration(runtime.NewString("1.5s")); err == nil {
		t.Fatal("CastDuration accepted a String")
	}
	if _, err := runtime.CastInt(duration); err == nil {
		t.Fatal("CastInt accepted a Duration")
	}
	if _, err := runtime.CastFloat(duration); err == nil {
		t.Fatal("CastFloat accepted a Duration")
	}
	if _, err := runtime.CastString(duration); err == nil {
		t.Fatal("CastString accepted a Duration")
	}
	if err := runtime.AssertNumber(duration); err == nil {
		t.Fatal("AssertNumber accepted a Duration")
	}
	if err := runtime.AssertDuration(duration); err != nil {
		t.Fatalf("AssertDuration: %v", err)
	}
}

func TestDurationArithmetic(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	second := runtime.NewDuration(time.Second)
	halfSecond := runtime.NewDuration(500 * time.Millisecond)

	assertValue := func(name string, actual runtime.Value, err error, expected runtime.Value) {
		t.Helper()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if compareValues(actual, expected) != 0 || runtime.TypeOf(actual) != runtime.TypeOf(expected) {
			t.Fatalf("%s = %v (%s), want %v (%s)", name, actual, runtime.TypeOf(actual), expected, runtime.TypeOf(expected))
		}
	}

	actual, err := runtime.Add(ctx, second, halfSecond)
	assertValue("add", actual, err, runtime.NewDuration(1500*time.Millisecond))
	actual, err = runtime.Subtract(ctx, second, halfSecond)
	assertValue("subtract", actual, err, halfSecond)
	actual, err = runtime.Multiply(ctx, second, runtime.NewFloat(1.5))
	assertValue("multiply", actual, err, runtime.NewDuration(1500*time.Millisecond))
	actual, err = runtime.Multiply(ctx, runtime.NewInt(2), halfSecond)
	assertValue("reverse multiply", actual, err, second)
	actual, err = runtime.Divide(ctx, second, runtime.NewInt(2))
	assertValue("divide number", actual, err, halfSecond)
	actual, err = runtime.Divide(ctx, second, halfSecond)
	assertValue("exact ratio", actual, err, runtime.NewInt(2))
	actual, err = runtime.Divide(ctx, second, runtime.NewDuration(3*time.Second))
	assertValue("fractional ratio", actual, err, runtime.NewFloat(1.0/3.0))

	actual, err = runtime.Multiply(ctx, runtime.NewDuration(time.Nanosecond), runtime.NewFloat(0.5))
	assertValue("positive tie", actual, err, runtime.ZeroDuration)
	actual, err = runtime.Multiply(ctx, runtime.NewDuration(-time.Nanosecond), runtime.NewFloat(0.5))
	assertValue("negative tie", actual, err, runtime.ZeroDuration)
	actual, err = runtime.Divide(ctx, runtime.NewDuration(time.Nanosecond), runtime.NewInt(2))
	assertValue("positive integer division tie", actual, err, runtime.ZeroDuration)
	actual, err = runtime.Divide(ctx, runtime.NewDuration(-time.Nanosecond), runtime.NewInt(2))
	assertValue("negative integer division tie", actual, err, runtime.ZeroDuration)
	actual, err = runtime.Positive(second)
	assertValue("unary positive", actual, err, second)
	actual, err = runtime.Negative(second)
	assertValue("unary negative", actual, err, runtime.NewDuration(-time.Second))
	actual, err = runtime.Multiply(ctx, runtime.Duration(math.MaxInt64), runtime.NewFloat(1))
	assertValue("maximum float identity", actual, err, runtime.Duration(math.MaxInt64))
	actual, err = runtime.Divide(ctx, runtime.Duration(math.MinInt64), runtime.NewFloat(1))
	assertValue("minimum float identity", actual, err, runtime.Duration(math.MinInt64))

	invalid := []struct {
		fn   func() (runtime.Value, error)
		name string
	}{
		{name: "duration multiplication", fn: func() (runtime.Value, error) { return runtime.Multiply(ctx, second, second) }},
		{name: "reverse division", fn: func() (runtime.Value, error) { return runtime.Divide(ctx, runtime.NewInt(1), second) }},
		{name: "zero division", fn: func() (runtime.Value, error) { return runtime.Divide(ctx, second, runtime.ZeroInt) }},
		{name: "zero duration division", fn: func() (runtime.Value, error) {
			return runtime.Divide(ctx, second, runtime.ZeroDuration)
		}},
		{name: "nan scaling", fn: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, second, runtime.NewFloat(math.NaN()))
		}},
		{name: "infinite division", fn: func() (runtime.Value, error) {
			return runtime.Divide(ctx, second, runtime.NewFloat(math.Inf(1)))
		}},
		{name: "overflow", fn: func() (runtime.Value, error) {
			return runtime.Add(ctx, runtime.Duration(math.MaxInt64), runtime.Duration(1))
		}},
		{name: "subtraction overflow", fn: func() (runtime.Value, error) {
			return runtime.Subtract(ctx, runtime.Duration(math.MinInt64), runtime.Duration(1))
		}},
		{name: "multiplication overflow", fn: func() (runtime.Value, error) {
			return runtime.Multiply(ctx, runtime.Duration(math.MaxInt64), runtime.NewInt(2))
		}},
		{name: "negation overflow", fn: func() (runtime.Value, error) {
			return runtime.Negative(runtime.Duration(math.MinInt64))
		}},
		{name: "increment", fn: func() (runtime.Value, error) { return runtime.Increment(ctx, second) }},
		{name: "decrement", fn: func() (runtime.Value, error) { return runtime.Decrement(ctx, second) }},
		{name: "modulo", fn: func() (runtime.Value, error) { return runtime.Modulo(ctx, second, runtime.NewInt(2)) }},
	}

	for _, test := range invalid {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.fn()
			if err == nil || (!errors.Is(err, runtime.ErrInvalidOperation) && !errors.Is(err, runtime.ErrRange)) {
				t.Fatalf("got %v, want duration operation error", err)
			}
		})
	}
}
