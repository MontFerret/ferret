package vm_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	pkgdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestLiteralSpread(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`RETURN [0, ...[1, 2], 3, ...[], ...NONE]`, []any{0, 1, 2, 3}, "array spread preserves entry order"),
		Object(`RETURN { a: 0, ...{ a: 1, b: 2 }, a: 3, ...{ c: 4 }, ...NONE }`, map[string]any{
			"a": 3,
			"b": 2,
			"c": 4,
		}, "object spread preserves later-property precedence"),
		Array(`
LET source = [1]
LET copied = [...source]
copied[0] = 9
RETURN [source, copied]
`, []any{[]any{1}, []any{9}}, "array spread does not mutate its source container"),
		Array(`
LET nested = { value: 1 }
LET source = [nested]
LET copied = [...source]
copied[0].value = 2
RETURN [source[0].value, copied[0].value]
`, []any{2, 2}, "array spread is shallow"),
		Array(`
LET source = { value: 1 }
LET copied = { ...source }
copied.value = 9
RETURN [source.value, copied.value]
`, []any{1, 9}, "object spread does not mutate its source container"),
		Array(`
LET nested = { value: 1 }
LET source = { nested }
LET copied = { ...source }
copied.nested.value = 2
RETURN [source.nested.value, copied.nested.value]
`, []any{2, 2}, "object spread is shallow"),
		Array(`RETURN [...HOST_ARRAY()]`, []any{1, 2}, "host-returned Array can be spread"),
		Object(`RETURN { ...HOST_OBJECT() }`, map[string]any{"value": 1}, "host-returned Object can be spread"),
		Error(`RETURN [...(1..2)]`, "generic iterable values cannot be spread into arrays"),
		Error(`RETURN { ...HOST_MAP() }`, "Map values without ObjectLike cannot be spread into objects"),
		Array(`RETURN ([...FAIL()]) ON ERROR RETURN ["recovered"]`, []any{"recovered"}, "spread-source failures use normal recovery"),
	},
		vm.WithFunction("HOST_ARRAY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)), nil
		}),
		vm.WithFunction("HOST_OBJECT", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}), nil
		}),
		vm.WithFunction("HOST_MAP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return genericMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)})}, nil
		}),
		vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return runtime.None, errors.New("spread source failed")
		}),
	)
}

