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

var collectionOperations = map[string]func(context.Context, runtime.Value) (runtime.Value, error){
	"choice":  Choice,
	"shuffle": Shuffle,
}

func TestRandomCollectionValues(t *testing.T) {
	for _, test := range []struct {
		order  []int
		size   int
		choice int
	}{
		{size: 0, choice: -1, order: []int{}},
		{size: 1, choice: 0, order: []int{0}},
		{size: 2, choice: 0, order: []int{0, 1}},
		{size: 4, choice: 0, order: []int{0, 1, 2, 3}},
		{size: 12, choice: 4, order: []int{3, 8, 2, 4, 7, 5, 6, 9, 0, 1, 11, 10}},
	} {
		for name, call := range collectionOperations {
			for _, host := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/%d/host=%t", name, test.size, host), func(t *testing.T) {
					values := collectionValues()[:test.size]
					original := append([]runtime.Value(nil), values...)
					var input runtime.List = runtime.NewArrayOf(values)
					list := &traversalList{values: values}
					if host {
						input = list
					}

					ctx := rnd.WithContext(t.Context(), rnd.NewSeed(42))
					got, err := call(ctx, input)
					if err != nil {
						t.Fatal(err)
					}

					if name == "choice" {
						var want runtime.Value = runtime.None
						if test.choice >= 0 {
							want = original[test.choice]
						}

						if !sameCollectionValue(got, want) {
							t.Fatal("choice did not return the expected original value")
						}
					} else {
						result, ok := got.(*runtime.Array)
						if !ok || result == input {
							t.Fatalf("shuffle returned %T instead of a fresh Array", got)
						}

						shuffled := arrayValues(t, result)
						want := make([]runtime.Value, len(test.order))
						for i, index := range test.order {
							want[i] = original[index]
						}

						assertCollectionValues(t, shuffled, want)
						assertCollectionPermutation(t, shuffled, original)

						if len(shuffled) == 0 {
							err = result.Append(ctx, runtime.True)
						} else {
							err = result.SetAt(ctx, 0, runtime.False)
						}

						if err != nil {
							t.Fatal(err)
						}
					}

					if host {
						assertCollectionValues(t, list.values, original)
						if list.calls != 1 || list.cursorCloses != 1 || list.closes != 0 {
							t.Fatalf("traversals %d, iterator closes %d, source closes %d", list.calls, list.cursorCloses, list.closes)
						}
					} else {
						assertCollectionValues(t, arrayValues(t, input.(*runtime.Array)), original)
					}

					for _, value := range original {
						if borrowed, ok := value.(*borrowedValue); ok && borrowed.closes != 0 {
							t.Fatal("closed a borrowed value")
						}
					}
				})
			}
		}
	}
}

func TestRandomCollectionTrivialCallsDoNotDraw(t *testing.T) {
	for name, call := range collectionOperations {
		for size := 0; size <= 1; size++ {
			t.Run(fmt.Sprintf("%s/%d", name, size), func(t *testing.T) {
				src, control := rnd.NewSeed(42), rnd.NewSeed(42)
				ctx := rnd.WithContext(t.Context(), src)
				for range 3 {
					if _, err := call(ctx, &traversalList{values: collectionValues()[:size]}); err != nil {
						t.Fatal(err)
					}

					if src.Float64() != control.Float64() {
						t.Fatal("deterministic operation advanced the source")
					}
				}
			})
		}
	}
}

func TestRandomCollectionValidation(t *testing.T) {
	for name, call := range collectionOperations {
		t.Run(name, func(t *testing.T) {
			for _, value := range []runtime.Value{runtime.None, runtime.Int(1), runtime.Float(1), runtime.True, runtime.String("a"), runtime.NewObject()} {
				src, control := rnd.NewSeed(42), rnd.NewSeed(42)
				got, err := call(rnd.WithContext(t.Context(), src), value)
				position, ok, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !ok || position != 0 {
					t.Fatalf("invalid argument: result %v, error %v", got, err)
				}

				if src.Float64() != control.Float64() {
					t.Fatal("argument validation advanced the source")
				}
			}

			for size := 0; size <= 2; size++ {
				list := &traversalList{values: collectionValues()[:size]}
				got, err := call(t.Context(), list)
				if got != runtime.None || !errors.Is(err, runtime.ErrUnexpected) || list.calls != 0 {
					t.Fatalf("missing source: result %v, error %v, traversals %d", got, err, list.calls)
				}
			}
		})
	}
}

