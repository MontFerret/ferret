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

func TestMutableArrayRegistrations(t *testing.T) {
	functions := arrayFunctions(t)
	for _, test := range []struct {
		names []string
		want  []string
	}{
		{functions.A1().Names(), []string{"clear", "pop", "shift", "sort"}},
		{functions.A2().Names(), []string{"push", "remove", "remove_at", "unshift"}},
		{functions.A3().Names(), []string{"insert", "set"}},
		{functions.Var().Names(), nil},
	} {
		var got []string
		for _, name := range test.names {
			if strings.HasPrefix(name, "arrays::mut::") {
				got = append(got, strings.TrimPrefix(name, "arrays::mut::"))
			} else if strings.HasPrefix(name, "mut::") {
				t.Fatalf("unexpected alias: %s", name)
			}
		}

		slices.Sort(got)
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("mutable registrations = %v, want %v", got, test.want)
		}
	}
}

func TestMutableArrayOperations(t *testing.T) {
	for _, test := range []struct {
		removed runtime.Value
		name    string
		want    string
		input   []runtime.Int
		args    []runtime.Value
	}{
		{name: "push", args: []runtime.Value{runtime.Int(1)}, want: "[1]"},
		{name: "push", input: []runtime.Int{1}, args: []runtime.Value{runtime.Int(1)}, want: "[1,1]"},
		{name: "unshift", args: []runtime.Value{runtime.Int(1)}, want: "[1]"},
		{name: "unshift", input: []runtime.Int{2, 3}, args: []runtime.Value{runtime.Int(1)}, want: "[1,2,3]"},
		{name: "pop", want: "[]", removed: runtime.None},
		{name: "pop", input: []runtime.Int{1}, want: "[]", removed: runtime.Int(1)},
		{name: "pop", input: []runtime.Int{1, 2}, want: "[1]", removed: runtime.Int(2)},
		{name: "shift", want: "[]", removed: runtime.None},
		{name: "shift", input: []runtime.Int{1}, want: "[]", removed: runtime.Int(1)},
		{name: "shift", input: []runtime.Int{1, 2}, want: "[2]", removed: runtime.Int(1)},
		{name: "set", input: []runtime.Int{1}, args: []runtime.Value{runtime.Int(0), runtime.Int(3)}, want: "[3]"},
		{name: "set", input: []runtime.Int{1, 2}, args: []runtime.Value{runtime.Int(1), runtime.Int(3)}, want: "[1,3]"},
		{name: "insert", args: []runtime.Value{runtime.Int(0), runtime.Int(1)}, want: "[1]"},
		{name: "insert", input: []runtime.Int{2, 3}, args: []runtime.Value{runtime.Int(0), runtime.Int(1)}, want: "[1,2,3]"},
		{name: "insert", input: []runtime.Int{1, 3}, args: []runtime.Value{runtime.Int(1), runtime.Int(2)}, want: "[1,2,3]"},
		{name: "insert", input: []runtime.Int{1, 2}, args: []runtime.Value{runtime.Int(2), runtime.Int(3)}, want: "[1,2,3]"},
		{name: "remove", args: []runtime.Value{runtime.Int(1)}, want: "[]"},
		{name: "remove", input: []runtime.Int{1}, args: []runtime.Value{runtime.Int(1)}, want: "[]"},
		{name: "remove", input: []runtime.Int{1, 2, 1, 3, 3, 1}, args: []runtime.Value{runtime.Float(1)}, want: "[2,3,3]"},
		{name: "remove", input: []runtime.Int{1, 2}, args: []runtime.Value{runtime.Int(9)}, want: "[1,2]"},
		{name: "remove", input: []runtime.Int{1, 1, 1}, args: []runtime.Value{runtime.Int(1)}, want: "[]"},
		{name: "remove_at", input: []runtime.Int{1}, args: []runtime.Value{runtime.Int(0)}, want: "[]", removed: runtime.Int(1)},
		{name: "remove_at", input: []runtime.Int{1, 2}, args: []runtime.Value{runtime.Int(0)}, want: "[2]", removed: runtime.Int(1)},
		{name: "remove_at", input: []runtime.Int{1, 2, 3}, args: []runtime.Value{runtime.Int(1)}, want: "[1,3]", removed: runtime.Int(2)},
		{name: "remove_at", input: []runtime.Int{1, 2}, args: []runtime.Value{runtime.Int(1)}, want: "[1]", removed: runtime.Int(2)},
		{name: "clear", want: "[]"},
		{name: "clear", input: []runtime.Int{1}, want: "[]"},
		{name: "clear", input: []runtime.Int{1, 2, 3}, want: "[]"},
		{name: "sort", want: "[]"},
		{name: "sort", input: []runtime.Int{1}, want: "[1]"},
		{name: "sort", input: []runtime.Int{3, 1, 2, 1}, want: "[1,1,2,3]"},
	} {
		for _, host := range []bool{false, true} {
			t.Run(test.name, func(t *testing.T) {
				native := runtime.NewArray(len(test.input) + 4)
				for _, value := range test.input {
					_ = native.Append(t.Context(), value)
				}

				var target runtime.List = native
				if host {
					custom := newMutableHostList(native)
					custom.expectedContext = t.Context()
					target = custom
				}

				args := append([]runtime.Value{target}, test.args...)
				result, err := callLegacy("arrays::mut::"+test.name, t.Context(), args...)
				if err != nil || native.String() != test.want {
					t.Fatalf("host=%v: array=%s, result=%v, err=%v, want=%s", host, native, result, err, test.want)
				}

				want := runtime.Value(target)
				if test.removed != nil {
					want = test.removed
				}

				if result != want {
					t.Fatalf("result=%v, want identity/value=%v", result, want)
				}

				if custom, ok := target.(*mutableHostList); ok {
					assertNoMutableCopies(t, custom)
				}
			})
		}
	}
}

