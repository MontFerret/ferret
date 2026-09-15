package random

import (
	"context"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestShuffleDestinationFailures(t *testing.T) {
	primary, cleanup := errors.New("backend failed"), errors.New("destination close failed")
	negative, minimum := runtime.Int(-1), runtime.Int(math.MinInt64)
	for _, test := range []struct {
		length *runtime.Int
		name   string
		op     string
		failAt int
	}{
		{name: "factory", op: "New", failAt: 1},
		{name: "first append", op: "Append", failAt: 1},
		{name: "partial append", op: "Append", failAt: 2},
		{name: "length", op: "Length", failAt: 1},
		{name: "negative length", op: "Length", failAt: 1, length: &negative},
		{name: "minimum length", op: "Length", failAt: 1, length: &minimum},
		{name: "first swap", op: "Swap", failAt: 1},
		{name: "later swap", op: "Swap", failAt: 2},
	} {
		for _, closeErr := range []error{nil, cleanup} {
			for _, cancelAlongside := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/cleanup=%t/cancel=%t", test.name, closeErr != nil, cancelAlongside), func(t *testing.T) {
					rng, control := rnd.NewSeed(42), rnd.NewSeed(42)
					ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), rng))
					defer cancel()
					borrowed := &borrowedValue{}
					values := []runtime.Value{borrowed, runtime.Int(2), runtime.Int(3), borrowed}
					original := append([]runtime.Value(nil), values...)
					source := newShuffleList(&traversalList{values: values})
					source.expected, source.cleanupErr = ctx, closeErr
					source.failOp, source.failAt, source.failure = test.op, test.failAt, primary
					source.length = test.length
					wantErr := primary
					if test.length != nil {
						source.failOp = ""
						wantErr = runtime.ErrInvalidOperation
					}

					calls := 0
					source.after = func(operation string) {
						if operation == test.op {
							calls++
							if cancelAlongside && calls == test.failAt {
								cancel()
							}
						}
					}

					got, err := Shuffle(ctx, source)
					if got != runtime.None || !errors.Is(err, wantErr) || errors.Is(err, cleanup) != (closeErr != nil) || errors.Is(err, context.Canceled) {
						t.Fatalf("operation error lost: result %T, error %v", got, err)
					}

					if source.operations["New"] != 1 || source.created == nil || source.created.closes != 1 || source.closes != 0 || borrowed.closes != 0 {
						t.Fatal("incorrect factory, destination, or borrowed-value ownership")
					}

					wantTraversals, wantRollbacks := 1, 0
					wantAppends, wantLengths, wantSwaps := 4, 1, 0
					switch test.op {
					case "New":
						wantTraversals, wantRollbacks = 0, 1
						wantAppends, wantLengths = 0, 0
					case "Append":
						wantAppends, wantLengths = test.failAt, 0
					case "Swap":
						wantSwaps = test.failAt
					}

					if source.calls != wantTraversals || source.cursorCloses != wantTraversals || source.rollbacks != wantRollbacks {
						t.Fatalf("traversals %d, iterator closes %d, factory rollbacks %d", source.calls, source.cursorCloses, source.rollbacks)
					}

					for operation, count := range map[string]int{"Append": wantAppends, "Length": wantLengths, "Swap": wantSwaps} {
						if source.created.operations[operation] != count {
							t.Fatalf("%s calls = %d, want %d", operation, source.created.operations[operation], count)
						}
					}

					for i := 0; i < wantSwaps; i++ {
						control.Int64(0, int64(3-i))
					}

					if rng.Float64() != control.Float64() {
						t.Fatal("failure consumed an unexpected number of draws")
					}

					assertCollectionValues(t, source.values, original)
				})
			}
		}
	}
}

