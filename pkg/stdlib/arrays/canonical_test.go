package arrays_test

import (
	"context"
	"errors"
	"math"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func TestCanonicalArrayRegistrations(t *testing.T) {
	functions := arrayFunctions(t)
	for _, test := range []struct {
		got  []string
		want []string
	}{
		{functions.A1().Names(), []string{"first", "flatten", "last", "sorted", "unique"}},
		{functions.A2().Names(), []string{"append", "at", "contains", "flatten", "index_of", "remove", "remove_any", "remove_at", "slice"}},
		{functions.A3().Names(), []string{"slice"}},
		{functions.Var().Names(), []string{"concat", "difference", "intersection", "symmetric_difference", "union"}},
	} {
		var names []string
		for _, name := range test.got {
			if strings.HasPrefix(name, "arrays::") {
				names = append(names, strings.TrimPrefix(name, "arrays::"))
			}
		}

		slices.Sort(names)
		if !reflect.DeepEqual(names, test.want) {
			t.Fatalf("canonical registrations = %v, want %v", names, test.want)
		}
	}

	for _, test := range [][2]string{{"union", "arrays::concat"}, {"union_distinct", "arrays::union"}, {"minus", "arrays::difference"}, {"intersection", "arrays::intersection"}} {
		global, _ := functions.Var().Get(test[0])
		canonical, _ := functions.Var().Get(test[1])
		if reflect.ValueOf(global).Pointer() != reflect.ValueOf(canonical).Pointer() {
			t.Fatalf("%s must share the canonical implementation", test[0])
		}
	}

	if functions.A0().Size() != 0 || functions.A4().Size() != 0 {
		t.Fatal("unexpected array overloads")
	}
}

func TestCanonicalArrayOperations(t *testing.T) {
	input := runtime.NewArrayWith(runtime.Int(3), runtime.Int(1), runtime.Int(3), runtime.Int(2))
	other := runtime.NewArrayWith(runtime.Int(2), runtime.Int(4), runtime.Int(4))
	for _, test := range []struct {
		name string
		run  func() (runtime.Value, error)
		want string
	}{
		{name: "first", run: func() (runtime.Value, error) { return arrays.First(t.Context(), input) }, want: "3"},
		{name: "last", run: func() (runtime.Value, error) { return arrays.Last(t.Context(), input) }, want: "2"},
		{name: "at", run: func() (runtime.Value, error) { return arrays.At(t.Context(), input, runtime.Int(1)) }, want: "1"},
		{name: "append", run: func() (runtime.Value, error) { return arrays.Append(t.Context(), input, runtime.Int(3)) }, want: "[3,1,3,2,3]"},
		{name: "append array as element", run: func() (runtime.Value, error) { return arrays.Append(t.Context(), input, other) }, want: "[3,1,3,2,[2,4,4]]"},
		{name: "concat", run: func() (runtime.Value, error) { return arrays.Concat(t.Context(), input, other) }, want: "[3,1,3,2,2,4,4]"},
		{name: "slice", run: func() (runtime.Value, error) { return arrays.Slice(t.Context(), input, runtime.Int(1), runtime.Int(2)) }, want: "[1,3]"},
		{name: "flatten", run: func() (runtime.Value, error) { return arrays.Flatten(t.Context(), runtime.NewArrayWith(input, other)) }, want: "[3,1,3,2,2,4,4]"},
		{name: "contains", run: func() (runtime.Value, error) { return arrays.Contains(t.Context(), input, runtime.Float(3)) }, want: "true"},
		{name: "index_of", run: func() (runtime.Value, error) { return arrays.IndexOf(t.Context(), input, runtime.Int(3)) }, want: "0"},
		{name: "missing index", run: func() (runtime.Value, error) { return arrays.IndexOf(t.Context(), input, runtime.Int(99)) }, want: "-1"},
		{name: "unique", run: func() (runtime.Value, error) { return arrays.Unique(t.Context(), input) }, want: "[3,1,2]"},
		{name: "sorted", run: func() (runtime.Value, error) { return arrays.Sorted(t.Context(), input) }, want: "[1,2,3,3]"},
		{name: "remove", run: func() (runtime.Value, error) { return arrays.Remove(t.Context(), input, runtime.Int(3)) }, want: "[1,2]"},
		{name: "remove_at", run: func() (runtime.Value, error) { return arrays.RemoveAt(t.Context(), input, runtime.Int(2)) }, want: "[3,1,2]"},
		{name: "remove_any", run: func() (runtime.Value, error) { return arrays.RemoveAny(t.Context(), input, other) }, want: "[3,1,3]"},
		{name: "union", run: func() (runtime.Value, error) { return arrays.Union(t.Context(), input, other) }, want: "[3,1,2,4]"},
		{name: "intersection", run: func() (runtime.Value, error) { return arrays.Intersection(t.Context(), input, other) }, want: "[2]"},
		{name: "difference", run: func() (runtime.Value, error) { return arrays.Difference(t.Context(), input, other) }, want: "[3,1]"},
		{name: "symmetric_difference", run: func() (runtime.Value, error) { return arrays.SymmetricDifference(t.Context(), input, other) }, want: "[3,1,4]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.run()
			if err != nil || result.String() != test.want {
				t.Fatalf("result = %v, %v, want %s", result, err, test.want)
			}

			if input.String() != "[3,1,3,2]" || other.String() != "[2,4,4]" {
				t.Fatalf("input mutated: %s, %s", input, other)
			}
		})
	}
}

