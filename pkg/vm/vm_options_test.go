package vm

import (
	"errors"
	"testing"

	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func TestNewWithAppliesOrderedOptionsAndPreservesNoOps(t *testing.T) {
	t.Parallel()

	instance, err := NewWith(
		newPoolTestProgram(),
		nil,
		WithShapeCacheLimit(11),
		WithShapeCacheLimit(17),
		WithShapeCacheLimit(0),
		WithFastObjectDictThreshold(19),
		WithFastObjectDictThreshold(23),
		WithFastObjectDictThreshold(-1),
		WithPanicPolicy(PanicRecover),
		WithPanicPolicy(PanicPropagate),
		WithPanicPolicy(PanicPolicy(255)),
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
