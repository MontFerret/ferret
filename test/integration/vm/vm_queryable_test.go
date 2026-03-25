package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"github.com/MontFerret/ferret/v2/test/spec/mock"
)

func TestQueryable(t *testing.T) {
	queryable := mock.NewQueryable(runtime.NewArrayWith(runtime.NewString("ok")))

	RunSpecs(t, []spec.Spec{
		Array("RETURN @doc[~ css`.items`]", []any{"ok"}, "Should apply query literal"),
		S("RETURN @doc[~ css`.items`][0]", "ok", "Should apply query literal and index tail"),
		Array("RETURN QUERY `.items` IN @doc USING css", []any{"ok"}, "Should apply query expression"),
		Array("RETURN QUERY @q IN @doc USING css", []any{"ok"}, "Should apply query expression with param payload"),
		Array("LET q = \".dynamic-var\"\nRETURN QUERY q IN @doc USING css", []any{"ok"}, "Should apply query expression with variable payload"),
		Array("RETURN @doc[~ sql`SELECT * FROM products`({ c: \"laptops\" })]", []any{"ok"}, "Should apply query literal with params"),
		Array("RETURN QUERY `SELECT * FROM products` IN @doc USING sql WITH { c: \"phones\" }", []any{"ok"}, "Should apply query expression with options"),
		Array("RETURN @doc[~ text]", []any{"ok"}, "Should apply query literal with no string payload"),
		spec.New("RETURN @val[~ css`x`]").Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "Invalid type",
		}),
	}, vm.WithParams(map[string]runtime.Value{
		"doc": queryable,
		"q":   runtime.NewString(".dynamic-param"),
		"val": runtime.NewInt(1),
	}))

	t.Run("Should receive correct queries", func(t *testing.T) {
		var hasCSS bool
		var hasCSSParam bool
		var hasCSSVar bool
		var hasSQLParams bool
		var hasSQLQueryExpr bool
		var hasText bool

		for _, q := range queryable.MockQueries() {
			switch q.Kind {
			case runtime.NewString("css"):
				if q.Payload == runtime.NewString(".items") {
					hasCSS = true
				}
				if q.Payload == runtime.NewString(".dynamic-param") {
					hasCSSParam = true
				}
				if q.Payload == runtime.NewString(".dynamic-var") {
					hasCSSVar = true
				}
			case runtime.NewString("text"):
				if q.Payload == runtime.EmptyString {
					hasText = true
				}
			case runtime.NewString("sql"):
				if q.Payload == runtime.NewString("SELECT * FROM products") {
					params, err := runtime.ToMap(context.Background(), q.Options)
					if err != nil {
						t.Fatalf("sql query expression params decode failed: %v", err)
					}

					value, err := params.Get(context.Background(), runtime.NewString("c"))
					if err == nil && value == runtime.NewString("phones") {
						hasSQLQueryExpr = true
					}
				}

				params, err := runtime.ToMap(context.Background(), q.Options)
				if err != nil {
					t.Fatalf("sql params decode failed: %v", err)
				}

				value, err := params.Get(context.Background(), runtime.NewString("c"))
				if err == nil && value == runtime.NewString("laptops") {
					hasSQLParams = true
				}
			}
		}

		if !hasCSS {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and payload %q", "css", ".items"))
		}
		if !hasCSSParam {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and payload %q", "css", ".dynamic-param"))
		}
		if !hasCSSVar {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and payload %q", "css", ".dynamic-var"))
		}
		if !hasText {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and empty payload", "text"))
		}
		if !hasSQLQueryExpr {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and params containing %q=%q", "sql", "c", "phones"))
		}
		if !hasSQLParams {
			t.Fatal(fmt.Sprintf("expected to receive a query with kind %q and params containing %q=%q", "sql", "c", "laptops"))
		}
	})
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
		spec.New("RETURN [@qA, 1][~ text]", "Should fail when list element is not queryable").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "Invalid type"},
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

	RunSpecs(t, []spec.Spec{
		S("RETURN QUERY EXISTS `.items` IN @many USING css", true, "EXISTS should return true for non-empty result"),
		S("RETURN QUERY EXISTS `.items` IN @empty USING css", false, "EXISTS should return false for empty result"),
		S("RETURN QUERY COUNT `.items` IN @many USING css", 2, "COUNT should return result length"),
		S("RETURN QUERY COUNT `.items` IN @empty USING css", 0, "COUNT should return zero for empty result"),
		S("RETURN QUERY ANY `.items` IN @many USING css", "a", "ANY should return first result"),
		Nil("RETURN QUERY ANY `.items` IN @empty USING css", "ANY should return NONE for empty result"),
		S("RETURN QUERY VALUE `.items` IN @many USING css", "a", "VALUE should return first result"),
		spec.New("RETURN QUERY VALUE `.items` IN @empty USING css", "VALUE should fail for empty result").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"QUERY VALUE expected at least one match"}},
		),
		S("RETURN QUERY ONE `.items` IN @one USING css", "only", "ONE should return the only result"),
		spec.New("RETURN QUERY ONE `.items` IN @empty USING css", "ONE should fail for empty result").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"QUERY ONE expected exactly one match"}},
		),
		spec.New("RETURN QUERY ONE `.items` IN @many USING css", "ONE should fail for multiple results").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"QUERY ONE expected exactly one match"}},
		),
		Nil("LET maybe = (QUERY VALUE `.items` IN @empty USING css)?\nRETURN maybe", "VALUE assertion should be catchable with optional operator"),
		Nil("LET maybe = (QUERY ONE `.items` IN @empty USING css)?\nRETURN maybe", "ONE assertion should be catchable for empty result with optional operator"),
		Nil("LET maybe = (QUERY ONE `.items` IN @many USING css)?\nRETURN maybe", "ONE assertion should be catchable for multi result with optional operator"),
		spec.New("LET maybe = (QUERY ONE `.items` IN @empty USING css)?\nRETURN maybe.foo", "Catch should not swallow the first instruction after a guarded QUERY ONE").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Contains: []string{"Cannot read property", "\"foo\""}},
		),
	}, vm.WithParams(map[string]runtime.Value{
		"many":  queryableMany,
		"one":   queryableOne,
		"empty": queryableEmpty,
	}))
}
