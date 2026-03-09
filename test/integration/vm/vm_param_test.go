package vm_test

import (
	"testing"
)

func TestParam(t *testing.T) {
	RunUseCases(t,
		[]UseCase{
			CaseRuntimeErrorStr(`RETURN @foo`, "Missing parameter"),
			CaseRuntimeErrorStr(`
FUNC read() => @foo
RETURN read()
`, "Missing parameter"),
			Case(`RETURN @str`, "bar", "Should return a value of a parameter"),
			Case(`RETURN @int + @int`, 2, "Should return a sum of two parameters"),
			Case(`RETURN @obj.str1 + @obj.str2`, "foobar", "Should return a concatenated string of two parameter properties"),
			CaseArray(`FOR i IN @values1 RETURN i`, []any{1, 2, 3, 4}, "Should iterate over an array parameter"),
			CaseArray(`FOR i IN @values2 SORT i RETURN i`, []any{"a", "b", "c", "d"}, "Should iterate over an object parameter"),
			CaseArray(`FOR i IN @start..@end RETURN i`, []any{1, 2, 3, 4, 5}, "Should iterate over a range parameter"),
			Case(`RETURN @obj.str1`, "foo", "Should be possible to use in member expression"),
		},
		WithParam("str", "bar"),
		WithParam("int", 1),
		WithParam("bool", true),
		WithParam("obj", map[string]interface{}{"str1": "foo", "str2": "bar"}),
		WithParam("values1", []int{1, 2, 3, 4}),
		WithParam("values2", map[string]interface{}{"a": "a", "b": "b", "c": "c", "d": "d"}),
		WithParam("start", 1),
		WithParam("end", 5),
	)
}

func TestParamInNestedUdf(t *testing.T) {
	expr := `
FUNC outer() (
  FUNC middle() (
    FUNC inner() => @foo
    RETURN inner()
  )
  RETURN middle()
)
RETURN outer()
`

	RunUseCases(t,
		[]UseCase{
			CaseRuntimeErrorStr(expr, "Missing parameter", "Should report missing parameter used only in nested UDF path"),
			Options(Case(expr, "bar", "Should resolve parameter in nested UDF path when provided"), WithParam("foo", "bar")),
		},
	)
}

func TestParamDifferentInUdf(t *testing.T) {
	expr := `
LET x = @alpha
FUNC f() => @beta
RETURN [x, f()]
`

	RunUseCases(t,
		[]UseCase{
			Case(expr, []any{1, 2}, "Should resolve different parameters in main body and UDF body"),
		},
		WithParam("alpha", 1),
		WithParam("beta", 2),
	)
}
