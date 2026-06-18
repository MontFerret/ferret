package vm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestQueryable(t *testing.T) {
	queryable := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("ok")))
	queryableWithoutDefault := mock.NewQueryableWithoutDefault(runtime.NewArrayWith(runtime.NewString("ok")))
	sectionA := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("one")))
	sectionB := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("two")))

	RunSpecs(t, []spec.Spec{
		Array("RETURN @doc[~ css`.items`]", []any{"ok"}, "Should apply query literal"),
		Array(`RETURN @doc[~ "shortcut"]`, []any{"ok"}, "Should apply raw-string query shorthand with default dialect"),
		S(`RETURN @doc[~? "shortcut-one"]`, "ok", "Should apply raw-string query-one shorthand with default dialect"),
		S("RETURN @doc[~ css`.items`][0]", "ok", "Should apply query literal and index tail"),
		Array("RETURN QUERY `.default` IN @doc", []any{"ok"}, "Should apply query expression with default dialect"),
		Array("RETURN QUERY `.items` IN @doc USING css", []any{"ok"}, "Should apply query expression"),
		Array("RETURN QUERY @q IN @doc USING css", []any{"ok"}, "Should apply query expression with param expression"),
		Array("LET q = \".dynamic-var\"\nRETURN QUERY q IN @doc USING css", []any{"ok"}, "Should apply query expression with variable expression"),
		Array("RETURN QUERY `.default-with` IN @doc WITH { value: 4 }", []any{"ok"}, "Should apply default query expression with params"),
		Array("RETURN QUERY `.with` IN @doc USING css WITH { value: 1 }", []any{"ok"}, "Should apply query expression with params"),
		Array("RETURN QUERY `.options` IN @doc USING css OPTIONS { timeout: 5000 }", []any{"ok"}, "Should apply query expression with options"),
		Array("RETURN QUERY `.both` IN @doc USING css WITH { value: 2 } OPTIONS { timeout: 6000 }", []any{"ok"}, "Should apply query expression with params and options"),
		S("RETURN QUERY ONE `.one-both` IN @doc USING css WITH { value: 3 } OPTIONS { timeout: 7000 }", "ok", "Should apply query-one expression with params and options"),
		Array("RETURN @doc[~ sql`SELECT * FROM products`({ c: \"laptops\" })]", []any{"ok"}, "Should apply query literal with params"),
		S("RETURN @doc[~? sql`SELECT * FROM featured`({ c: \"tablets\" })]", "ok", "Should apply query-one literal with params"),
		Array("RETURN QUERY `SELECT * FROM products` IN @doc USING sql WITH { c: \"phones\" }", []any{"ok"}, "Should apply query expression with options"),
		Array("RETURN @doc[~ text]", []any{"ok"}, "Should apply query literal with no string expression"),
		Array(`RETURN @sections[* RETURN (QUERY "a" IN . USING css)][**]`, []any{"one", "two"}, "Should apply query expression to implicit current source"),
		spec.NewSpec("RETURN @val[~ css`x`]").Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "invalid type",
		}),
		spec.NewSpec("RETURN QUERY `.x` IN @noDefault", "Should fail when default dialect is unsupported").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"query dialect is required for this value; use USING <dialect>"}},
		),
		Array("RETURN QUERY `.x` IN @noDefault USING css", []any{"ok"}, "Should preserve explicit dialect for values without default query behavior"),
	}, vm.WithParams(map[string]runtime.Value{
		"doc":       queryable,
		"noDefault": queryableWithoutDefault,
		"q":         runtime.NewString(".dynamic-param"),
		"sections":  runtime.NewArrayWith(sectionA, sectionB),
		"val":       runtime.NewInt(1),
	}))

	t.Run("Should receive correct queries", func(t *testing.T) {
		var hasDefault bool
		var hasDefaultWith bool
		var hasDefaultShortcut bool
		var hasDefaultShortcutOne bool
		var hasCSS bool
		var hasCSSParam bool
		var hasCSSVar bool
		var hasSQLParams bool
		var hasSQLParamsOne bool
		var hasSQLQueryExpr bool
		var hasText bool
		var hasNoClauses bool
		var hasWithOnly bool
		var hasOptionsOnly bool
		var hasBoth bool
		var hasOneBoth bool

		for _, q := range queryable.MockQueries() {
			switch q.Kind {
			case runtime.EmptyString:
				if q.Expression == runtime.NewString(".default") {
					hasDefault = q.Params == runtime.None && q.Options == runtime.None
				}
				if q.Expression == runtime.NewString(".default-with") {
					hasDefaultWith = queryMapValue(t, q.Params, "value") == runtime.NewInt(4) && q.Options == runtime.None
				}
				if q.Expression == runtime.NewString("shortcut") {
					hasDefaultShortcut = q.Params == runtime.None && q.Options == runtime.None
				}
				if q.Expression == runtime.NewString("shortcut-one") {
					hasDefaultShortcutOne = q.Params == runtime.None && q.Options == runtime.None
				}
			case runtime.NewString("css"):
				if q.Expression == runtime.NewString(".items") {
					hasCSS = true
					hasNoClauses = q.Params == runtime.None && q.Options == runtime.None
				}
				if q.Expression == runtime.NewString(".dynamic-param") {
					hasCSSParam = true
				}
				if q.Expression == runtime.NewString(".dynamic-var") {
					hasCSSVar = true
				}
				if q.Expression == runtime.NewString(".with") {
					hasWithOnly = queryMapValue(t, q.Params, "value") == runtime.NewInt(1) && q.Options == runtime.None
				}
				if q.Expression == runtime.NewString(".options") {
					hasOptionsOnly = q.Params == runtime.None && queryMapValue(t, q.Options, "timeout") == runtime.NewInt(5000)
				}
				if q.Expression == runtime.NewString(".both") {
					hasBoth = queryMapValue(t, q.Params, "value") == runtime.NewInt(2) &&
						queryMapValue(t, q.Options, "timeout") == runtime.NewInt(6000)
				}
				if q.Expression == runtime.NewString(".one-both") {
					hasOneBoth = queryMapValue(t, q.Params, "value") == runtime.NewInt(3) &&
						queryMapValue(t, q.Options, "timeout") == runtime.NewInt(7000)
				}
			case runtime.NewString("text"):
				if q.Expression == runtime.EmptyString {
					hasText = true
				}
			case runtime.NewString("sql"):
				if q.Expression == runtime.NewString("SELECT * FROM products") {
					params, err := runtime.ToMap(context.Background(), q.Params)
					if err != nil {
						t.Fatalf("sql query expression params decode failed: %v", err)
					}

					value, err := params.Get(context.Background(), runtime.NewString("c"))
					if err == nil && value == runtime.NewString("phones") {
						hasSQLQueryExpr = true
					}
				}
				if q.Expression == runtime.NewString("SELECT * FROM featured") {
					params, err := runtime.ToMap(context.Background(), q.Params)
					if err != nil {
						t.Fatalf("sql query-one params decode failed: %v", err)
					}

					value, err := params.Get(context.Background(), runtime.NewString("c"))
					if err == nil && value == runtime.NewString("tablets") && q.Options == runtime.None {
						hasSQLParamsOne = true
					}
				}

				params, err := runtime.ToMap(context.Background(), q.Params)
				if err != nil {
					t.Fatalf("sql params decode failed: %v", err)
				}

				value, err := params.Get(context.Background(), runtime.NewString("c"))
				if err == nil && value == runtime.NewString("laptops") && q.Options == runtime.None {
					hasSQLParams = true
				}
			}
		}

		if !hasDefault {
			t.Fatal("expected omitted USING to produce an empty query kind")
		}
		if !hasDefaultWith {
			t.Fatal("expected default query WITH to populate params without options")
		}
		if !hasDefaultShortcut {
			t.Fatal("expected raw-string query shorthand to produce an empty query kind")
		}
		if !hasDefaultShortcutOne {
			t.Fatal("expected raw-string query-one shorthand to produce an empty query kind")
		}
		if !hasCSS {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and expression %q", "css", ".items"))
		}
		if !hasCSSParam {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and expression %q", "css", ".dynamic-param"))
		}
		if !hasCSSVar {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and expression %q", "css", ".dynamic-var"))
		}
		if !hasText {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and empty expression", "text"))
		}
		if !hasNoClauses {
			t.Fatal("expected omitted WITH and OPTIONS clauses to produce NONE")
		}
		if !hasWithOnly {
			t.Fatal("expected WITH to populate params without options")
		}
		if !hasOptionsOnly {
			t.Fatal("expected OPTIONS to populate options without params")
		}
		if !hasBoth {
			t.Fatal("expected WITH and OPTIONS to remain distinct")
		}
		if !hasOneBoth {
			t.Fatal("expected QUERY ONE to receive distinct params and options")
		}
		if !hasSQLQueryExpr {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and params containing %q=%q", "sql", "c", "phones"))
		}
		if !hasSQLParams {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and params containing %q=%q", "sql", "c", "laptops"))
		}
		if !hasSQLParamsOne {
			t.Fatal(fmt.Sprintf("expected to receive a query-one shorthand with kind %q and params containing %q=%q", "sql", "c", "tablets"))
		}
	})
}