func TestMutableArrayIndexValidation(t *testing.T) {
	for _, size := range []int{0, 1, 3} {
		for _, name := range []string{"set", "insert", "remove_at"} {
			indexes := []runtime.Value{runtime.Int(math.MinInt64), runtime.Int(-1), runtime.Int(size + 1), runtime.Int(math.MaxInt64), runtime.Float(0), runtime.True, runtime.None, nil}
			if name != "insert" {
				indexes = append(indexes, runtime.Int(size))
			}

			for _, index := range indexes {
				target := newMutableHostList(runtime.NewSizedArray(size))
				args := []runtime.Value{target, index}
				if name != "remove_at" {
					args = append(args, runtime.Int(9))
				}

				got, err := callLegacy("arrays::mut::"+name, t.Context(), args...)
				cause := runtime.ErrInvalidType
				if _, ok := index.(runtime.Int); ok {
					cause = runtime.ErrInvalidOperation
				}

				pos, attributed, _ := runtime.InvalidArgumentDetails(err)
				if got != runtime.None || !errors.Is(err, cause) || !attributed || pos != 1 {
					t.Fatalf("%s size=%d index=%v: result=%v err=%v", name, size, index, got, err)
				}

				if target.calls["SetAt"]+target.calls["Insert"]+target.calls["RemoveAt"] != 0 {
					t.Fatalf("invalid index reached host mutation: %v", target.calls)
				}
			}
		}
	}
}

