package runtime_test

import (
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestArrayBounds(t *testing.T) {
	ctx := t.Context()
	for _, values := range [][]runtime.Value{nil, {runtime.Int(7)}, {runtime.Int(7), runtime.Int(8)}} {
		for _, index := range []runtime.Int{math.MinInt64, -1, runtime.Int(len(values)), math.MaxInt64} {
			array := runtime.NewArrayOf(values).Copy().(*runtime.Array)
			before := array.String()
			value, err := array.At(ctx, index)
			if err != nil || value != runtime.None {
				t.Fatalf("At(%d) = %v, %v", index, value, err)
			}

			value, found, err := array.LookupAt(ctx, index)
			if err != nil || found || value != runtime.None {
				t.Fatalf("LookupAt(%d) = %v, %v, %v", index, value, found, err)
			}

			value, err = array.RemoveAt(ctx, index)
			if err != nil || value != runtime.None || array.String() != before {
				t.Fatalf("RemoveAt(%d) = %v, %v; array=%s", index, value, err, array)
			}

			for _, err := range []error{array.SetAt(ctx, index, runtime.Int(9)), array.Swap(ctx, index, 0), array.Swap(ctx, 0, index)} {
				if !errors.Is(err, runtime.ErrInvalidOperation) || array.String() != before {
					t.Fatalf("invalid write at %d: %v; array=%s", index, err, array)
				}
			}

			if index != runtime.Int(len(values)) {
				err := array.Insert(ctx, index, runtime.Int(9))
				if !errors.Is(err, runtime.ErrInvalidOperation) || array.String() != before {
					t.Fatalf("Insert(%d): %v; array=%s", index, err, array)
				}
			}
		}
	}

	array := runtime.EmptyArray()
	if err := array.Insert(ctx, 0, runtime.Int(7)); err != nil {
		t.Fatal(err)
	}

	if err := array.Insert(ctx, 1, runtime.Int(8)); err != nil || array.String() != "[7,8]" {
		t.Fatalf("insert at end: %s, %v", array, err)
	}

	if err := array.Swap(ctx, 0, 1); err != nil || array.String() != "[8,7]" {
		t.Fatalf("valid swap: %s, %v", array, err)
	}

	value, err := array.RemoveAt(ctx, 1)
	if err != nil || value != runtime.Int(7) {
		t.Fatalf("remove last: %v, %v", value, err)
	}

	value, err = array.RemoveAt(ctx, 0)
	if err != nil || value != runtime.Int(8) || array.String() != "[]" {
		t.Fatalf("remove only element: %v, %v; array=%s", value, err, array)
	}
}

func TestArraySliceBounds(t *testing.T) {
	for _, test := range []struct {
		want  string
		start runtime.Int
		end   runtime.Int
	}{
		{start: -1, end: 1, want: "[]"}, {start: 0, end: -1, want: "[]"}, {start: 2, end: 1, want: "[]"}, {start: 3, end: 4, want: "[]"},
		{start: math.MinInt64, end: math.MaxInt64, want: "[]"}, {start: math.MaxInt64, end: math.MaxInt64, want: "[]"},
		{start: 0, end: math.MaxInt64, want: "[1,2,3]"}, {start: 1, end: math.MaxInt64, want: "[2,3]"}, {start: 1, end: 1, want: "[]"}, {start: 1, end: 2, want: "[2]"},
	} {
		array := runtime.NewArrayWith(runtime.Int(1), runtime.Int(2), runtime.Int(3))
		result, err := array.Slice(t.Context(), test.start, test.end)
		if err != nil || result.String() != test.want || array.String() != "[1,2,3]" {
			t.Fatalf("Slice(%d, %d) = %v, %v; source=%s", test.start, test.end, result, err, array)
		}

		empty, err := runtime.EmptyArray().Slice(t.Context(), test.start, test.end)
		if err != nil || empty.String() != "[]" {
			t.Fatalf("empty Slice(%d, %d) = %v, %v", test.start, test.end, empty, err)
		}
	}
}

func TestArraySliceIndependentStorage(t *testing.T) {
	for _, test := range []struct {
		mutate func(runtime.List) error
		name   string
	}{
		{name: "replace", mutate: func(list runtime.List) error { return list.SetAt(t.Context(), 0, runtime.Int(99)) }},
		{name: "append", mutate: func(list runtime.List) error { return list.Append(t.Context(), runtime.Int(99)) }},
		{name: "remove", mutate: func(list runtime.List) error {
			_, err := list.RemoveAt(t.Context(), 0)

			return err
		}},
		{name: "sort", mutate: func(list runtime.List) error { return list.SortAsc(t.Context()) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := runtime.NewArrayWith(runtime.Int(4), runtime.Int(3), runtime.Int(2), runtime.Int(1)).CopyWithGrowth(8)
			left, err := source.Slice(t.Context(), 1, 3)
			if err != nil {
				t.Fatal(err)
			}

			right, err := source.Slice(t.Context(), 1, 3)
			if err != nil {
				t.Fatal(err)
			}

			if err := test.mutate(left); err != nil {
				t.Fatal(err)
			}

			if source.String() != "[4,3,2,1]" || right.String() != "[3,2]" {
				t.Fatalf("slice mutation changed source or sibling: %s, %s", source, right)
			}

			before := left.String()
			if err := source.SetAt(t.Context(), 1, runtime.Int(88)); err != nil {
				t.Fatal(err)
			}

			if left.String() != before || right.String() != "[3,2]" {
				t.Fatalf("source mutation changed slices: %s, %s", left, right)
			}

			if err := test.mutate(source); err != nil {
				t.Fatal(err)
			}

			if left.String() != before || right.String() != "[3,2]" {
				t.Fatalf("source %s changed slices: %s, %s", test.name, left, right)
			}
		})
	}
}

func TestArraySliceIsShallow(t *testing.T) {
	nested := runtime.NewArrayWith(runtime.Int(1))
	source := runtime.NewArrayWith(nested)
	result, err := source.Slice(t.Context(), 0, 1)
	if err != nil {
		t.Fatal(err)
	}

	value, err := result.At(t.Context(), 0)
	if err != nil || value != nested {
		t.Fatalf("slice must retain nested value identity: %v, %v", value, err)
	}
}
