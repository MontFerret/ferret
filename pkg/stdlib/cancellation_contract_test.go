package stdlib_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/random"
	stdlibstrings "github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestStdlibTraversalsCompleteWhenHostsIgnoreCancellation(t *testing.T) {
	for _, operation := range []string{
		"count", "distinct", "includes", "sum", "filtered-sum", "median", "choice", "shuffle",
		"entries", "values", "from-entries", "zip", "merge-arguments", "key-arguments", "join",
	} {
		for _, cancelAt := range []int{1, 4096} {
			t.Run(fmt.Sprintf("%s/after=%d", operation, cancelAt), func(t *testing.T) {
				ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), rnd.NewSeed(42)))
				defer cancel()
				list := &contextList{
					expected: ctx, cancel: cancel, cancelAt: cancelAt, stage: "visit",
					values: make([]runtime.Value, 4096), calls: make(map[string]int),
				}
				var item runtime.Value = runtime.Int(0)
				switch operation {
				case "filtered-sum":
					item = runtime.None
				case "from-entries":
					item = runtime.NewArrayWith(runtime.String("key"), runtime.Int(0))
				case "zip", "key-arguments", "join":
					item = runtime.String("key")
				case "merge-arguments":
					item = runtime.NewObject()
				}

				for index := range list.values {
					list.values[index] = item
				}

				stream := &contextIterable{Value: runtime.None, list: list}
				object := &contextMap{Map: runtime.NewObject(), list: list}
				var result runtime.Value
				var err error
				switch operation {
				case "count":
					result, err = collections.Count(ctx, stream)
				case "distinct":
					result, err = collections.CountDistinct(ctx, stream)
				case "includes":
					result, err = collections.Includes(ctx, stream, runtime.Int(1))
				case "sum", "filtered-sum":
					result, err = math.Sum(ctx, list)
				case "median":
					result, err = math.Median(ctx, list)
				case "choice":
					result, err = random.Choice(ctx, list)
				case "shuffle":
					result, err = random.Shuffle(ctx, list)
				case "entries":
					result, err = objects.Entries(ctx, object)
				case "values":
					result, err = objects.Values(ctx, object)
				case "from-entries":
					result, err = objects.FromEntries(ctx, stream)
				case "zip":
					list.stage = "At"
					result, err = objects.Zip(ctx, list, list)
				case "merge-arguments":
					result, err = objects.Merge(ctx, list)
				case "key-arguments":
					result, err = objects.KeepKeys(ctx, runtime.NewObject(), list)
				case "join":
					result, err = stdlibstrings.Join(ctx, list, runtime.String(","))
				}

				// Normalizers finish traversal before downstream runtime operations
				// independently reject the canceled context.
				if operation == "merge-arguments" || operation == "key-arguments" {
					if !errors.Is(err, context.Canceled) {
						t.Fatalf("downstream error = %v, want cancellation", err)
					}
				} else if err != nil {
					t.Fatalf("completed operation returned %v", err)
				}

				total := len(list.values)
				if operation == "zip" {
					total *= 2
				}

				if visits := list.calls[list.stage]; visits != total {
					t.Fatalf("visited %d values, want %d", visits, total)
				}

				if list.closes != 0 || list.created != nil && list.created.closes != 0 {
					t.Fatal("successful operation closed borrowed or returned values")
				}

				switch operation {
				case "count":
					if result != runtime.Int(len(list.values)) {
						t.Fatalf("count = %v", result)
					}
				case "distinct":
					if result != runtime.Int(1) {
						t.Fatalf("distinct count = %v", result)
					}
				case "includes":
					if result != runtime.False {
						t.Fatalf("membership = %v", result)
					}
				case "sum", "filtered-sum", "median":
					if result != runtime.Float(0) {
						t.Fatalf("numeric result = %v", result)
					}
				case "choice":
					if result != item {
						t.Fatal("choice lost the original value")
					}
				case "entries", "values", "shuffle":
					destination, ok := result.(runtime.List)
					if !ok {
						t.Fatalf("result type = %T, want List", result)
					}

					size, err := destination.Length(ctx)
					if err != nil || size != runtime.Int(len(list.values)) {
						t.Fatalf("result length = %d, %v", size, err)
					}
				case "from-entries", "zip":
					value, found, err := result.(runtime.Map).Lookup(ctx, runtime.String("key"))
					var want runtime.Value = runtime.Int(0)
					if operation == "zip" {
						want = runtime.String("key")
					}

					if err != nil || !found || value != want {
						t.Fatalf("result entry = %v, %t, %v", value, found, err)
					}
				case "join":
					want := strings.Repeat("key,", len(list.values)-1) + "key"
					if result != runtime.String(want) {
						t.Fatal("join lost values or separators")
					}
				}
			})
		}
	}
}

