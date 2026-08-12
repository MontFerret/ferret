package vm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

type observableKeyReadable struct {
	err    error
	values map[string]runtime.Value
	runtime.Int
	reads int
}

func (v *observableKeyReadable) Get(_ context.Context, key runtime.Value) (runtime.Value, error) {
	v.reads++
	if v.err != nil {
		return runtime.None, v.err
	}

	name, ok := key.(runtime.String)
	if !ok {
		return runtime.None, runtime.ErrNotFound
	}

	value, ok := v.values[string(name)]
	if !ok {
		return runtime.None, nil
	}

	return value, nil
}

type observableIndexReadable struct {
	err    error
	values []runtime.Value
	runtime.Int
	reads int
}

func (v *observableIndexReadable) At(_ context.Context, index runtime.Int) (runtime.Value, error) {
	v.reads++
	if v.err != nil {
		return runtime.None, v.err
	}

	if index < 0 || int(index) >= len(v.values) {
		return runtime.None, runtime.ErrNotFound
	}

	return v.values[index], nil
}

func TestDestructuringBindings(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
LET { name, age: years, address: { city }, values: [first, _, third] } = {
    name: "Ada",
    age: 37,
    address: { city: "London" },
    values: [1, 2, 3],
    extra: true
}
RETURN [name, years, city, first, third]
`, []any{"Ada", 37, "London", 1, 3}),
		Array(`
LET { name, missing, nested: { value } } = { name: "Ada" }
LET [first, second, third] = [1]
RETURN [name, missing, value, first, second, third]
`, []any{"Ada", nil, nil, 1, nil, nil}),
		Array(`
LET { value, nested: { child } } = NONE
LET [first, { second }] = NONE
RETURN [value, child, first, second]
`, []any{nil, nil, nil, nil}),
		Array(`
LET [first] = [1, 2, 3]
RETURN [first]
`, []any{1}),
		S(`
LET {} = {}
LET [] = []
VAR _ = 42
LET { ignored: _, nested: [_, _] } = { ignored: 1, nested: [2, 3] }
RETURN 1
`, 1),
		Array(`
VAR { count, label: name } = { count: 1, label: "start" }
count += 2
name = "done"
RETURN [count, name]
`, []any{3, "done"}),
		Array(`
VAR calls = 0
FUNC source() {
    calls += 1
    RETURN { left: 1, right: 2 }
}
LET { left, right } = source()
RETURN [left, right, calls]
`, []any{1, 2, 1}),
		Array(`
LET { value, ignored: _, nested: [_, value2] } = {
    value: 1,
    ignored: 2,
    nested: [3, 4]
}
RETURN [value, value2]
`, []any{1, 4}),
		S(`
LET { value } = { value: 7 }
FUNC read() => value
RETURN read()
`, 7),
		Array(`
VAR { value } = { value: 1 }
FUNC increment() {
    value += 1
    RETURN value
}
RETURN [increment(), increment(), value]
`, []any{2, 3, 3}),
		Array(`
LET { value } = { value: 1 }
FUNC inner() {
    LET { value } = { value: 2 }
    RETURN value
}
RETURN [value, inner()]
`, []any{1, 2}),
	})
}

func TestForDestructuringBindings(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
FOR { name, score }, counter IN [
    { name: "a", score: 1 },
    { name: "b", score: 3 },
    { name: "c", score: 2 }
]
    FILTER score > 1
    RETURN [name, score, counter]
`, []any{[]any{"b", 3, 1}, []any{"c", 2, 2}}),
		Array(`
FOR [name, { score }] IN [["a", { score: 2 }], ["b", { score: 1 }]] {
    RETURN [name, score]
}
`, []any{[]any{"a", 2}, []any{"b", 1}}),
		Array(`
FOR { name, score } IN [
    { name: "a", score: 2 },
    { name: "b", score: 1 },
    { name: "c", score: 3 }
]
    SORT score
    RETURN name
`, []any{"b", "a", "c"}),
		Array(`
FOR { gender } IN [{ gender: "m" }, { gender: "f" }, { gender: "m" }]
    COLLECT group = gender
    SORT group
    RETURN group
`, []any{"f", "m"}),
		Array(`
FOR { value, nested: [child] } IN [NONE]
    RETURN [value, child]
`, []any{[]any{nil, nil}}),
		Array(`
VAR calls = 0
FUNC source() {
    calls += 1
    RETURN [{ value: 1 }, { value: 2 }]
}
LET values = (
    FOR { value } IN source()
        RETURN value
)
RETURN [values, calls]
`, []any{[]any{1, 2}, 1}),
		Array(`
RETURN (
    FOR { value } IN [1]
        RETURN value
) ON ERROR RETURN []
`, []any{}),
	})
}

