package collections_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func TestRegisteredCompatibility(t *testing.T) {
	functions := collectionFunctions(t)
	for _, tc := range []struct {
		input runtime.Value
		want  runtime.Value
		name  string
	}{
		{runtime.NewArrayWith(runtime.Int(1), runtime.Float(1)), runtime.Int(2), "count"},
		{runtime.NewArrayWith(runtime.Int(1), runtime.Float(1)), runtime.Int(1), "count_distinct"},
		{runtime.String("123"), runtime.True, "includes"},
		{runtime.String("a狐犬"), runtime.String("犬狐a"), "reverse"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var canonicalError error
			for _, prefix := range []string{"collections::", ""} {
				call := registeredCollectionFunction(t, functions, prefix+tc.name)
				got, err := call(t.Context(), tc.input)
				if err != nil || got != tc.want {
					t.Fatalf("%s%s = %v, %v; want %v", prefix, tc.name, got, err, tc.want)
				}

				_, err = call(t.Context(), runtime.Int(1))
				if !errors.Is(err, runtime.ErrInvalidType) {
					t.Fatalf("%s%s lost type error: %v", prefix, tc.name, err)
				}

				if prefix != "" {
					canonicalError = err
				} else if err.Error() != canonicalError.Error() {
					t.Fatalf("alias error %v differs from canonical error %v", err, canonicalError)
				}
			}
		})
	}
}

func TestRegisteredCompatibilityPreservesHostFailures(t *testing.T) {
	functions := collectionFunctions(t)
	cleanup := errors.New("cleanup failed")
	for _, primary := range []error{errors.New("host failed"), context.Canceled} {
		for _, prefix := range []string{"collections::", ""} {
			for _, name := range []string{"count", "count_distinct", "includes", "reverse"} {
				t.Run(prefix+name+"/"+primary.Error(), func(t *testing.T) {
					call := registeredCollectionFunction(t, functions, prefix+name)
					source := &scanSource{expected: t.Context(), nextErr: primary, closeErr: cleanup}
					list := &configuredList{
						Array: runtime.NewArrayWith(runtime.Int(1)), expected: t.Context(),
						calls: make(map[string]int), failures: map[string]error{"Append": primary}, closeErr: cleanup,
					}
					var input runtime.Value = source
					if name == "reverse" {
						input = list
					}

					_, err := call(t.Context(), input)
					if !errors.Is(err, primary) || !errors.Is(err, cleanup) {
						t.Fatalf("lost host or cleanup error: %v", err)
					}

					if name == "reverse" {
						if list.created == nil || list.created.closes != 1 || list.closes != 0 {
							t.Fatal("incorrect destination cleanup")
						}
					} else if source.iterations != 1 || source.closes != 1 || source.sourceCloses != 0 {
						t.Fatal("incorrect iterator cleanup")
					}
				})
			}
		}
	}
}

func collectionFunctions(t *testing.T) *runtime.Functions {
	t.Helper()

	ns := runtime.NewLibrary()
	collections.RegisterLib(ns)
	functions, err := ns.Build()
	if err != nil {
		t.Fatal(err)
	}

	return functions
}

func registeredCollectionFunction(t *testing.T, functions *runtime.Functions, name string) runtime.Function1 {
	t.Helper()

	if call, ok := functions.A1().Get(name); ok {
		return call
	}

	call, ok := functions.A2().Get(name)
	if !ok {
		t.Fatalf("missing collection function %s", name)
	}

	return func(ctx context.Context, value runtime.Value) (runtime.Value, error) {
		return call(ctx, value, runtime.Int(123))
	}
}