func TestRandomCollectionTraversalErrors(t *testing.T) {
	primary, cleanup := errors.New("host traversal failed"), errors.New("iterator cleanup failed")
	for name, call := range collectionOperations {
		for _, failAt := range []int{0, 1, 3} {
			for _, closeErr := range []error{nil, cleanup} {
				t.Run(fmt.Sprintf("%s/after=%d/cleanup=%t", name, failAt, closeErr != nil), func(t *testing.T) {
					src, control := rnd.NewSeed(42), rnd.NewSeed(42)
					ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), src))
					defer cancel()
					value := &borrowedValue{}
					list := &traversalList{values: []runtime.Value{value, value, value}, closeErr: closeErr}
					list.visit = func(_ context.Context, intIndex int) error {
						if intIndex == failAt {
							cancel()

							return primary
						}

						return nil
					}

					got, err := call(ctx, list)
					if got != runtime.None || err != list.traversalErr || !errors.Is(err, primary) || errors.Is(err, cleanup) != (closeErr != nil) {
						t.Fatalf("traversal error lost: result %v, error %v", got, err)
					}

					if list.calls != 1 || list.cursorCloses != 1 || list.closes != 0 || value.closes != 0 {
						t.Fatal("incorrect traversal or borrowed-resource cleanup")
					}

					if name == "choice" {
						for n := 2; n <= failAt; n++ {
							control.Int64(0, int64(n-1))
						}
					}

					if src.Float64() != control.Float64() {
						t.Fatal("failed traversal consumed an unexpected number of draws")
					}
				})
			}
		}

		t.Run(name+"/cleanup-only", func(t *testing.T) {
			list := &traversalList{values: []runtime.Value{&borrowedValue{}}, closeErr: cleanup}
			got, err := call(rnd.WithContext(t.Context(), rnd.NewSeed(42)), list)
			if got != runtime.None || err != list.traversalErr || !errors.Is(err, cleanup) || list.cursorCloses != 1 || list.closes != 0 {
				t.Fatalf("cleanup error lost: result %v, error %v", got, err)
			}
		})
	}
}

func TestRandomCollectionCancellation(t *testing.T) {
	for name, call := range collectionOperations {
		for _, stage := range []string{"before", "first callback", "later callback", "after traversal", "empty traversal"} {
			t.Run(name+"/"+stage, func(t *testing.T) {
				src, control := rnd.NewSeed(42), rnd.NewSeed(42)
				ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), src))
				defer cancel()
				list := &traversalList{values: collectionValues()[:3]}
				cancelAt := 0
				switch stage {
				case "before":
					cancel()
				case "later callback":
					cancelAt = 2
				case "after traversal":
					cancelAt = 3
				case "empty traversal":
					list.values = nil
				}

				list.visit = func(_ context.Context, intIndex int) error {
					if intIndex == cancelAt {
						cancel()
					}

					return nil
				}

				got, err := call(ctx, list)
				if got != runtime.None || !errors.Is(err, context.Canceled) || list.closes != 0 {
					t.Fatalf("cancellation: result %v, error %v", got, err)
				}

				wantTraversals := 1
				if stage == "before" {
					wantTraversals = 0
				}

				if list.calls != wantTraversals || list.cursorCloses != wantTraversals {
					t.Fatalf("traversals %d, iterator closes %d", list.calls, list.cursorCloses)
				}

				if name == "choice" {
					for n := 2; n <= cancelAt; n++ {
						control.Int64(0, int64(n-1))
					}
				}

				if src.Float64() != control.Float64() {
					t.Fatal("canceled operation continued drawing")
				}
			})
		}
	}
}

func collectionValues() []runtime.Value {
	borrowed := &borrowedValue{}

	return []runtime.Value{
		borrowed, runtime.Int(1<<53 + 1), runtime.Float(math.Copysign(0, -1)), runtime.None,
		runtime.String("browser"), runtime.True, runtime.NewArrayWith(runtime.Int(5)), runtime.NewObject(),
		borrowed, runtime.Float(math.Float64frombits(0x7ff8000000000042)), runtime.Int(7), runtime.Float(7),
	}
}

func arrayValues(t *testing.T, array *runtime.Array) []runtime.Value {
	t.Helper()

	var values []runtime.Value
	err := array.ForEach(t.Context(), func(_ context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		values = append(values, value)

		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return values
}

func sameCollectionValue(got, want runtime.Value) bool {
	if floating, ok := want.(runtime.Float); ok {
		actual, ok := got.(runtime.Float)

		return ok && math.Float64bits(float64(actual)) == math.Float64bits(float64(floating))
	}

	return got == want
}

func assertCollectionValues(t *testing.T, got, want []runtime.Value) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("value count %d, want %d", len(got), len(want))
	}

	for i := range want {
		if !sameCollectionValue(got[i], want[i]) {
			t.Fatalf("value %d lost its original identity or representation", i)
		}
	}
}

func assertCollectionPermutation(t *testing.T, got, original []runtime.Value) {
	t.Helper()

	if len(got) != len(original) {
		t.Fatalf("permutation length %d, want %d", len(got), len(original))
	}

	matched := make([]bool, len(original))
	for _, value := range got {
		found := false
		for i, candidate := range original {
			if !matched[i] && sameCollectionValue(value, candidate) {
				matched[i], found = true, true

				break
			}
		}

		if !found {
			t.Fatal("shuffle changed a value or its multiplicity")
		}
	}
}