func TestMutableArrayTargetAndCancellationValidation(t *testing.T) {
	var typedNil *mutableHostList
	for _, name := range mutableOperationNames() {
		for _, target := range []runtime.Value{nil, typedNil, (*runtime.Array)(nil), runtime.None, runtime.Int(1), runtime.NewObject(), runtime.NewRange(1, 3), &readableArray{Value: runtime.None}} {
			got, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, target)...)
			pos, attributed, _ := runtime.InvalidArgumentDetails(err)
			if got != runtime.None || !errors.Is(err, runtime.ErrInvalidType) || !attributed || pos != 0 {
				t.Fatalf("%s target=%T: result=%v err=%v", name, target, got, err)
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		target := newMutableHostList(runtime.NewArrayWith(runtime.Int(0), runtime.Int(1)))
		got, err := callLegacy("arrays::mut::"+name, ctx, mutableOperationArgs(name, target)...)
		if got != runtime.None || !errors.Is(err, context.Canceled) || len(target.calls) != 0 {
			t.Fatalf("%s canceled: result=%v calls=%v err=%v", name, got, target.calls, err)
		}
	}
}

func TestMutableArrayHostErrors(t *testing.T) {
	for _, test := range []struct{ name, method string }{
		{"push", "Append"}, {"unshift", "Insert"}, {"set", "SetAt"}, {"insert", "Insert"},
		{"pop", "RemoveAt"}, {"shift", "RemoveAt"}, {"remove_at", "RemoveAt"},
		{"clear", "Clear"}, {"sort", "SortAsc"}, {"sort", "Swap"},
		{"remove", "At"}, {"remove", "SetAt"}, {"remove", "RemoveAt"},
	} {
		for _, failure := range []error{errors.New("host failure"), runtime.ErrNotSupported, context.Canceled} {
			target := newMutableHostList(runtime.NewArrayWith(runtime.Int(0), runtime.Int(2), runtime.Int(1)))
			target.failure = failure
			target.failAt[test.method] = 1
			got, err := callLegacy("arrays::mut::"+test.name, t.Context(), mutableOperationArgs(test.name, target)...)
			if got != runtime.None || !errors.Is(err, failure) {
				t.Fatalf("%s/%s: result=%v err=%v", test.name, test.method, got, err)
			}

			assertNoMutableCopies(t, target)
		}
	}

	for _, name := range []string{"unshift", "pop", "shift", "set", "insert", "remove_at", "remove", "sort"} {
		for _, negative := range []bool{false, true} {
			target := newMutableHostList(runtime.NewArrayWith(runtime.Int(1)))
			wantErr := target.failure
			if negative {
				length := runtime.Int(-1)
				target.length = &length
				wantErr = runtime.ErrInvalidOperation
			} else {
				target.failAt["Length"] = 1
			}

			got, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, target)...)
			if got != runtime.None || !errors.Is(err, wantErr) || target.String() != "[1]" {
				t.Fatalf("%s negative=%v: result=%v target=%s err=%v", name, negative, got, target, err)
			}
		}
	}
}

func TestMutableArrayNoOpDoesNotProbeWrites(t *testing.T) {
	for _, name := range []string{"pop", "shift", "remove"} {
		target := newMutableHostList(runtime.EmptyArray())
		target.failAt["RemoveAt"], target.failAt["SetAt"] = 1, 1
		result, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, target)...)
		if err != nil || target.calls["RemoveAt"]+target.calls["SetAt"] != 0 {
			t.Fatalf("%s: result=%v calls=%v err=%v", name, result, target.calls, err)
		}
	}

	target := newMutableHostList(runtime.NewArrayWith(runtime.Int(1)))
	target.failAt["RemoveAt"], target.failAt["SetAt"] = 1, 1
	result, err := arrays.RemoveMutable(t.Context(), target, runtime.Int(0))
	if err != nil || result != target || target.calls["RemoveAt"]+target.calls["SetAt"] != 0 {
		t.Fatalf("no matches: result=%v calls=%v err=%v", result, target.calls, err)
	}
}

