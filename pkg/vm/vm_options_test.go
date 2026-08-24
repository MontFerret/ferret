package vm

import (
	"errors"
	"testing"

	"github.com/ziflex/go-options"

	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func TestNewWithAppliesOptionsInOrder(t *testing.T) {
	t.Parallel()

	instance, err := NewWith(
		newPoolTestProgram(),
		nil,
		WithShapeCacheLimit(11),
		WithShapeCacheLimit(17),
		WithFastObjectDictThreshold(19),
		WithFastObjectDictThreshold(23),
		WithPanicPolicy(PanicRecover),
		WithPanicPolicy(PanicPropagate),
	)
	if err != nil {
		t.Fatalf("expected VM construction to succeed, got: %v", err)
	}
	t.Cleanup(func() {
		_ = instance.Close()
	})

	if got, want := instance.config.shapeCacheLimit, 17; got != want {
		t.Fatalf("unexpected shape cache limit: got %d, want %d", got, want)
	}

	if got, want := instance.config.fastObjectDictThreshold, 23; got != want {
		t.Fatalf("unexpected fast object dict threshold: got %d, want %d", got, want)
	}

	if got, want := instance.config.panicPolicy, PanicPropagate; got != want {
		t.Fatalf("unexpected panic policy: got %d, want %d", got, want)
	}
}

func TestNewWithRejectsInvalidPositiveOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		option Option
		name   string
		field  string
		value  string
	}{
		{
			name:   "zero shape cache limit",
			field:  "shape cache limit",
			value:  "0",
			option: WithShapeCacheLimit(0),
		},
		{
			name:   "negative shape cache limit",
			field:  "shape cache limit",
			value:  "-1",
			option: WithShapeCacheLimit(-1),
		},
		{
			name:   "zero fast object dict threshold",
			field:  "fast object dict threshold",
			value:  "0",
			option: WithFastObjectDictThreshold(0),
		},
		{
			name:   "negative fast object dict threshold",
			field:  "fast object dict threshold",
			value:  "-1",
			option: WithFastObjectDictThreshold(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			instance, err := NewWith(newPoolTestProgram(), tt.option)
			if instance != nil {
				t.Fatal("expected invalid option to return no VM")
			}

			var validationErr options.ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected options.ValidationError, got: %v", err)
			}

			if validationErr.Field != tt.field || validationErr.Value != tt.value {
				t.Fatalf("unexpected validation error: %+v", validationErr)
			}

			var nested options.ValidationError
			if errors.As(validationErr.Reason, &nested) {
				t.Fatalf("expected flat validation error, got nested error: %+v", nested)
			}

			if validationErr.Reason == nil || validationErr.Reason.Error() != "must be positive" {
				t.Fatalf("unexpected validation reason: %v", validationErr.Reason)
			}
		})
	}
}

func TestNewWithRejectsInvalidPanicPolicy(t *testing.T) {
	t.Parallel()

	instance, err := NewWith(newPoolTestProgram(), WithPanicPolicy(PanicPolicy(123)))
	if instance != nil {
		t.Fatal("expected invalid panic policy to return no VM")
	}

	var validationErr options.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected options.ValidationError, got: %v", err)
	}

	if validationErr.Field != "panic policy" || validationErr.Value != "123" {
		t.Fatalf("unexpected validation error: %+v", validationErr)
	}

	var nested options.ValidationError
	if errors.As(validationErr.Reason, &nested) {
		t.Fatalf("expected flat validation error, got nested error: %+v", nested)
	}

	if validationErr.Reason == nil || validationErr.Reason.Error() != "must be one of [0 1]" {
		t.Fatalf("unexpected validation reason: %v", validationErr.Reason)
	}
}

