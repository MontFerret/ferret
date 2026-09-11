package arrays_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func TestRemovalComparisonContract(t *testing.T) {
	for _, mode := range []string{"immutable", "mutable", "legacy"} {
		for _, limit := range []runtime.Int{math.MinInt64, -5, -1, 0, 1, 2} {
			if mode != "legacy" && limit != -1 {
				continue
			}

			for _, failure := range []error{nil, errors.New("late equality failure"), context.Canceled} {
				t.Run(mode+"/"+limit.String(), func(t *testing.T) {
					ctx := t.Context()
					target := &removalComparisonValue{}
					var seen []int
					first := &removalComparisonValue{target: target, expectedContext: ctx, seen: &seen, index: 0, equal: true}
					second := &removalComparisonValue{target: target, expectedContext: ctx, seen: &seen, index: 1}
					last := &removalComparisonValue{target: target, expectedContext: ctx, seen: &seen, index: 2, equal: true, failure: failure}
					source := runtime.NewArrayWith(first, second, last)
					host := &removalFilterList{List: source, expectedContext: ctx}

					var result runtime.Value
					var err error
					switch mode {
					case "immutable":
						result, err = arrays.Remove(ctx, host, target)
					case "mutable":
						result, err = arrays.RemoveMutable(ctx, host, target)
					case "legacy":
						result, err = callLegacy("remove_value", ctx, host, target, limit)
					}

					if !reflect.DeepEqual(seen, []int{0, 1, 2}) {
						t.Fatalf("comparison traversal changed: %v", seen)
					}

					if mode != "mutable" {
						if host.filters != 1 {
							t.Fatalf("host Filter calls=%d", host.filters)
						}

						for index, want := range []runtime.Value{first, second, last} {
							got, accessErr := source.At(ctx, runtime.Int(index))
							if accessErr != nil || got != want {
								t.Fatalf("immutable source changed at %d: %v, %v", index, got, accessErr)
							}
						}
					} else if host.filters != 0 {
						t.Fatal("mutable removal invoked immutable filtering")
					}

					if failure != nil {
						if !errors.Is(err, failure) {
							t.Fatalf("lost late comparison error: %v", err)
						}

						if mode == "mutable" && result != runtime.None {
							t.Fatalf("mutable failure result=%v", result)
						}

						return
					}

					if err != nil {
						t.Fatal(err)
					}

					resultList, ok := result.(runtime.List)
					if !ok {
						t.Fatalf("result is not a list: %T", result)
					}

					if mode == "mutable" && result != host {
						t.Fatal("mutable target identity changed")
					}

					if mode != "mutable" && (result == host || result == source) {
						t.Fatal("immutable removal returned the source")
					}

					want := []runtime.Value{second}
					if limit == 0 {
						want = []runtime.Value{first, second, last}
					} else if limit == 1 {
						want = []runtime.Value{second, last}
					}

					size, err := resultList.Length(ctx)
					if err != nil || int(size) != len(want) {
						t.Fatalf("result length=%d, want=%d err=%v", size, len(want), err)
					}

					for index, expected := range want {
						value, err := resultList.At(ctx, runtime.Int(index))
						if err != nil || value != expected {
							t.Fatalf("survivor %d changed identity/order: %v, %v", index, value, err)
						}
					}
				})
			}
		}
	}
}