func TestQueryExpressionClausesEvaluatedOnce(t *testing.T) {
	queryable := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("ok")))
	paramsCalls := 0
	optionsCalls := 0
	var clauseOrder []string

	RunSpecs(t, []spec.Spec{
		Array(
			"RETURN QUERY `.once` IN @doc USING css WITH PARAMS_FN() OPTIONS OPTIONS_FN()",
			[]any{"ok"},
			"Should evaluate query params and options once",
		),
	},
		vm.WithParam("doc", queryable),
		vm.WithFunction("PARAMS_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			paramsCalls++
			clauseOrder = append(clauseOrder, "params")
			return runtime.NewObjectWith(map[string]runtime.Value{"value": runtime.NewInt(1)}), nil
		}),
		vm.WithFunction("OPTIONS_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			optionsCalls++
			clauseOrder = append(clauseOrder, "options")
			return runtime.NewObjectWith(map[string]runtime.Value{"timeout": runtime.NewInt(5000)}), nil
		}),
	)

	if paramsCalls != 2 || optionsCalls != 2 {
		t.Fatalf("expected params and options to be evaluated once per O0/O1 execution, got params=%d options=%d", paramsCalls, optionsCalls)
	}
	if got, want := strings.Join(clauseOrder, ","), "params,options,params,options"; got != want {
		t.Fatalf("unexpected clause evaluation order: got %q, want %q", got, want)
	}

	queries := queryable.MockQueries()
	if len(queries) != 2 {
		t.Fatalf("expected one captured query per O0/O1 execution, got %d", len(queries))
	}
	for _, query := range queries {
		if queryMapValue(t, query.Params, "value") != runtime.NewInt(1) {
			t.Fatal("expected captured params")
		}
		if queryMapValue(t, query.Options, "timeout") != runtime.NewInt(5000) {
			t.Fatal("expected captured options")
		}
	}
}

