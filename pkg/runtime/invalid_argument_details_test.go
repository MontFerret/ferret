package runtime_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestInvalidArgumentDetails(t *testing.T) {
	cause := runtime.TypeErrorOf(runtime.True, runtime.TypeString)
	err := runtime.ArgError(cause, 1)

	if !errors.Is(err, runtime.ErrInvalidArgument) {
		t.Fatalf("expected errors.Is(err, ErrInvalidArgument) to be true")
	}

	if !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("expected errors.Is(err, ErrInvalidType) to be true")
	}

	pos, ok, wrappedCause := runtime.InvalidArgumentDetails(err)
	if !ok {
		t.Fatal("expected invalid argument details")
	}

	if got, want := pos, 1; got != want {
		t.Fatalf("unexpected invalid argument position: got %d, want %d", got, want)
	}

	if wrappedCause == nil {
		t.Fatal("expected wrapped cause")
	}

	if !errors.Is(wrappedCause, runtime.ErrInvalidType) {
		t.Fatalf("expected wrapped cause to preserve invalid type classification")
	}

	if got, want := err.Error(), "invalid argument at position 2 - "+cause.Error(); got != want {
		t.Fatalf("unexpected invalid argument error string: got %q, want %q", got, want)
	}
}

func TestInvalidArgumentDetailsReturnsOutermostArgument(t *testing.T) {
	inner := runtime.ArgError(runtime.TypeErrorOf(runtime.True, runtime.TypeString), 1)
	outer := runtime.ArgError(inner, 0)

	pos, ok, cause := runtime.InvalidArgumentDetails(outer)
	if !ok {
		t.Fatal("expected invalid argument details")
	}

	if got, want := pos, 0; got != want {
		t.Fatalf("unexpected outer invalid argument position: got %d, want %d", got, want)
	}

	if cause == nil {
		t.Fatal("expected outer invalid argument cause")
	}

	if got, want := cause.Error(), inner.Error(); got != want {
		t.Fatalf("unexpected outer invalid argument cause: got %q, want %q", got, want)
	}
}