func TestLiteralSpreadRuntimeCollections(t *testing.T) {
	query := `
LET genericList = HOST_LIST()
LET snapshotList = HOST_SNAPSHOT_LIST()
LET genericObject = HOST_OBJECT_LIKE()
LET snapshotObject = HOST_SNAPSHOT_OBJECT()

RETURN {
  lists: [[...genericList], [...snapshotList]],
  objects: [{...genericObject}, {...snapshotObject}]
}
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			fallbackErr := errors.New("generic fallback must not be used")
			snapshotList := &spreadSnapshotList{
				List:        runtime.NewArrayWith(runtime.NewInt(99)),
				snapshot:    runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				fallbackErr: fallbackErr,
			}
			snapshotObject := &spreadSnapshotObject{
				Map: runtime.NewObjectWith(map[string]runtime.Value{"wrong": runtime.True}),
				snapshot: runtime.NewObjectWith(map[string]runtime.Value{
					"active": runtime.True,
					"value":  runtime.NewInt(1),
				}),
				fallbackErr: fallbackErr,
			}

			RunSpecsWith(t, "runtime collections", compiler.New(compiler.WithOptimizationLevel(level)), []spec.Spec{
				Object(query, map[string]any{
					"lists": []any{
						[]any{1, 2},
						[]any{1, 2},
					},
					"objects": []any{
						map[string]any{"active": true, "value": 1},
						map[string]any{"active": true, "value": 1},
					},
				}, "runtime collection fallbacks and snapshots are equivalent"),
			},
				vm.WithFunction("HOST_LIST", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return struct{ runtime.List }{
						List: runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
					}, nil
				}),
				vm.WithFunction("HOST_SNAPSHOT_LIST", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return snapshotList, nil
				}),
				vm.WithFunction("HOST_OBJECT_LIKE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return &spreadObjectLike{
						Map: runtime.NewObjectWith(map[string]runtime.Value{
							"active": runtime.True,
							"value":  runtime.NewInt(1),
						}),
					}, nil
				}),
				vm.WithFunction("HOST_SNAPSHOT_OBJECT", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return snapshotObject, nil
				}),
			)

			if snapshotList.snapshotCalls != 1 || snapshotList.fallbackCalls != 0 {
				t.Fatalf(
					"unexpected list access counts: snapshot=%d fallback=%d",
					snapshotList.snapshotCalls,
					snapshotList.fallbackCalls,
				)
			}

			if snapshotObject.snapshotCalls != 1 || snapshotObject.fallbackCalls != 0 {
				t.Fatalf(
					"unexpected object access counts: snapshot=%d fallback=%d",
					snapshotObject.snapshotCalls,
					snapshotObject.fallbackCalls,
				)
			}
		})
	}
}

func TestLiteralSpreadSnapshotsRemainShallow(t *testing.T) {
	query := `
LET listSource = HOST_SNAPSHOT_LIST()
LET listCopy = [...listSource]
listCopy[0].value = 2

LET objectSource = HOST_SNAPSHOT_OBJECT()
LET objectCopy = {...objectSource}
objectCopy.nested.value = 3

RETURN [listSource[0].value, listCopy[0].value, objectSource.nested.value, objectCopy.nested.value]
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			fallbackErr := errors.New("generic fallback must not be used")
			listValues := runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{
				"value": runtime.NewInt(1),
			}))
			objectValues := runtime.NewObjectWith(map[string]runtime.Value{
				"nested": runtime.NewObjectWith(map[string]runtime.Value{
					"value": runtime.NewInt(1),
				}),
			})
			listSource := &spreadSnapshotList{
				List:        listValues,
				snapshot:    listValues,
				fallbackErr: fallbackErr,
			}
			objectSource := &spreadSnapshotObject{
				Map:         objectValues,
				snapshot:    objectValues,
				fallbackErr: fallbackErr,
			}

			RunSpecsWith(t, "snapshot shallow copy", compiler.New(compiler.WithOptimizationLevel(level)), []spec.Spec{
				Array(query, []any{2, 2, 3, 3}),
			},
				vm.WithFunction("HOST_SNAPSHOT_LIST", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return listSource, nil
				}),
				vm.WithFunction("HOST_SNAPSHOT_OBJECT", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return objectSource, nil
				}),
			)

			if listSource.snapshotCalls != 1 || listSource.fallbackCalls != 0 {
				t.Fatalf(
					"unexpected list access counts: snapshot=%d fallback=%d",
					listSource.snapshotCalls,
					listSource.fallbackCalls,
				)
			}

			if objectSource.snapshotCalls != 1 || objectSource.fallbackCalls != 0 {
				t.Fatalf(
					"unexpected object access counts: snapshot=%d fallback=%d",
					objectSource.snapshotCalls,
					objectSource.fallbackCalls,
				)
			}
		})
	}
}

func TestLiteralSpreadSnapshotCapabilityDoesNotDefineCompatibility(t *testing.T) {
	listSnapshotter := &spreadSnapshotList{
		List:     runtime.NewArrayWith(runtime.NewInt(1)),
		snapshot: runtime.NewArrayWith(runtime.NewInt(1)),
	}
	mapSnapshotter := &spreadSnapshotObject{
		Map:      runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}),
		snapshot: runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}),
	}

	RunSpecs(t, []spec.Spec{
		Error(`RETURN [...HOST_LIST_SNAPSHOT_ONLY()]`, "ListSnapshotter without List remains invalid"),
		Error(`RETURN {...HOST_MAP_SNAPSHOT_ONLY()}`, "MapSnapshotter without ObjectLike remains invalid"),
	},
		vm.WithFunction("HOST_LIST_SNAPSHOT_ONLY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return struct {
				runtime.Value
				runtime.ListSnapshotter
			}{
				Value:           runtime.NewString("snapshot-only list"),
				ListSnapshotter: listSnapshotter,
			}, nil
		}),
		vm.WithFunction("HOST_MAP_SNAPSHOT_ONLY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return struct {
				runtime.Value
				runtime.MapSnapshotter
			}{
				Value:          runtime.NewString("snapshot-only map"),
				MapSnapshotter: mapSnapshotter,
			}, nil
		}),
	)

	if listSnapshotter.snapshotCalls != 0 {
		t.Fatalf("snapshot-only list was invoked %d times", listSnapshotter.snapshotCalls)
	}

	if mapSnapshotter.snapshotCalls != 0 {
		t.Fatalf("snapshot-only map was invoked %d times", mapSnapshotter.snapshotCalls)
	}
}