func TestMutableArrayPartialRemovalFailure(t *testing.T) {
	for _, method := range []string{"SetAt", "RemoveAt"} {
		target := newMutableHostList(runtime.NewArrayWith(runtime.Int(0), runtime.Int(1), runtime.Int(0), runtime.Int(2)))
		target.failAt[method] = 2
		got, err := arrays.RemoveMutable(t.Context(), target, runtime.Int(0))
		want := "[1,1,0,2]"
		if method == "RemoveAt" {
			want = "[1,2,0]"
		}

		if got != runtime.None || !errors.Is(err, target.failure) || target.String() != want {
			t.Fatalf("%s: result=%v target=%s err=%v", method, got, target, err)
		}
	}

	sentinel := errors.New("comparison failed")
	target := runtime.NewArrayWith(failingEqualityValue{err: sentinel})
	got, err := arrays.RemoveMutable(t.Context(), target, failingEqualityValue{err: sentinel})
	if got != runtime.None || !errors.Is(err, sentinel) {
		t.Fatalf("comparison: result=%v err=%v", got, err)
	}
}

func TestMutableArrayRemovalCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	target := newMutableHostList(runtime.NewArrayWith(runtime.Int(0), runtime.Int(1), runtime.Int(0)))
	target.after = func(name string) {
		if name == "SetAt" {
			cancel()
		}
	}
	got, err := arrays.RemoveMutable(ctx, target, runtime.Int(0))
	if got != runtime.None || !errors.Is(err, context.Canceled) || target.calls["RemoveAt"] != 0 {
		t.Fatalf("result=%v calls=%v err=%v", got, target.calls, err)
	}
}

func TestMutableArraySortParityAndStability(t *testing.T) {
	first := runtime.NewArrayWith(runtime.Int(1))
	second := runtime.NewArrayWith(runtime.Float(1))
	for _, host := range []bool{false, true} {
		input := runtime.NewArrayWith(second, runtime.NewArrayWith(runtime.Int(0)), first)
		expected, err := arrays.Sorted(t.Context(), input)
		if err != nil {
			t.Fatal(err)
		}

		var target runtime.List = input
		if host {
			target = newMutableHostList(input)
		}

		result, err := arrays.SortMutable(t.Context(), target)
		equal, compareErr := runtime.EqualValues(t.Context(), result, expected)
		if err != nil || compareErr != nil || !equal || result != target {
			t.Fatalf("sort: result=%v err=%v comparison=%v", result, err, compareErr)
		}

		middle, _ := target.At(t.Context(), 1)
		last, _ := target.At(t.Context(), 2)
		if middle != second || last != first {
			t.Fatal("sort lost stable ordering or nested identities")
		}
	}
}

func TestMutableArrayMalformedRemovalAndOrderingErrors(t *testing.T) {
	for _, name := range []string{"pop", "shift", "remove_at"} {
		target := newMutableHostList(runtime.NewArrayWith(runtime.Int(1)))
		target.nilRemoval = true
		result, err := callLegacy("arrays::mut::"+name, t.Context(), mutableOperationArgs(name, target)...)
		if result != runtime.None || !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("%s: result=%v err=%v", name, result, err)
		}
	}

	sentinel := errors.New("host ordering failed")
	value := failingOrderingValue{Value: runtime.None, failure: sentinel}
	target := runtime.NewArrayWith(value, value)
	result, err := arrays.SortMutable(t.Context(), target)
	if result != runtime.None || !errors.Is(err, sentinel) {
		t.Fatalf("result=%v err=%v", result, err)
	}
}

func mutableOperationNames() []string {
	return []string{"push", "pop", "shift", "unshift", "set", "insert", "remove", "remove_at", "clear", "sort"}
}

func mutableOperationArgs(name string, target runtime.Value) []runtime.Value {
	args := []runtime.Value{target}
	switch name {
	case "push", "unshift", "remove", "remove_at":
		args = append(args, runtime.Int(0))
	case "set", "insert":
		args = append(args, runtime.Int(0), runtime.Int(9))
	}

	return args
}

func assertNoMutableCopies(t *testing.T, target *mutableHostList) {
	t.Helper()
	if target.calls["Copy"]+target.calls["Clone"]+target.calls["Empty"] != 0 {
		t.Fatalf("mutable operation copied its target: %v", target.calls)
	}
}