func TestSymmetricDifferenceParityAndRepresentatives(t *testing.T) {
	for count := 2; count <= 4; count++ {
		inputs := []runtime.Value{runtime.NewArrayWith(runtime.Int(7), runtime.Int(3), runtime.Int(7))}
		for index := 1; index < count; index++ {
			inputs = append(inputs, runtime.NewArrayWith(runtime.Float(7), runtime.Float(7)))
		}

		result, err := arrays.SymmetricDifference(t.Context(), inputs...)
		if err != nil {
			t.Fatal(err)
		}

		want := "[3]"
		if count%2 == 1 {
			want = "[7,3]"
		}

		if result.String() != want {
			t.Fatalf("%d-source XOR = %s, want %s", count, result, want)
		}

		if count%2 == 1 {
			first, err := result.(runtime.List).At(t.Context(), 0)
			if err != nil || first != runtime.Int(7) {
				t.Fatalf("lost first representative: %T, %v", first, err)
			}
		}
	}

	first := distinctCollisionValue{label: "first"}
	second := distinctCollisionValue{label: "second"}
	result, err := arrays.SymmetricDifference(t.Context(), runtime.NewArrayWith(first, first, second), runtime.NewArrayWith(first), runtime.NewArrayWith(first))
	if err != nil {
		t.Fatal(err)
	}

	assertDistinctCollisionValues(t, t.Context(), result.(runtime.List), first, second)
}

func TestCanonicalArrayEdgeCases(t *testing.T) {
	for _, list := range []runtime.List{runtime.EmptyArray(), runtime.NewArrayWith(runtime.Int(1))} {
		before := list.String()
		for _, index := range []runtime.Int{math.MinInt64, -1, math.MaxInt64} {
			value, err := arrays.At(t.Context(), list, index)
			if err != nil || value != runtime.None {
				t.Fatalf("at(%s,%d): %v, %v", list, index, value, err)
			}

			result, err := arrays.RemoveAt(t.Context(), list, index)
			if err != nil || result.String() != before || result == list {
				t.Fatalf("remove_at(%s,%d): %v, %v", list, index, result, err)
			}
		}
	}

	input := runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Int(3))
	result, err := arrays.Slice(t.Context(), input, runtime.Int(1), runtime.Int(math.MaxInt64))
	if err != nil || result.String() != "[2,3]" {
		t.Fatalf("overflow slice: %v, %v", result, err)
	}

	probe := &sliceProbe{List: runtime.EmptyArray()}
	result, err = callLegacy("pop", t.Context(), probe)
	if err != nil || result.String() != "[]" || probe.called {
		t.Fatalf("empty host pop: %v, %v, sliced=%t", result, err, probe.called)
	}

	for _, fn := range []runtime.Function{arrays.Concat, arrays.Union, arrays.Intersection, arrays.Difference, arrays.SymmetricDifference} {
		if _, err := fn(t.Context(), input); err == nil {
			t.Fatal("set/concat accepted fewer than two inputs")
		}

		if _, err := fn(t.Context(), runtime.EmptyArray(), runtime.EmptyArray(), runtime.Int(1)); err == nil {
			t.Fatal("empty inputs bypassed type validation")
		}
	}
}

