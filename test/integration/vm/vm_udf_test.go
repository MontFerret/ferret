package vm_test

import "testing"

func TestUDFs(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
FUNC id(x)
  RETURN x
RETURN id(1)
`, 1, "UDF without parentheses"),
		Case(`
FUNC id(x) (
  RETURN x
)
RETURN id(2)
`, 2, "UDF with parentheses"),
		Case(`
LET base = 5
FUNC getBase()
  RETURN base
RETURN getBase()
`, 5, "Capture global variable"),
		Case(`
FUNC outer(x)
  FUNC inner(y) (
    RETURN x + y
  )
  RETURN inner(1)
RETURN outer(2)
`, 3, "Nested capture"),
		Case(`
FUNC fact(n)
  RETURN MATCH n (
    0 => 1,
    _ => n * fact(n - 1),
  )
RETURN fact(5)
`, 120, "Recursion"),
		CaseNil(`
FUNC risky()
  RETURN T::FAIL()
RETURN risky()?
`, "Protected UDF call"),
		Case(`
FUNC LENGTH(x)
  RETURN 42
RETURN LENGTH([1,2,3])
`, 42, "UDF shadows builtin"),
	})
}
