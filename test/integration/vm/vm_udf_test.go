package vm_test

import "testing"

func TestUDFs(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
FUNC id(x) => x
RETURN id(1)
`, 1, "UDF arrow body"),
		Case(`
FUNC id(x) (
  RETURN x
)
RETURN id(2)
`, 2, "UDF with parentheses"),
		Case(`
LET base = 5
FUNC getBase() => base
RETURN getBase()
`, 5, "Capture global variable"),
		Case(`
FUNC outer(x) (
  FUNC inner(y) (
    RETURN x + y
  )
  RETURN inner(1)
)
RETURN outer(2)
`, 3, "Nested capture"),
		Case(`
FUNC outer(a) (
  FUNC inner(b) (
    RETURN b
  )
  LET v = inner(1)
  RETURN v
)
RETURN outer(2)
`, 1, "Nested LET before return"),
		Case(`
FUNC fact(n) (
  RETURN MATCH n (
    0 => 1,
    _ => n * fact(n - 1),
  )
)
RETURN fact(5)
`, 120, "Recursion"),
		CaseNil(`
FUNC risky() (
  RETURN T::FAIL()
)
RETURN risky()?
`, "Protected UDF call"),
		Case(`
FUNC LENGTH(x) (
  RETURN 42
)
RETURN LENGTH([1,2,3])
`, 42, "UDF shadows builtin"),
		CaseArray(`
FUNC a0() => 0
FUNC a1(a) => a
FUNC a2(a, b) => a + b
FUNC a3(a, b, c) => a + b + c
FUNC a4(a, b, c, d) => a + b + c + d
FUNC a6(a, b, c, d, e, f) => a + b + c + d + e + f
RETURN [a0(), a1(1), a2(1, 2), a3(1, 2, 3), a4(1, 2, 3, 4), a6(1, 2, 3, 4, 5, 6)]
`, []any{0, 1, 3, 6, 10, 21}, "UDF arity coverage"),
		Case(`
FUNC outer(a, b, c, d, e, f) (
  FUNC inner(x, y, z, w, u, v) (
    RETURN x + y + z + w + u + v
  )
  RETURN inner(a, b, c, d, e, f)
)
RETURN outer(1, 2, 3, 4, 5, 6)
`, 21, "Nested frame arity preservation"),
		Case(`
FUNC sum(n, acc) (
  RETURN MATCH n (
    0 => acc,
    _ => sum(n - 1, acc + n),
  )
)
RETURN sum(10, 0)
`, 55, "Tail recursion semantics"),
	})
}
