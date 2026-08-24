package test

import (
	"errors"
	"testing"
)

func TestNewOptionsUsesDefaultsAndIgnoresNilOptions(t *testing.T) {
	t.Parallel()

	cfg, err := NewOptions([]Option{nil})
	if err != nil {
		t.Fatalf("expected testing options to succeed, got: %v", err)
	}

	if cfg.BenchmarkMode {
		t.Fatal("expected benchmark mode to be disabled by default")
	}
}

func TestNewOptionsAppliesOptionsInOrder(t *testing.T) {
	t.Parallel()

	cfg, err := NewOptions([]Option{
		WithBenchmarkMode(),
		func(opts *Options) error {
			opts.BenchmarkMode = false

			return nil
		},
		WithBenchmarkMode(),
	})
	if err != nil {
		t.Fatalf("expected testing options to succeed, got: %v", err)
	}

	if !cfg.BenchmarkMode {
		t.Fatal("expected the later option to enable benchmark mode")
	}
}

func TestNewOptionsJoinsErrorsAndReturnsZeroOptions(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first option")
	secondErr := errors.New("second option")
	laterApplied := false
	cfg, err := NewOptions([]Option{
		func(opts *Options) error {
			opts.BenchmarkMode = true

			return firstErr
		},
		nil,
		func(*Options) error {
			laterApplied = true

			return secondErr
		},
	})
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected both option errors, got: %v", err)
	}

	if cfg != (Options{}) {
		t.Fatalf("expected zero options on failure, got: %+v", cfg)
	}

	if !laterApplied {
		t.Fatal("expected later options to execute after an error")
	}
}

func TestNewTestingPropagatesOptionErrorsAndReturnsZeroTesting(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("testing option")
	instance, err := NewTesting[*pointerCloser]([]Option{
		func(*Options) error {
			return wantErr
		},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected testing option error, got: %v", err)
	}

	if instance.Benchmark != nil || instance.Options != (Options{}) {
		t.Fatalf("expected zero testing state on failure, got: %+v", instance)
	}
}
