package vm_test

import (
	"context"
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestFor(t *testing.T) {
	// Should not allocate memory if NONE is a return statement
	//{
	//	`FOR i IN 0..100
	//		RETURN NONE`,
	//	[]any{},
	//	ShouldEqualJSON,
	//},
	RunUseCases(t, []UseCase{
		SkipCaseCompilationError(`
			FOR foo IN foo
				RETURN foo
		`, "Should not compile FOR foo IN foo"),
		CaseArray(`
FOR i IN 1..5
	RETURN i
`, []any{1, 2, 3, 4, 5}),
		CaseArray(
			`
		FOR i IN 1..5
			LET x = i * 2
			RETURN x
		`,
			[]any{2, 4, 6, 8, 10},
		),
		CaseArray(
			`
		FOR val, counter IN 1..5
			LET x = val
			TEST_FN(counter)
			LET y = counter
			RETURN [x, y]
				`,
			[]any{[]any{1, 0}, []any{2, 1}, []any{3, 2}, []any{4, 3}, []any{5, 4}},
		),
		CaseArray(
			`FOR i IN [] RETURN i
				`,
			[]any{},
		),
		CaseArray(
			`FOR i IN [1, 2, 3] RETURN i
				`,
			[]any{1, 2, 3},
		),
		CaseArray(
			`FOR i, k IN [1, 2, 3] RETURN k`,
			[]any{0, 1, 2},
		),
		CaseArray(
			`FOR i IN ['foo', 'bar', 'qaz'] RETURN i`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseItems(
			`FOR i IN {a: 'bar', b: 'foo', c: 'qaz'} RETURN i`,
			[]any{"bar", "foo", "qaz"},
		),
		CaseArray(
			`FOR i, k IN {a: 'foo', b: 'bar', c: 'qaz'} RETURN k`,
			[]any{"a", "b", "c"},
		),
		CaseArray(
			`FOR i IN [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
		CaseArray(
			`FOR i IN { items: [{name: 'foo'}, {name: 'bar'}, {name: 'qaz'}] }.items RETURN i.name`,
			[]any{"foo", "bar", "qaz"},
		),
	}, vm.WithFunction("TEST_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return nil, nil
	}))
}