func TestLiteralSpreadSnapshotFailures(t *testing.T) {
	tests := []struct {
		newValue func(snapshotErr, fallbackErr error) (runtime.Value, func() (int, int))
		name     string
		query    string
		span     string
	}{
		{
			name:  "list",
			query: `RETURN [...SOURCE()]`,
			span:  `...SOURCE()`,
			newValue: func(snapshotErr, fallbackErr error) (runtime.Value, func() (int, int)) {
				value := &spreadSnapshotList{
					List:        runtime.NewArrayWith(runtime.NewInt(1)),
					snapshotErr: snapshotErr,
					fallbackErr: fallbackErr,
				}

				return value, func() (int, int) {
					return value.snapshotCalls, value.fallbackCalls
				}
			},
		},
		{
			name:  "object",
			query: `RETURN {...SOURCE()}`,
			span:  `...SOURCE()`,
			newValue: func(snapshotErr, fallbackErr error) (runtime.Value, func() (int, int)) {
				value := &spreadSnapshotObject{
					Map:         runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}),
					snapshotErr: snapshotErr,
					fallbackErr: fallbackErr,
				}

				return value, func() (int, int) {
					return value.snapshotCalls, value.fallbackCalls
				}
			},
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.name), func(t *testing.T) {
				snapshotErr := errors.New("snapshot failed")
				fallbackErr := errors.New("generic fallback must not be used")
				value, calls := test.newValue(snapshotErr, fallbackErr)

				runtimeErr := runLiteralSpreadError(t, level, test.query, value)
				if !errors.Is(runtimeErr, snapshotErr) {
					t.Fatalf("expected snapshot error cause, got %v", runtimeErr)
				}

				assertRuntimeErrorMainSpan(t, runtimeErr, test.query, test.span)

				snapshotCalls, fallbackCalls := calls()
				if snapshotCalls != 1 || fallbackCalls != 0 {
					t.Fatalf("unexpected access counts: snapshot=%d fallback=%d", snapshotCalls, fallbackCalls)
				}

				recoveryValue, recoveryCalls := test.newValue(snapshotErr, fallbackErr)
				recoveryQuery := fmt.Sprintf(`RETURN (%s) ON ERROR RETURN ["recovered"]`, test.query[len("RETURN "):])

				RunSpecsWith(t, "snapshot recovery", compiler.New(compiler.WithOptimizationLevel(level)), []spec.Spec{
					Array(recoveryQuery, []any{"recovered"}),
				}, vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return recoveryValue, nil
				}))

				snapshotCalls, fallbackCalls = recoveryCalls()
				if snapshotCalls != 1 || fallbackCalls != 0 {
					t.Fatalf(
						"unexpected recovery access counts: snapshot=%d fallback=%d",
						snapshotCalls,
						fallbackCalls,
					)
				}
			})
		}
	}
}

func TestLiteralSpreadNilSnapshots(t *testing.T) {
	tests := []struct {
		newValue func() runtime.Value
		name     string
		query    string
		note     string
	}{
		{
			name:  "list",
			query: `RETURN [...SOURCE()]`,
			note:  "ListSnapshotter returned nil Array",
			newValue: func() runtime.Value {
				return &spreadSnapshotList{List: runtime.EmptyArray()}
			},
		},
		{
			name:  "object",
			query: `RETURN {...SOURCE()}`,
			note:  "MapSnapshotter returned nil Object",
			newValue: func() runtime.Value {
				return &spreadSnapshotObject{Map: runtime.NewObject()}
			},
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.name), func(t *testing.T) {
				runtimeErr := runLiteralSpreadError(t, level, test.query, test.newValue())
				if !errors.Is(runtimeErr, runtime.ErrInvalidOperation) {
					t.Fatalf("expected invalid-operation cause, got %v", runtimeErr)
				}

				if runtimeErr.Note != test.note {
					t.Fatalf("unexpected note: got %q, want %q", runtimeErr.Note, test.note)
				}

				assertRuntimeErrorMainSpan(t, runtimeErr, test.query, `...SOURCE()`)
			})
		}
	}
}

func TestLiteralSpreadEvaluationOrder(t *testing.T) {
	query := `RETURN [
  TRACE(1, 10),
  ...TRACE(2, [20, 30]),
  TRACE(3, 40),
  { first: TRACE(4, 1), ...TRACE(5, { second: 2 }), third: TRACE(6, 3) }
]`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			var trace []int

			RunSpecsWith(t, "trace", compiler.New(compiler.WithOptimizationLevel(level)), []spec.Spec{
				Array(query, []any{10, 20, 30, 40, map[string]any{
					"first":  1,
					"second": 2,
					"third":  3,
				}}),
			}, vm.WithFunction("TRACE", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
				trace = append(trace, int(args[0].(runtime.Int)))

				return args[1], nil
			}))

			if want := []int{1, 2, 3, 4, 5, 6}; !slices.Equal(trace, want) {
				t.Fatalf("unexpected evaluation trace: got %v, want %v", trace, want)
			}
		})
	}
}

