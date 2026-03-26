package vm_test

import (
	"testing"

	spec "github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestParam(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ErrorStr(`RETURN @foo`, "Missing parameter"),
		ErrorStr(`
FUNC read() => @foo
RETURN read()
`, "Missing parameter"),
		S(`RETURN @str`, "bar", "Should return a value of a parameter"),
		S(`RETURN @int + @int`, 2, "Should return a sum of two parameters"),
		S(`RETURN @obj.str1 + @obj.str2`, "foobar", "Should return a concatenated string of two parameter properties"),
		Array(`FOR i IN @values1 RETURN i`, []any{1, 2, 3, 4}, "Should iterate over an array parameter"),
		Array(`FOR i IN @values2 SORT i RETURN i`, []any{"a", "b", "c", "d"}, "Should iterate over an object parameter"),
		Array(`FOR i IN @start..@end RETURN i`, []any{1, 2, 3, 4, 5}, "Should iterate over a range parameter"),
		S(`RETURN @obj.str1`, "foo", "Should be possible to use in member expression"),
	},
		spec.WithParam("str", "bar"),
		spec.WithParam("int", 1),
		spec.WithParam("bool", true),
		spec.WithParam("obj", map[string]interface{}{"str1": "foo", "str2": "bar"}),
		spec.WithParam("values1", []int{1, 2, 3, 4}),
		spec.WithParam("values2", map[string]interface{}{"a": "a", "b": "b", "c": "c", "d": "d"}),
		spec.WithParam("start", 1),
		spec.WithParam("end", 5),
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

	RunSpecs(t, []spec.Spec{
		ErrorStr(expr, "Missing parameter", "Should report missing parameter used only in nested UDF path"),
		S(expr, "bar", "Should resolve parameter in nested UDF path when provided").Env(spec.WithParam("foo", "bar")),
	},
	)
}

func TestParamUdfSlotAlignment(t *testing.T) {
	expr1 := `
LET x = @alpha
FUNC f() => @beta
RETURN x + f()
`

	RunSpecs(t, []spec.Spec{
		S(expr1, 30, "Should keep UDF @param slots aligned with program param ordering"),
	},
		spec.WithParam("alpha", 10),
		spec.WithParam("beta", 20),
	)

	expr2 := `
LET x = @alpha
FUNC f() => @beta
RETURN [x, f()]
`

	RunSpecs(t, []spec.Spec{
		Array(expr2, []any{1, 2}, "Should resolve different parameters in main body and UDF body"),
	},
		spec.WithParam("alpha", 1),
		spec.WithParam("beta", 2),
	)
}
