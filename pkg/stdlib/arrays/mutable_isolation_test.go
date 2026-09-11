package arrays_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

func TestMutableArrayCopyIsolation(t *testing.T) {
	for _, name := range mutableOperationNames() {
		for _, copyName := range []string{"copy", "slice", "append", "remove", "remove_at", "sorted"} {
			t.Run(name+"/"+copyName, func(t *testing.T) {
				source := runtime.NewArray(16)
				for _, value := range []runtime.Value{runtime.Int(3), runtime.Int(0), runtime.Int(1)} {
					_ = source.Append(t.Context(), value)
				}

				sibling, err := source.Slice(t.Context(), 0, 3)
				if err != nil {
					t.Fatal(err)
				}

				var copied runtime.Value
				switch copyName {
				case "copy":
					copied, err = runtime.Copy[runtime.List](source)
				case "slice":
					copied, err = arrays.Slice(t.Context(), source, runtime.Int(0), runtime.Int(3))
				case "append":
					copied, err = arrays.Append(t.Context(), source, runtime.Int(4))
				case "remove":
					copied, err = arrays.Remove(t.Context(), source, runtime.Int(99))
				case "remove_at":
					copied, err = arrays.RemoveAt(t.Context(), source, runtime.Int(-1))
				case "sorted":
					copied, err = arrays.Sorted(t.Context(), source)
				}

				if err != nil {
					t.Fatal(err)
				}

				if _, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, copied)...); err != nil {
					t.Fatal(err)
				}

				if source.String() != "[3,0,1]" || sibling.String() != "[3,0,1]" {
					t.Fatalf("mutation leaked: source=%s sibling=%s", source, sibling)
				}

				copiedBefore := copied.String()
				if _, err := arrays.SetMutable(t.Context(), source, runtime.Int(0), runtime.Int(88)); err != nil {
					t.Fatal(err)
				}

				if copied.String() != copiedBefore || sibling.String() != "[3,0,1]" {
					t.Fatal("source mutation leaked into a copy")
				}
			})
		}
	}
}

func TestMutableArrayNestedValuesAndBorrowedResources(t *testing.T) {
	for _, operation := range []func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error){
		arrays.PushMutable, arrays.UnshiftMutable,
		func(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
			return arrays.SetMutable(ctx, array, runtime.Int(0), value)
		},
		func(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
			return arrays.InsertMutable(ctx, array, runtime.Int(0), value)
		},
	} {
		nested := runtime.NewArrayWith(runtime.Int(1))
		target := runtime.NewArrayWith(runtime.None)
		got, err := operation(t.Context(), target, nested)
		if err != nil || got != target {
			t.Fatalf("insert: result=%v err=%v", got, err)
		}

		_, _ = arrays.PushMutable(t.Context(), nested, runtime.Int(2))
		found, err := target.Contains(t.Context(), runtime.NewArrayWith(runtime.Int(1), runtime.Int(2)))
		if err != nil || !found {
			t.Fatalf("nested value was copied: target=%s err=%v", target, err)
		}
	}

	nested := runtime.NewArrayWith(runtime.Int(1))
	target := runtime.NewArrayWith(nested, runtime.NewArrayWith(runtime.Float(1)), runtime.Int(2))
	if _, err := arrays.RemoveMutable(t.Context(), target, nested); err != nil || target.String() != "[2]" {
		t.Fatalf("nested equality: target=%s err=%v", target, err)
	}

	for _, name := range []string{"pop", "shift", "remove_at", "clear"} {
		resource := &mutableBorrowedValue{Value: runtime.None}
		target := runtime.NewArrayWith(resource)
		got, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, target)...)
		if err != nil || resource.closes != 0 || resource.copies != 0 {
			t.Fatalf("%s: result=%v copies=%d closes=%d err=%v", name, got, resource.copies, resource.closes, err)
		}

		if name != "clear" && got != resource {
			t.Fatalf("%s changed the removed value's identity", name)
		}
	}
}
