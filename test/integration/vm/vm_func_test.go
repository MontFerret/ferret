package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestFunctionCall(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN TYPENAME(1)", "int"),
		Case("RETURN TYPENAME(1.1)", "float"),
		Case("WAIT(10) RETURN 1", 1),
		Case("RETURN LENGTH([1,2,3])", 3),
		Case("RETURN CONCAT('a', 'b', 'c')", "abc"),
		Case("RETURN CONCAT(CONCAT('a', 'b'), 'c', CONCAT('d', 'e'))", "abcde", "Nested calls"),
		SkipCaseArray(`
		LET arr = []
		LET a = 1
		LET res = APPEND(arr, a)
		RETURN res
		`,
			[]any{1}, "Append to array"),
		Case("LET duration = 10 WAIT(duration) RETURN 1", 1),
		CaseNil("RETURN (FALSE OR T::FAIL())?"),
		CaseNil("RETURN T::FAIL()?"),
		CaseArray(`FOR i IN [1, 2, 3, 4]
				LET duration = 10
		
				WAIT(duration)
		
				RETURN i * 2`,
			[]any{2, 4, 6, 8}),

		Case(`RETURN FIRST((FOR i IN 1..10 RETURN i * 2))`, 2),
		CaseArray(`RETURN UNION((FOR i IN 0..5 RETURN i), (FOR i IN 6..10 RETURN i))`, []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}),
	})
}

func TestFunctionCall0(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN TEST0()", "test0", "Should call a function with no arguments"),
		Case("RETURN TEST()", "test", "Should call a function with no arguments using fallback"),
	}, vm.WithFunctions(runtime.NewFunctionsBuilder().
		Set("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.String("test"), nil
		}).
		Set0("TEST0", func(ctx context.Context) (runtime.Value, error) {
			return runtime.String("test0"), nil
		}).
		Build()))
}

func TestBuiltinFunctions(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN LENGTH([1,2,3])", 3),
		Case("RETURN TYPENAME([1,2,3])", "list"),
		Case("RETURN TYPENAME({ a: 1, b: 2 })", "map"),
		Case("WAIT(10) RETURN 1", 1),
	})
}
