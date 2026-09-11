package arrays_test

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func TestLegacyArrayRegistrations(t *testing.T) {
	functions := arrayFunctions(t)
	for _, test := range []struct {
		name string
		got  []string
		want []string
	}{
		{"unary", functions.A1().Names(), []string{"first", "flatten", "last", "pop", "shift", "sorted", "sorted_unique", "unique"}},
		{"binary", functions.A2().Names(), []string{"append", "flatten", "nth", "position", "push", "remove_nth", "remove_value", "remove_values", "slice", "unshift"}},
		{"ternary", functions.A3().Names(), []string{"append", "position", "push", "remove_value", "slice", "unshift"}},
		{"variadic", functions.Var().Names(), []string{"intersection", "minus", "outersection", "union", "union_distinct"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			global := slices.DeleteFunc(test.got, func(name string) bool { return strings.Contains(name, "::") })
			slices.Sort(global)
			if !reflect.DeepEqual(global, test.want) {
				t.Fatalf("registrations = %v, want %v", global, test.want)
			}
		})
	}
}

func TestLegacyOutersectionPresence(t *testing.T) {
	outersection, _ := arrayFunctions(t).Var().Get("outersection")
	// Membership is exactly one source, not one occurrence or odd source parity.
	for _, test := range []struct {
		name   string
		want   string
		inputs []runtime.Value
	}{
		{name: "duplicates in one source", inputs: []runtime.Value{runtime.NewArrayWith(runtime.Int(1), runtime.Int(1)), runtime.EmptyArray()}, want: "[1]"},
		{name: "two sources", inputs: []runtime.Value{runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1))}, want: "[]"},
		{name: "three sources", inputs: []runtime.Value{runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1))}, want: "[]"},
		{name: "four sources", inputs: []runtime.Value{runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Int(1))}, want: "[]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := outersection(context.Background(), test.inputs...)
			if err != nil || result.String() != test.want {
				t.Fatalf("outersection = %v, %v, want %s", result, err, test.want)
			}
		})
	}
}

func BenchmarkArraySets(b *testing.B) {
	functions := arrayFunctions(b)
	for _, size := range []int{16, 1024} {
		for _, name := range []string{"union_distinct", "intersection", "minus", "outersection"} {
			b.Run(name+"/"+runtime.Int(size).String(), func(b *testing.B) {
				left := runtime.NewArray(size)
				right := runtime.NewArray(size)
				for index := range size {
					_ = left.Append(b.Context(), runtime.Int(index))
					_ = right.Append(b.Context(), runtime.Int(index+size/2))
				}

				fn, _ := functions.Var().Get(name)
				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					if _, err := fn(b.Context(), left, right); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	for _, size := range []int{16, 1024} {
		b.Run(runtime.Int(size).String(), func(b *testing.B) {
			input := runtime.NewArray(size)
			for index := range size {
				_ = input.Append(b.Context(), runtime.Int(index/2))
			}

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if _, err := arrays.SymmetricDifference(b.Context(), input, input, input); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Legacy tests exercise the actual registered functions after Go API migration.
func callLegacy(name string, ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	library := runtime.NewLibrary()
	arrays.RegisterLib(library)
	functions, err := library.Build()
	if err != nil {
		return runtime.None, err
	}

	switch len(args) {
	case 1:
		if fn, ok := functions.A1().Get(name); ok {
			return fn(ctx, args[0])
		}
	case 2:
		if fn, ok := functions.A2().Get(name); ok {
			return fn(ctx, args[0], args[1])
		}
	case 3:
		if fn, ok := functions.A3().Get(name); ok {
			return fn(ctx, args[0], args[1], args[2])
		}
	}

	if fn, ok := functions.Var().Get(name); ok {
		return fn(ctx, args...)
	}

	return runtime.None, fmt.Errorf("no registered %s/%d", name, len(args))
}

func arrayFunctions(t testing.TB) *runtime.Functions {
	t.Helper()

	library := runtime.NewLibrary()
	arrays.RegisterLib(library)
	functions, err := library.Build()
	if err != nil {
		t.Fatal(err)
	}

	return functions
}
