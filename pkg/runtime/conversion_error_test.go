package runtime

import (
	"context"
	"errors"
	"math"
	"strconv"
	"testing"
	"time"
)

func TestToAPIsMarkConversionOwnedErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		operation func() error
		target    Type
		name      string
	}{
		{
			name:   "duration",
			target: TypeDuration,
			operation: func() error {
				_, err := ToDuration(t.Context(), NewString("invalid"))
				return err
			},
		},
		{
			name:   "date time",
			target: TypeDateTime,
			operation: func() error {
				_, err := ToDateTime(t.Context(), NewString("invalid"))
				return err
			},
		},
		{
			name:   "date time epoch",
			target: TypeDateTime,
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewInt(1), NewString("invalid"))
				return err
			},
		},
		{
			name:   "integer",
			target: TypeInt,
			operation: func() error {
				_, err := ToInt(t.Context(), NewString("invalid"))
				return err
			},
		},
		{
			name:   "float",
			target: TypeFloat,
			operation: func() error {
				_, err := ToFloat(t.Context(), NewString("invalid"))
				return err
			},
		},
		{
			name:   "number delegates to float",
			target: TypeFloat,
			operation: func() error {
				_, err := ToNumber(t.Context(), NewString("invalid"))
				return err
			},
		},
		{
			name:   "integer default delegates to integer",
			target: TypeInt,
			operation: func() error {
				_, err := ToIntDefault(t.Context(), NewString("invalid"), NewInt(1))
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.operation()
			if err == nil {
				t.Fatal("conversion unexpectedly succeeded")
			}
			if !isConversionErrorTo(err, test.target) {
				t.Fatalf("error %v is not marked for %s", err, test.target)
			}
			if test.target != TypeDuration && isConversionErrorTo(err, TypeDuration) {
				t.Fatalf("error %v was also marked for Duration", err)
			}
		})
	}
}

func TestConversionErrorsPreservePublicCauses(t *testing.T) {
	t.Parallel()

	_, expectedIntErr := strconv.ParseInt("invalid", 10, 64)
	_, intErr := ToInt(t.Context(), NewString("invalid"))
	if intErr.Error() != expectedIntErr.Error() {
		t.Fatalf("ToInt error = %q, want %q", intErr, expectedIntErr)
	}
	var intParseErr *strconv.NumError
	if !errors.As(intErr, &intParseErr) {
		t.Fatalf("ToInt error %v does not retain strconv.NumError", intErr)
	}

	_, expectedFloatErr := strconv.ParseFloat("invalid", 64)
	_, floatErr := ToFloat(t.Context(), NewString("invalid"))
	if floatErr.Error() != expectedFloatErr.Error() {
		t.Fatalf("ToFloat error = %q, want %q", floatErr, expectedFloatErr)
	}
	var floatParseErr *strconv.NumError
	if !errors.As(floatErr, &floatParseErr) {
		t.Fatalf("ToFloat error %v does not retain strconv.NumError", floatErr)
	}

	if _, err := ToInt(t.Context(), NewObject()); !errors.Is(err, ErrInvalidType) {
		t.Fatalf("ToInt unsupported type error = %v, want ErrInvalidType", err)
	}
	if _, err := ToFloat(t.Context(), NewObject()); !errors.Is(err, ErrInvalidType) {
		t.Fatalf("ToFloat unsupported type error = %v, want ErrInvalidType", err)
	}
	if _, err := ToDateTime(t.Context(), NewString("invalid")); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("ToDateTime malformed value error = %v, want ErrInvalidArgument", err)
	}
	if _, err := ToDateTime(t.Context(), NewInt(1)); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("ToDateTime numeric value error = %v, want ErrInvalidArgument", err)
	}
	if _, err := ToDateTimeEpoch(t.Context(), NewInt64(math.MaxInt64), NewString("s")); !errors.Is(err, ErrRange) {
		t.Fatalf("ToDateTimeEpoch overflow error = %v, want ErrRange", err)
	}
	if _, err := ToDuration(t.Context(), NewInt64(9_223_372_036_855)); !errors.Is(err, ErrRange) {
		t.Fatalf("ToDuration overflow error = %v, want ErrRange", err)
	}
}

func TestToDateTimeEpochMarksConversionFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		category  error
		operation func() error
		name      string
	}{
		{
			name: "missing unit",
			operation: func() error {
				_, err := ToDateTime(t.Context(), NewInt(1))
				return err
			},
			category: ErrInvalidArgument,
		},
		{
			name: "unit type",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewInt(1), None)
				return err
			},
			category: ErrInvalidType,
		},
		{
			name: "input type",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewObject(), NewString("s"))
				return err
			},
			category: ErrInvalidType,
		},
		{
			name: "unit misuse",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewString("1"), NewString("s"))
				return err
			},
			category: ErrInvalidArgument,
		},
		{
			name: "unknown unit",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewInt(1), NewString("minutes"))
				return err
			},
			category: ErrInvalidArgument,
		},
		{
			name: "non-finite",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewFloat(math.NaN()), NewString("s"))
				return err
			},
			category: ErrInvalidArgument,
		},
		{
			name: "range",
			operation: func() error {
				_, err := ToDateTimeEpoch(t.Context(), NewInt64(math.MaxInt64), NewString("s"))
				return err
			},
			category: ErrRange,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.operation()
			if !errors.Is(err, test.category) {
				t.Fatalf("error = %v, want %v", err, test.category)
			}
			if !isConversionErrorTo(err, TypeDateTime) {
				t.Fatalf("error %v is not marked for DateTime", err)
			}
			if isConversionErrorTo(err, TypeDuration) {
				t.Fatalf("error %v was also marked for Duration", err)
			}
		})
	}
}

func TestToAPIsLeaveOperationalErrorsUnmarked(t *testing.T) {
	t.Parallel()

	operationErr := errors.New("conversion source failed")
	list := newConversionErrorList(operationErr, operationErr, NewInt(1))
	iterable := newConversionErrorIterable(operationErr)
	tests := []struct {
		operation func() error
		name      string
	}{
		{
			name: "duration list length",
			operation: func() error {
				_, err := ToDuration(t.Context(), list)
				return err
			},
		},
		{
			name: "integer list iteration",
			operation: func() error {
				_, err := ToInt(t.Context(), list)
				return err
			},
		},
		{
			name: "float list iteration",
			operation: func() error {
				_, err := ToFloat(t.Context(), list)
				return err
			},
		},
		{
			name: "list iterable setup",
			operation: func() error {
				_, err := ToList(t.Context(), iterable)
				return err
			},
		},
		{
			name: "map iterable setup",
			operation: func() error {
				_, err := ToMap(t.Context(), iterable)
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.operation()
			if !errors.Is(err, operationErr) {
				t.Fatalf("operation error = %v, want %v", err, operationErr)
			}
			for _, target := range []Type{TypeDuration, TypeDateTime, TypeInt, TypeFloat} {
				if isConversionErrorTo(err, target) {
					t.Fatalf("operational error %v was marked for %s", err, target)
				}
			}
		})
	}
}

func TestEqualValuesDoesNotInspectOpaqueDurationLists(t *testing.T) {
	t.Parallel()

	duration := NewDuration(time.Second)
	durationCause := errors.New("duration conversion failed")
	durationErr := newConversionError(TypeDuration, durationCause)
	equal, err := EqualValues(
		t.Context(),
		duration,
		newConversionErrorList(durationErr, nil),
	)
	if err != nil || equal {
		t.Fatalf("Duration conversion equality = %v, %v; want false, nil", equal, err)
	}

	intCause := errors.New("integer conversion failed")
	intErr := newConversionError(TypeInt, intCause)
	equal, err = EqualValues(
		t.Context(),
		duration,
		newConversionErrorList(intErr, nil),
	)
	if err != nil || equal {
		t.Fatalf("opaque non-Duration conversion equality = %v, %v; want false, nil", equal, err)
	}

	equal, err = EqualValues(
		context.Background(),
		duration,
		newConversionErrorList(context.Canceled, nil),
	)
	if err != nil || equal {
		t.Fatalf("opaque stored cancellation equality = %v, %v; want false, nil", equal, err)
	}
}

func isConversionErrorTo(err error, target Type) bool {
	var conversionErr *conversionError
	return errors.As(err, &conversionErr) && conversionErr.targets(target)
}
