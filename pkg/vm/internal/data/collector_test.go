package data_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type staticIterator struct{}

func (it *staticIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.None, runtime.None, io.EOF
}

func TestNewCollectorDoesNotPanicOnUnsupportedTypes(t *testing.T) {
	assertNoPanic(t, func() {
		if collector := data.NewCollector(bytecode.CollectorTypeAggregate); collector == nil {
			t.Fatal("expected deterministic fallback collector for aggregate type without plan")
		}
	})

	assertNoPanic(t, func() {
		if collector := data.NewCollector(bytecode.CollectorType(255)); collector == nil {
			t.Fatal("expected deterministic fallback collector for unknown type")
		}
	})
}

func TestNewCollectorSafeReturnsErrorsForUnsupportedTypes(t *testing.T) {
	if _, err := data.NewCollectorSafe(bytecode.CollectorTypeAggregate); err == nil {
		t.Fatal("expected aggregate-without-plan error")
	}

	if _, err := data.NewCollectorSafe(bytecode.CollectorType(255)); err == nil {
		t.Fatal("expected unknown collector type error")
	}
}

func TestIteratorUnsupportedOperationsAreNonPanicking(t *testing.T) {
	iterator := data.NewIterator(&staticIterator{})

	assertNoPanic(t, func() {
		_, err := iterator.MarshalJSON()
		if err == nil {
			t.Fatal("expected marshal json error")
		}
	})

	assertNoPanic(t, func() {
		if hash := iterator.Hash(); hash == 0 {
			t.Fatal("expected deterministic non-zero hash")
		}
	})

	assertNoPanic(t, func() {
		copied := iterator.Copy()
		if copied == nil {
			t.Fatal("expected copy fallback value")
		}
	})
}

func TestRegexpCopyDoesNotPanic(t *testing.T) {
	regexp, err := data.NewRegexp(runtime.NewString("foo"))
	if err != nil {
		t.Fatalf("regexp compile failed: %v", err)
	}

	assertNoPanic(t, func() {
		if copied := regexp.Copy(); copied == nil {
			t.Fatal("expected non-nil copied regexp")
		}
	})
}

func assertNoPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	fn()
}

func TestNoopCollectorDeterministicFallback(t *testing.T) {
	collector := data.NewNoopCollector()
	if collector == nil {
		t.Fatal("expected noop collector")
	}

	ctx := context.Background()
	if err := collector.Set(ctx, runtime.NewString("k"), runtime.NewInt(1)); err != nil {
		t.Fatalf("expected noop set success, got %v", err)
	}

	if _, err := collector.Get(ctx, runtime.NewString("missing")); !errors.Is(err, runtime.ErrNotFound) {
		t.Fatalf("expected not found from noop collector get, got %v", err)
	}

	length, err := collector.Length(ctx)
	if err != nil {
		t.Fatalf("expected noop length success, got %v", err)
	}
	if length != 0 {
		t.Fatalf("expected noop length 0, got %d", length)
	}
}
