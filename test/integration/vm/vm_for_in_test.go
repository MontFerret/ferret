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

func TestForIn(t *testing.T) {
	queryable := mock.NewQueryable(runtime.NewArrayWith(
		runtime.NewObjectWith(map[string]runtime.Value{"id": runtime.NewString("a")}),
		runtime.NewObjectWith(map[string]runtime.Value{"id": runtime.NewString("b")}),
	))

	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunSpecs(t, []spec.Spec{
		Array(`
			FOR i IN 1..5
				RETURN i
`, []any{1, 2, 3, 4, 5}),
		Array(`
			FOR _ IN 1..5
				RETURN 1
`, []any{1, 1, 1, 1, 1}, "Should use ignored variable"),
		Array(
			`
			FOR i IN 1..5
				LET x = i * 2
				RETURN x
		`,
			[]any{2, 4, 6, 8, 10},
		),
		Array(
			`
			FOR val, counter IN 1..5
				LET x = val
				TEST_FN(counter)
				LET y = counter
				RETURN [x, y]
				`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
		),
		Array(
			`
			FOR i IN [] RETURN i
				`,
			[]any{},
		),
		Array(
			`
			FOR i IN [1, 2, 3] RETURN i
				`,
			[]any{1, 2, 3},
		),
		Array(`
			FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
		),
		Array(`
			FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
		),
		Fn(`
			FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			func(actual any) error {
				hashMap := make(map[string]bool)
				expectedArr := []any{"bar", "foo", "qaz"}
				actualArr := actual.([]any)

				for _, v := range expectedArr {
					hashMap[v.(string)] = false
				}

				for _, v := range actualArr {
					if _, ok := hashMap[v.(string)]; !ok {
						return fmt.Errorf("unexpected value: %s", v.(string))
					}

					hashMap[v.(string)] = true
				}

				return nil
			},
		),
		Fn(`
			FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			func(actual any) error {
				hashMap := make(map[string]bool)
				expectedArr := []any{"a", "b", "c"}
				actualArr := actual.([]any)

				for _, v := range expectedArr {
					hashMap[v.(string)] = false
				}

				for _, v := range actualArr {
					if _, ok := hashMap[v.(string)]; !ok {
						return fmt.Errorf("unexpected value: %s", v.(string))
					}

					hashMap[v.(string)] = true
				}

				return nil
			},
		),
		Array(`
			FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		Array(`
			FOR i IN { items: [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] }.items RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		Array(`
			FOR i IN GET_ITEMS()
				RETURN i
		`, []any{1, 2, 3}, "Should iterate over a function call source"),
		Array(`
			FOR order IN QUERY "/orders" IN @api
				RETURN order.id
		`, []any{"a", "b"}, "Should iterate over a query source with nested IN"),
		Array(`
			FOR order IN QUERY "/orders" IN @api WITH {
				query: {
					status: "open"
				}
			}
				RETURN order.id
		`, []any{"a", "b"}, "Should iterate over a query source with WITH"),
		Array(`
			FOR order IN (
				QUERY "/orders" IN @api WITH {
					query: {
						status: "open"
					}
				}
			)
				RETURN order.id
		`, []any{"a", "b"}, "Should iterate over a parenthesized query source"),
		Array(`
			FOR item IN WAITFOR VALUE [1, 2]
				RETURN item
		`, []any{1, 2}, "Should iterate over a WAITFOR VALUE source"),
		Array(`
			FOR item IN MATCH "a" ("a" => [1], _ => [2])
				RETURN item
		`, []any{1}, "Should iterate over a MATCH source"),
	}, vm.WithFunction("TEST_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, nil
	}), vm.WithFunction("GET_ITEMS", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)), nil
	}), vm.WithParam("api", queryable))
}

func TestForInNonIterableSourceErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(`
			FOR item IN 42
				RETURN item
		`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid type",
			Contains: []string{"expected Iterable", "got Int"},
		}),
		spec.NewSpec(`
			FOR item IN TRUE
				RETURN item
		`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid type",
			Contains: []string{"expected Iterable", "got Boolean"},
		}),
	})
}
