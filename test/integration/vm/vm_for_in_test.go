package vm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForIn(t *testing.T) {
	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunSpecs(t, []Spec{
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
	}, vm.WithFunction("TEST_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, nil
	}))
}
