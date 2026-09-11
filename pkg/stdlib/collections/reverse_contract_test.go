package collections_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func TestReversePreservesListImplementationAndOwnership(t *testing.T) {
	for _, size := range []int{0, 1, 3} {
		ctx := t.Context()
		resource := &scanValue{}
		values := []runtime.Value{resource, runtime.NewObject(), runtime.Int(3)}[:size]
		source := &configuredList{Array: runtime.NewArrayWith(values...), backend: &listBackend{encoding: "binary", limit: 100}, expected: ctx, calls: make(map[string]int)}
		out, err := collections.Reverse(ctx, source)
		if err != nil {
			t.Fatal(err)
		}

		destination, ok := out.(*configuredList)
		if !ok || destination == source || destination != source.created || destination.backend != source.backend {
			t.Fatalf("factory identity/configuration lost: %T", out)
		}

		length, err := destination.Length(ctx)
		if err != nil || length != runtime.Int(size) {
			t.Fatalf("length: %d, %v", length, err)
		}

		for i, value := range values {
			original, err := source.At(ctx, runtime.Int(i))
			if err != nil || original != value {
				t.Fatal("source changed")
			}

			reversed, err := destination.At(ctx, runtime.Int(size-i-1))
			if err != nil || reversed != value {
				t.Fatal("element identity/order changed")
			}
		}

		if err := destination.Append(ctx, runtime.True); err != nil {
			t.Fatal(err)
		}

		length, err = source.Length(ctx)
		if err != nil || length != runtime.Int(size) {
			t.Fatal("destination aliases source storage")
		}

		if source.calls["New"] != 1 || source.closes != 0 || destination.closes != 0 || resource.closes+resource.copies+resource.clones != 0 {
			t.Fatal("incorrect ownership transfer")
		}

		if err := destination.Close(); err != nil {
			t.Fatal(err)
		}

		if destination.closes != 1 || source.closes != 0 || resource.closes != 0 {
			t.Fatal("caller cleanup changed borrowed ownership")
		}
	}
}

func TestReverseClosesPopulatedDestination(t *testing.T) {
	primary, cleanup := errors.New("operation"), errors.New("close")
	for _, operation := range []string{"At", "Append"} {
		for _, failAt := range []int{1, 2, 3} {
			t.Run(fmt.Sprintf("%s/%d", operation, failAt), func(t *testing.T) {
				resource := &scanValue{}
				source := &configuredList{
					Array:    runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), resource),
					expected: t.Context(), calls: make(map[string]int),
					failures: map[string]error{operation: primary},
					failAt:   map[string]int{operation: failAt}, closeErr: cleanup,
				}
				out, err := collections.Reverse(t.Context(), source)
				if out != runtime.None || !errors.Is(err, primary) || !errors.Is(err, cleanup) {
					t.Fatalf("lost failure: %v, %v", out, err)
				}

				if source.created == nil || source.created.closes != 1 || source.closes != 0 || resource.closes != 0 {
					t.Fatal("incorrect partial-result ownership")
				}

				length, err := source.created.Array.Length(t.Context())
				if err != nil || length != runtime.Int(failAt-1) {
					t.Fatalf("unexpected partial result: %d, %v", length, err)
				}
			})
		}
	}
}

func TestReverseRejectsNegativeHostLengthBeforeConstruction(t *testing.T) {
	for _, size := range []runtime.Int{-1, math.MinInt64} {
		source := &configuredList{Array: runtime.NewArray(0), expected: t.Context(), calls: make(map[string]int), length: &size}
		result, err := collections.Reverse(t.Context(), source)
		if result != runtime.None || !errors.Is(err, runtime.ErrInvalidOperation) || source.calls["New"] != 0 || source.closes != 0 {
			t.Fatalf("negative length: %v, %v", result, err)
		}
	}
}

func TestReverseFailureOwnership(t *testing.T) {
	primary, cleanup := errors.New("operation"), errors.New("cleanup")
	for _, stage := range []string{"Length", "New", "At", "Append", "pre-cancel", "cancel-Length", "cancel-New", "cancel-At", "cancel-Append"} {
		for _, closeFails := range []bool{false, true} {
			t.Run(stage+map[bool]string{false: "", true: "+cleanup"}[closeFails], func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				value := &scanValue{}
				source := &configuredList{Array: runtime.NewArrayWith(value), backend: &listBackend{}, expected: ctx, calls: make(map[string]int), failures: make(map[string]error)}
				want := primary
				constructed := stage == "At" || stage == "Append" || stage == "cancel-New" || stage == "cancel-At" || stage == "cancel-Append"
				if closeFails {
					source.closeErr, source.newCleanupErr = cleanup, cleanup
				}

				if stage == "pre-cancel" {
					cancel()
					want = context.Canceled
				} else if len(stage) > 7 && stage[:7] == "cancel-" {
					want = context.Canceled
					source.after = func(name string) {
						if name == stage[7:] {
							cancel()
						}
					}
				} else {
					source.failures[stage] = primary
				}

				out, err := collections.Reverse(ctx, source)
				if out != runtime.None || !errors.Is(err, want) {
					t.Fatalf("got %v, %v; want %v", out, err, want)
				}

				if closeFails && (constructed || stage == "New") && !errors.Is(err, cleanup) {
					t.Fatalf("cleanup error lost: %v", err)
				}

				if constructed {
					if source.created == nil || source.created.closes != 1 {
						t.Fatal("incomplete destination not closed exactly once")
					}
				} else if source.created != nil {
					t.Fatal("unexpected destination")
				}

				if source.closes != 0 || value.closes+value.clones+value.copies != 0 {
					t.Fatal("borrowed source/element consumed")
				}

				if stage == "New" && source.newRollbacks != 1 {
					t.Fatal("factory did not own failed construction")
				}

				if stage == "pre-cancel" && len(source.calls) != 0 {
					t.Fatal("pre-canceled call performed host work")
				}

				if stage == "cancel-Length" && source.calls["New"] != 0 {
					t.Fatal("constructed after cancellation")
				}

				if stage == "cancel-New" && source.calls["At"] != 0 {
					t.Fatal("read after cancellation")
				}

				if stage == "cancel-At" && source.created.calls["Append"] != 0 {
					t.Fatal("appended after cancellation")
				}
			})
		}
	}
}
