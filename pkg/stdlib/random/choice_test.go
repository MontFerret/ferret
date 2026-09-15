package random

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestReservoirExhaustive(t *testing.T) {
	values := []runtime.Value{runtime.Int(0), runtime.Int(1), runtime.Int(2), runtime.Int(3)}
	var selections [4]int
	for second := int64(0); second < 2; second++ {
		for third := int64(0); third < 3; third++ {
			for fourth := int64(0); fourth < 4; fourth++ {
				draws := []int64{second, third, fourth}
				calls := 0
				list := &traversalList{values: values, indices: []runtime.Int{99, 99, -4, 800}}
				got, err := choose(t.Context(), list, func(min, max int64) int64 {
					if calls >= len(draws) || min != 0 || max != int64(calls+1) {
						t.Fatalf("draw %d requested [%d, %d]", calls, min, max)
					}

					value := draws[calls]
					calls++

					return value
				})
				if err != nil {
					t.Fatal(err)
				}

				want := runtime.Int(0)
				switch {
				case fourth == 0:
					want = 3
				case third == 0:
					want = 2
				case second == 0:
					want = 1
				}

				if got != want || calls != 3 || list.calls != 1 {
					t.Fatalf("draws %v: got %v, want %v; draws %d, traversals %d", draws, got, want, calls, list.calls)
				}

				selections[got.(runtime.Int)]++
			}
		}
	}

	// All 2*3*4 draw combinations are equally likely; each element wins six.
	if selections != [4]int{6, 6, 6, 6} {
		t.Fatalf("selection counts = %v", selections)
	}
}

func TestReservoirPreservesEveryValue(t *testing.T) {
	values := collectionValues()
	for selected, want := range values {
		got, err := choose(t.Context(), &traversalList{values: values}, func(_, max int64) int64 {
			if max == int64(selected) {
				return 0
			}

			return max
		})
		if err != nil || !sameCollectionValue(got, want) {
			t.Fatalf("element %d was not preserved (error %v)", selected, err)
		}
	}
}

func TestReservoirCompletesAfterCancellationDuringDraw(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	list := &traversalList{values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3)}}
	calls := 0
	got, err := choose(ctx, list, func(_, _ int64) int64 {
		calls++
		cancel()

		return 0
	})
	if got != runtime.Int(3) || err != nil || calls != 2 || list.cursorCloses != 1 {
		t.Fatalf("cancellation: result %v, error %v, draws %d, closes %d", got, err, calls, list.cursorCloses)
	}
}