func queryMapValue(t *testing.T, value runtime.Value, key string) runtime.Value {
	t.Helper()

	out, err := runtime.ToMap(context.Background(), value)
	if err != nil {
		t.Fatalf("failed to decode query map: %v", err)
	}

	item, err := out.Get(context.Background(), runtime.NewString(key))
	if err != nil {
		t.Fatalf("failed to read query map key %q: %v", key, err)
	}

	return item
}

func TestComplexQueries(t *testing.T) {
	queryableDoc := mock.NewNode("doc")
	queryableDB := mock.NewDBQueryable()
	queryableJSON := mock.NewJSONQueryable()

	RunSpecs(t, []spec.Spec{
		Array(
			"RETURN @doc\n    [~ css`.product`]\n    [~ css`.title`]\n    [~ text]",
			[]any{"title"},
			"Should chain apply operators",
		),
		Array(
			"RETURN\n  @db[~ sql`\n    SELECT name, price\n    FROM products\n    WHERE category = $c\n  `({ c: \"laptops\" })]",
			[]any{
				map[string]any{"name": "Laptop Pro", "price": 200},
			},
			"Should pass params to query",
		),
		Array(
			"RETURN @doc\n    [~ css`.product`]\n    [~ css`.title`]\n    [~ text]\n    [* FILTER . != \"\" RETURN UPPER(.)]",
			[]any{"TITLE"},
			"Should apply array operators after query chain",
		),
		Array(
			"RETURN @db[~ sql`SELECT name, price FROM products`]\n    [* FILTER .price > 100 RETURN .name]",
			[]any{"Laptop Pro"},
			"Should filter results from query",
		),
		Array(
			"RETURN @db[~ sql`SELECT name, price FROM products`]\n    [* FILTER .price > 100 RETURN .name]",
			[]any{"Laptop Pro"},
			"Should filter results from query using shorthand",
		),
		Array(
			"RETURN @doc[~ css`.product`]\n    [* RETURN .[~ css`.title`][~ text]]",
			[]any{
				[]any{"title"},
			},
			"Should apply query shorthand inside implicit current",
		),
		Array(
			"RETURN @doc\n    [~ css`.product`]\n    [* FILTER FIRST(.[~ css`.price`][~ text]) != \"\"]\n    [~ css`.title`]\n    [~ text]",
			[]any{
				[]any{"title"},
			},
			"Should combine query apply inside array filter",
		),
		Array(
			"RETURN @json\n    [~ jp`$.orders[*]`]\n    [* FILTER .total > 100]\n    [* RETURN {\n         id: .id,\n         items: .items[* RETURN .name]\n       }]",
			[]any{
				map[string]any{
					"id":    1,
					"items": []any{"Item A", "Item B"},
				},
			},
			"Should project nested array operators",
		),
		Array(
			"RETURN @doc\n    [~ css`.product`]\n    [* RETURN {\n         title: FIRST(.[~ css`.title`][~ text]),\n         price: FIRST(.[~ css`.price`][~ text])\n       }]",
			[]any{
				map[string]any{"title": "title", "price": "price"},
			},
			"Should support nested apply inside projections",
		),
		Array("RETURN @doc[~ nil`foo`]", []any{}, "Should return empty array for queryable that returns empty list"),
		Array("RETURN @doc[~ nil`foo`]?.[*].name", []any{}, "Should return empty array for queryable that returns empty list"),
	}, vm.WithParams(map[string]runtime.Value{
		"doc":  queryableDoc,
		"db":   queryableDB,
		"json": queryableJSON,
	}))
}