func TestForDestructuringValidatesEveryItem(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ErrorStr(`FOR { value } IN [{ value: 1 }, 2] RETURN value`, "cannot destructure Int as Object"),
		ErrorStr(`FOR [value] IN [[1], false] RETURN value`, "cannot destructure Boolean as Array"),
	})
}

func TestDestructuringAcceptsCapabilityBasedValuesAndSkipsIgnoredLookups(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d/object", level), func(t *testing.T) {
			value := &observableKeyReadable{
				Int:    0,
				values: map[string]runtime.Value{"kept": runtime.Int(42)},
			}

			program, err := spec.Compile(`LET { kept, ignored: _ } = @value RETURN kept`, level)
			if err != nil {
				t.Fatal(err)
			}

			out, err := spec.Run(program, vm.WithParam("value", value))
			if err != nil {
				t.Fatal(err)
			}

			if got, want := string(out), "42"; got != want {
				t.Fatalf("result = %s, want %s", got, want)
			}

			if got, want := value.reads, 1; got != want {
				t.Fatalf("object reads = %d, want %d", got, want)
			}
		})

		t.Run(fmt.Sprintf("O%d/array", level), func(t *testing.T) {
			value := &observableIndexReadable{
				Int:    0,
				values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3)},
			}

			program, err := spec.Compile(`LET [first, _, third] = @value RETURN [first, third]`, level)
			if err != nil {
				t.Fatal(err)
			}

			out, err := spec.Run(program, vm.WithParam("value", value))
			if err != nil {
				t.Fatal(err)
			}

			if got, want := string(out), "[1,3]"; got != want {
				t.Fatalf("result = %s, want %s", got, want)
			}

			if got, want := value.reads, 2; got != want {
				t.Fatalf("array reads = %d, want %d", got, want)
			}
		})
	}
}

func TestDestructuringPropagatesGetterFailures(t *testing.T) {
	getterErr := errors.New("getter failed")
	value := &observableKeyReadable{Int: 0, err: getterErr}

	program, err := spec.Compile(`LET { value } = @input RETURN value`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = spec.Run(program, vm.WithParam("input", value))
	if !errors.Is(err, getterErr) {
		t.Fatalf("expected getter failure, got %v", err)
	}
}

func TestDestructuringWrongShapeReportsNestedPattern(t *testing.T) {
	query := `LET { nested: [value] } = { nested: 42 }
RETURN value`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		program, err := spec.Compile(query, level)
		if err != nil {
			t.Fatal(err)
		}

		_, err = spec.Run(program)
		if err == nil {
			t.Fatalf("O%d expected destructuring error", level)
		}

		var runtimeErr *vm.RuntimeError
		if !errors.As(err, &runtimeErr) {
			t.Fatalf("O%d error type = %T, want *vm.RuntimeError", level, err)
		}

		formatted := runtimeErr.Format()
		if !strings.Contains(formatted, "cannot destructure Int as Array") {
			t.Fatalf("O%d unexpected error:\n%s", level, formatted)
		}

		if !strings.Contains(formatted, "[value]") {
			t.Fatalf("O%d error does not point at nested pattern:\n%s", level, formatted)
		}
	}
}

func TestEmptyDestructuringPatternsStillValidateShape(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ErrorStr(`LET {} = 1 RETURN 1`, "cannot destructure Int as Object"),
		ErrorStr(`LET [] = true RETURN 1`, "cannot destructure Boolean as Array"),
		ErrorStr(`LET { nested: [_] } = { nested: 1 } RETURN 1`, "cannot destructure Int as Array"),
	})
}