func TestShuffleCancellationAfterHostOperations(t *testing.T) {
	cleanup := errors.New("destination close failed")
	for _, test := range []struct {
		name string
		op   string
		at   int
		size int
	}{
		{name: "before factory", size: 4},
		{name: "factory", op: "New", at: 1, size: 4},
		{name: "first append", op: "Append", at: 1, size: 4},
		{name: "last append", op: "Append", at: 4, size: 4},
		{name: "length", op: "Length", at: 1, size: 4},
		{name: "empty length", op: "Length", at: 1, size: 0},
		{name: "singleton length", op: "Length", at: 1, size: 1},
		{name: "first swap", op: "Swap", at: 1, size: 4},
		{name: "last swap", op: "Swap", at: 3, size: 4},
	} {
		for _, closeErr := range []error{nil, cleanup} {
			t.Run(fmt.Sprintf("%s/cleanup=%t", test.name, closeErr != nil), func(t *testing.T) {
				rng, control := rnd.NewSeed(42), rnd.NewSeed(42)
				ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), rng))
				defer cancel()
				source := newShuffleList(&traversalList{values: collectionValues()[:test.size]})
				source.expected, source.cleanupErr = ctx, closeErr
				calls := 0
				source.after = func(operation string) {
					if operation == test.op {
						calls++
						if calls == test.at {
							cancel()
						}
					}
				}

				if test.op == "" {
					cancel()
				}

				got, err := Shuffle(ctx, source)
				if got != runtime.None || !errors.Is(err, context.Canceled) || errors.Is(err, cleanup) != (closeErr != nil && test.op != "") {
					t.Fatalf("cancellation lost: result %T, error %v", got, err)
				}

				if test.op == "" {
					if source.created != nil || source.operations["New"] != 0 {
						t.Fatal("pre-canceled shuffle called the factory")
					}
				} else if source.created == nil || source.created.closes != 1 || source.rollbacks != 0 {
					t.Fatal("cancellation did not close the successfully constructed destination exactly once")
				}

				if source.closes != 0 {
					t.Fatal("cancellation closed the source")
				}

				for _, value := range source.values {
					if borrowed, ok := value.(*borrowedValue); ok && borrowed.closes != 0 {
						t.Fatal("cancellation closed a borrowed value")
					}
				}

				if test.op == "New" && source.calls != 0 {
					t.Fatal("traversal started after factory cancellation")
				}

				if test.op == "Append" && (source.created.operations["Append"] != test.at || source.created.operations["Length"] != 0) {
					t.Fatal("materialization continued after append cancellation")
				}

				if test.op == "Swap" {
					for i := 0; i < test.at; i++ {
						control.Int64(0, int64(test.size-1-i))
					}
				}

				if rng.Float64() != control.Float64() {
					t.Fatal("cancellation consumed an unexpected number of draws")
				}
			})
		}
	}
}

func TestShufflePreservesTraversalAndBothCleanupErrors(t *testing.T) {
	primary := errors.New("source traversal failed")
	iteratorCleanup, destinationCleanup := errors.New("iterator close failed"), errors.New("destination close failed")
	borrowed := &borrowedValue{}
	source := newShuffleList(&traversalList{
		values:   []runtime.Value{borrowed, borrowed},
		closeErr: iteratorCleanup,
		visit: func(_ context.Context, index int) error {
			if index == 2 {
				return primary
			}

			return nil
		},
	})
	source.cleanupErr = destinationCleanup
	rng, control := rnd.NewSeed(42), rnd.NewSeed(42)
	got, err := Shuffle(rnd.WithContext(t.Context(), rng), source)
	if got != runtime.None || !errors.Is(err, primary) || !errors.Is(err, iteratorCleanup) || !errors.Is(err, destinationCleanup) {
		t.Fatalf("joined failure lost: result %T, error %v", got, err)
	}

	if source.created.closes != 1 || source.cursorCloses != 1 || source.closes != 0 || borrowed.closes != 0 {
		t.Fatal("incorrect cleanup for source traversal failure")
	}

	if rng.Float64() != control.Float64() {
		t.Fatal("failed traversal consumed randomness")
	}
}