func TestCanonicalArrayErrors(t *testing.T) {
	sentinel := errors.New("host list failed")
	host := &failingList{List: runtime.NewArrayWith(runtime.Int(1)), err: sentinel}
	for _, fn := range []func() (runtime.Value, error){
		func() (runtime.Value, error) { return arrays.Contains(t.Context(), host, runtime.Int(1)) },
		func() (runtime.Value, error) { return arrays.IndexOf(t.Context(), host, runtime.Int(1)) },
		func() (runtime.Value, error) { return arrays.Append(t.Context(), host, runtime.Int(1)) },
		func() (runtime.Value, error) { return arrays.Remove(t.Context(), host, runtime.Int(1)) },
		func() (runtime.Value, error) { return arrays.RemoveAny(t.Context(), runtime.EmptyArray(), host) },
		func() (runtime.Value, error) { return arrays.Union(t.Context(), host, runtime.EmptyArray()) },
		func() (runtime.Value, error) { return arrays.Intersection(t.Context(), host, runtime.EmptyArray()) },
		func() (runtime.Value, error) { return arrays.Difference(t.Context(), host, runtime.EmptyArray()) },
		func() (runtime.Value, error) {
			return arrays.SymmetricDifference(t.Context(), host, runtime.EmptyArray())
		},
	} {
		if _, err := fn(); !errors.Is(err, sentinel) {
			t.Fatalf("lost host error: %v", err)
		}
	}

	bad := runtime.NewArrayWith(failingEqualityValue{err: sentinel})
	for _, fn := range []runtime.Function{arrays.Union, arrays.Intersection, arrays.Difference, arrays.SymmetricDifference} {
		if _, err := fn(t.Context(), bad, bad); !errors.Is(err, sentinel) {
			t.Fatalf("lost equality error: %v", err)
		}

		canceled, cancel := context.WithCancel(t.Context())
		cancel()
		if _, err := fn(canceled, bad, bad); !errors.Is(err, context.Canceled) {
			t.Fatalf("lost cancellation: %v", err)
		}
	}
}

func TestArrayHostLengthPolicy(t *testing.T) {
	sentinel := errors.New("length unavailable")
	list := &lengthFailingList{List: runtime.NewArrayWith(runtime.Int(1), runtime.Int(2)), err: sentinel}
	other := runtime.NewArrayWith(runtime.Int(2))
	for _, name := range []string{"intersection", "outersection"} {
		if _, err := callLegacy(name, t.Context(), list, other); err != nil {
			t.Fatalf("%s must not request host length for capacity: %v", name, err)
		}
	}

	result, err := callLegacy("remove_nth", t.Context(), list, runtime.Int(0))
	if err != nil || result.String() != "[2]" || list.String() != "[1,2]" {
		t.Fatalf("legacy indexed removal: %v, %v; source=%s", result, err, list)
	}

	if _, err := arrays.RemoveAt(t.Context(), list, runtime.Int(0)); !errors.Is(err, sentinel) {
		t.Fatalf("canonical bounds validation must preserve length errors: %v", err)
	}
}