func TestNewWithJoinsValidationErrorsInOptionOrder(t *testing.T) {
	t.Parallel()

	instance, err := NewWith(
		newPoolTestProgram(),
		WithShapeCacheLimit(0),
		WithFastObjectDictThreshold(-1),
		WithPanicPolicy(PanicPolicy(123)),
	)
	if instance != nil {
		t.Fatal("expected invalid options to return no VM")
	}

	const wantMessage = "shape cache limit: must be positive: value=0\n" +
		"fast object dict threshold: must be positive: value=-1\n" +
		"panic policy: must be one of [0 1]: value=123"
	if err == nil || err.Error() != wantMessage {
		t.Fatalf("unexpected joined error: got %q, want %q", err, wantMessage)
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined error, got: %T", err)
	}

	want := []struct {
		field  string
		value  string
		reason string
	}{
		{field: "shape cache limit", value: "0", reason: "must be positive"},
		{field: "fast object dict threshold", value: "-1", reason: "must be positive"},
		{field: "panic policy", value: "123", reason: "must be one of [0 1]"},
	}
	got := joined.Unwrap()
	if len(got) != len(want) {
		t.Fatalf("unexpected joined error count: got %d, want %d", len(got), len(want))
	}

	for i, wantErr := range want {
		var validationErr options.ValidationError
		if !errors.As(got[i], &validationErr) {
			t.Fatalf("joined error %d = %T, want options.ValidationError", i, got[i])
		}

		if validationErr.Field != wantErr.field || validationErr.Value != wantErr.value {
			t.Fatalf("unexpected joined validation error %d: %+v", i, validationErr)
		}

		if validationErr.Reason == nil || validationErr.Reason.Error() != wantErr.reason {
			t.Fatalf("unexpected joined validation reason %d: %v", i, validationErr.Reason)
		}
	}
}

func TestNewWithUsesDefaultConfig(t *testing.T) {
	t.Parallel()

	instance, err := NewWith(newPoolTestProgram())
	if err != nil {
		t.Fatalf("expected VM construction to succeed, got: %v", err)
	}
	t.Cleanup(func() {
		_ = instance.Close()
	})

	if got, want := instance.config.shapeCacheLimit, defaultShapeCacheLimit; got != want {
		t.Fatalf("unexpected default shape cache limit: got %d, want %d", got, want)
	}

	if got, want := instance.config.fastObjectDictThreshold, defaultFastObjectDictThreshold; got != want {
		t.Fatalf("unexpected default fast object dict threshold: got %d, want %d", got, want)
	}

	if got, want := instance.config.panicPolicy, PanicRecover; got != want {
		t.Fatalf("unexpected default panic policy: got %d, want %d", got, want)
	}
}

func TestNewWithJoinsOptionErrorsAndReturnsNoVM(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first option")
	secondErr := errors.New("second option")
	applied := make([]string, 0, 2)

	instance, err := NewWith(
		newPoolTestProgram(),
		func(cfg *config) error {
			applied = append(applied, "first")
			cfg.shapeCacheLimit = 17

			return firstErr
		},
		nil,
		func(cfg *config) error {
			applied = append(applied, "second")
			cfg.shapeCacheLimit = 23

			return secondErr
		},
	)
	if instance != nil {
		t.Fatal("expected option failure to return no VM")
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected both option errors, got: %v", err)
	}

	if got, want := err.Error(), "first option\nsecond option"; got != want {
		t.Fatalf("unexpected joined error order: got %q, want %q", got, want)
	}

	if got, want := len(applied), 2; got != want {
		t.Fatalf("expected later options to execute after an error: got %d applications, want %d", got, want)
	}

	if applied[0] != "first" || applied[1] != "second" {
		t.Fatalf("unexpected option application order: %v", applied)
	}
}

func TestNewWithValidatesProgramBeforeApplyingOptions(t *testing.T) {
	t.Parallel()

	optionErr := errors.New("option failure")
	applied := false
	instance, err := NewWith(nil, func(*config) error {
		applied = true

		return optionErr
	})
	if instance != nil {
		t.Fatal("expected invalid program to return no VM")
	}

	if err == nil {
		t.Fatal("expected invalid program to fail")
	}

	if errors.Is(err, optionErr) {
		t.Fatalf("expected program validation error to take precedence, got: %v", err)
	}

	if applied {
		t.Fatal("expected invalid program not to apply options")
	}
}

func TestNewWithPropagatesTestingOptionErrors(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("testing option")
	instance, err := NewWith(
		newPoolTestProgram(),
		WithTesting(func(*vmtest.Options) error {
			return wantErr
		}),
	)
	if instance != nil {
		t.Fatal("expected testing option failure to return no VM")
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected testing option error, got: %v", err)
	}
}