func TestQueryableListInput(t *testing.T) {
	queryableDoc := mock.NewNode("doc")
	queryableNil := mock.NewNilListQueryable()
	queryableA := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("a1"), runtime.NewString("a2")))
	queryableB := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("b1")))

	RunSpecs(t, []spec.Spec{
		Array(
			"RETURN [@qA, @qB][~ text]",
			[]any{"a1", "a2", "b1"},
			"Should apply query to list inputs and flatten results",
		),
		Array(
			"RETURN [@qNil, @qA][~ text]",
			[]any{"a1", "a2"},
			"Should ignore nil list results when flattening",
		),
		Array(
			"RETURN [@doc, @doc][~ css`.title`][~ text]",
			[]any{"title", "title"},
			"Should chain queries over list inputs",
		),
		spec.NewSpec("RETURN [@qA, 1][~ text]", "Should fail when list element is not queryable").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid type"},
		),
	}, vm.WithParams(map[string]runtime.Value{
		"doc":  queryableDoc,
		"qNil": queryableNil,
		"qA":   queryableA,
		"qB":   queryableB,
	}))
}

func TestQueryableModifiers(t *testing.T) {
	queryableMany := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("a"), runtime.NewString("b")))
	queryableOne := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("only")))
	queryableEmpty := mock.NewQueryable(runtime.NewArray(0))
	queryableObject := mock.NewQueryable(runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{
		"innerText": runtime.NewString("Title"),
	})))

	RunSpecs(t, []spec.Spec{
		S("RETURN QUERY EXISTS `.items` IN @many USING css", true, "EXISTS should return true for non-empty result"),
		S("RETURN QUERY EXISTS `.items` IN @empty USING css", false, "EXISTS should return false for empty result"),
		S("RETURN QUERY COUNT `.items` IN @many USING css", 2, "COUNT should return result length"),
		S("RETURN QUERY COUNT `.items` IN @empty USING css", 0, "COUNT should return zero for empty result"),
		S("RETURN QUERY ONE `.items` IN @one USING css", "only", "ONE should return the only result"),
		S("RETURN QUERY ONE `.items` IN @many USING css", "a", "ONE should return the first result"),
		Nil("RETURN QUERY ONE `.items` IN @empty USING css", "ONE should return NONE for empty result"),
		S("RETURN @many[~? css`.items`]", "a", "query-one shorthand should return the first result"),
		Nil("RETURN @empty[~? css`.items`]", "query-one shorthand should return NONE for empty result"),
		S("RETURN [@empty, @many][~? css`.items`]", "a", "query-one shorthand should short-circuit list source on first match"),
		S("RETURN @node[~? css`h1`]?.innerText", "Title", "query-one shorthand should compose with optional access"),
		Nil("RETURN @empty[~? css`.missing`]?.innerText", "query-one shorthand should allow optional access on NONE"),
		Nil("LET obj = NONE RETURN obj?.prop", "safe member access should remain valid"),
		spec.NewSpec("LET maybe = (QUERY ONE `.items` IN @empty USING css)?\nRETURN maybe.foo", "Optional query should not swallow the next instruction").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"cannot read property", "\"foo\""}},
		),
		spec.NewSpec("LET maybe = QUERY ONE `.items` IN @empty USING css ON ERROR RETURN NONE\nRETURN maybe.foo", "Explicit query recovery should not swallow the next instruction").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"cannot read property", "\"foo\""}},
		),
	}, vm.WithParams(map[string]runtime.Value{
		"many":  queryableMany,
		"one":   queryableOne,
		"empty": queryableEmpty,
		"node":  queryableObject,
	}))
}