func TestStdlibMutationsCompleteWhenHostsIgnoreCancellation(t *testing.T) {
	for _, operation := range []string{"reverse-read", "reverse-append", "shuffle-append", "shuffle-swap", "remove-read", "remove-tail"} {
		for _, cancelAt := range []int{1, 4095} {
			t.Run(fmt.Sprintf("%s/after=%d", operation, cancelAt), func(t *testing.T) {
				ctx, cancel := context.WithCancel(rnd.WithContext(t.Context(), rnd.NewSeed(42)))
				defer cancel()
				cleanup := errors.New("destination cleanup")
				list := &contextList{
					expected: ctx, cancel: cancel, cancelAt: cancelAt, cleanup: cleanup,
					values: make([]runtime.Value, 4096), calls: make(map[string]int),
				}
				for index := range list.values {
					list.values[index] = runtime.Int(0)
				}

				var result runtime.Value
				var err error
				switch operation {
				case "reverse-read", "reverse-append":
					list.stage = "At"
					if operation == "reverse-append" {
						list.stage = "Append"
					}

					result, err = collections.Reverse(ctx, list)
				case "shuffle-append", "shuffle-swap":
					list.stage = "Append"
					if operation == "shuffle-swap" {
						list.stage = "Swap"
					}

					result, err = random.Shuffle(ctx, list)
				case "remove-read", "remove-tail":
					list.stage = "At"
					if operation == "remove-tail" {
						list.stage = "RemoveAt"
					}

					result, err = arrays.RemoveMutable(ctx, list, runtime.Int(0))
				}

				if err != nil || list.closes != 0 {
					t.Fatalf("error %v, source closes %d", err, list.closes)
				}

				target := list
				total := 4096
				if list.created != nil {
					if result != list.created || list.created.closes != 0 || len(list.created.values) != total {
						t.Fatal("completed destination ownership or contents changed")
					}

					if list.stage != "At" {
						target = list.created
					}
				} else if result != list || len(list.values) != 0 {
					t.Fatal("removal did not finish on the original list")
				}

				if list.stage == "Swap" {
					total--
				}

				if visits := target.calls[list.stage]; visits != total {
					t.Fatalf("completed %d operations, want %d", visits, total)
				}
			})
		}
	}
}

func TestShortStdlibWrappersPropagateHostContextAndErrors(t *testing.T) {
	for _, failure := range []error{errors.New("host failed"), context.Canceled, context.DeadlineExceeded} {
		for _, operation := range []string{"has-key", "array-set", "array-push"} {
			t.Run(operation+"/"+failure.Error(), func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				list := &contextList{
					expected: ctx, cancel: cancel, cancelAt: 1, failure: failure,
					values: []runtime.Value{runtime.Int(0)}, calls: make(map[string]int),
				}
				var result runtime.Value
				var err error
				switch operation {
				case "has-key":
					list.stage = "ContainsKey"
					result, err = objects.HasKey(ctx, &contextMap{Map: runtime.NewObject(), list: list}, runtime.String("key"))
				case "array-set":
					list.stage = "SetAt"
					result, err = arrays.SetMutable(ctx, list, runtime.Int(0), runtime.Int(1))
				case "array-push":
					list.stage = "Append"
					result, err = arrays.PushMutable(ctx, list, runtime.Int(1))
				}

				if result != runtime.None || !errors.Is(err, failure) || list.calls[list.stage] != 1 {
					t.Fatalf("result %v, error %v; host error/context lost", result, err)
				}
			})
		}
	}
}