func TestLiteralSpreadTypeErrors(t *testing.T) {
	tests := []struct {
		query   string
		message string
		span    string
		hint    string
	}{
		{
			query:   `RETURN [..."value"]`,
			message: "cannot spread String into Array",
			span:    `..."value"`,
			hint:    "Spread a List value or none inside an array literal",
		},
		{
			query:   `RETURN {...1}`,
			message: "cannot spread Int into Object",
			span:    `...1`,
			hint:    "Spread a Map value or none inside an object literal",
		},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		for _, test := range tests {
			t.Run(fmt.Sprintf("O%d/%s", level, test.message), func(t *testing.T) {
				program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("spread.fql", test.query))
				if err != nil {
					t.Fatalf("compile failed: %v", err)
				}

				instance, err := vm.New(program)
				if err != nil {
					t.Fatalf("vm init failed: %v", err)
				}
				defer instance.Close()

				env, err := vm.NewEnvironment(nil)
				if err != nil {
					t.Fatalf("environment init failed: %v", err)
				}

				_, err = instance.Run(context.Background(), env)
				if err == nil {
					t.Fatal("expected runtime error")
				}

				var runtimeErr *vm.RuntimeError
				if !errors.As(err, &runtimeErr) {
					t.Fatalf("expected runtime error, got %T", err)
				}

				if got := runtimeErr.Message; got != test.message {
					t.Fatalf("unexpected message: got %q, want %q", got, test.message)
				}

				if got := runtimeErr.Kind; got != pkgdiagnostics.TypeError {
					t.Fatalf("unexpected kind: got %s, want %s", got, pkgdiagnostics.TypeError)
				}

				if got := runtimeErr.Hint; got != test.hint {
					t.Fatalf("unexpected hint: got %q, want %q", got, test.hint)
				}

				if got := runtimeErr.Note; got != test.message {
					t.Fatalf("unexpected note: got %q, want %q", got, test.message)
				}

				if !errors.Is(runtimeErr, runtime.ErrInvalidType) {
					t.Fatalf("expected invalid-type cause, got %v", runtimeErr)
				}

				assertSpreadErrorSpan(t, runtimeErr, test.query, test.span)
			})
		}
	}
}

func assertSpreadErrorSpan(t *testing.T, runtimeErr *vm.RuntimeError, query, want string) {
	t.Helper()

	for _, span := range runtimeErr.Spans {
		if !span.Main {
			continue
		}

		if got := query[span.Span.Start:span.Span.End]; got != want {
			t.Fatalf("unexpected main span: got %q, want %q", got, want)
		}

		if !slices.Contains([]string{
			"expected List or none in array literal",
			"expected Map or none in object literal",
		}, span.Label) {
			t.Fatalf("unexpected spread error label: %q", span.Label)
		}

		return
	}

	t.Fatal("expected a main spread error span")
}

func runLiteralSpreadError(
	t *testing.T,
	level compiler.OptimizationLevel,
	query string,
	value runtime.Value,
) *vm.RuntimeError {
	t.Helper()

	program, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.New("spread.fql", query))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	instance, err := vm.New(program)
	if err != nil {
		t.Fatalf("vm init failed: %v", err)
	}
	defer instance.Close()

	env, err := vm.NewEnvironment([]vm.EnvironmentOption{
		vm.WithFunction("SOURCE", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			return value, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment init failed: %v", err)
	}

	result, err := instance.Run(context.Background(), env)
	if result != nil {
		defer result.Close()
	}

	if err == nil {
		t.Fatal("expected runtime error")
	}

	var runtimeErr *vm.RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	return runtimeErr
}

func assertRuntimeErrorMainSpan(t *testing.T, runtimeErr *vm.RuntimeError, query, want string) {
	t.Helper()

	for _, span := range runtimeErr.Spans {
		if !span.Main {
			continue
		}

		if got := query[span.Span.Start:span.Span.End]; got != want {
			t.Fatalf("unexpected main span: got %q, want %q", got, want)
		}

		return
	}

	t.Fatal("expected a main runtime error span")
}
