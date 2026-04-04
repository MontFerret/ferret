package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUDFs(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`
FUNC id(x) => x
RETURN id(1)
`, 1, "UDF arrow body"),
		S(`
FUNC id(x) (
  RETURN x
)
RETURN id(2)
`, 2, "UDF with parentheses"),
		S(`
LET base = 5
FUNC getBase() => base
RETURN getBase()
`, 5, "Capture global variable"),
		S(`
FUNC outer(x) (
  FUNC inner(y) (
    RETURN x + y
  )
  RETURN inner(1)
)
RETURN outer(2)
`, 3, "Nested capture"),
		S(`
LET global = 100
FUNC outer(a) (
  LET outerLocal = 10
  FUNC middle(b) (
    FUNC inner(c) => global + a + outerLocal + b + c
    RETURN inner(1)
  )
  RETURN middle(2)
)
RETURN outer(3)
`, 116, "Multi-level capture propagation"),
		S(`
FUNC outer(a) (
  FUNC inner(b) (
    RETURN b
  )
  LET v = inner(1)
  RETURN v
)
RETURN outer(2)
`, 1, "Nested LET before return"),
		S(`
FUNC fact(n) (
  RETURN MATCH n (
    0 => 1,
    _ => n * fact(n - 1),
  )
)
RETURN fact(5)
`, 120, "Recursion"),
		Array(`
FUNC f() => "outer"
FUNC outer() (
  FUNC f() => "inner"
  RETURN f()
)
RETURN [outer(), f()]
`, []any{"inner", "outer"}, "Nested UDF shadows only within lexical scope"),
		Array(`
LET value = 1
FUNC outer() (
  LET value = 10
  FUNC inner() => value
  RETURN [inner(), value]
)
RETURN [outer(), value]
`, []any{[]any{10, 10}, 1}, "Nested UDF captures nearest shadowed local"),
		Nil(`
FUNC risky() (
  RETURN T::FAIL()
)
RETURN risky()?
`, "Protected UDF call"),
		S(`
FUNC LENGTH(x) (
  RETURN 42
)
RETURN LENGTH([1,2,3])
`, 42, "UDF shadows builtin"),
		S(`
FUNC a() => 1
FUNC A() => 2
RETURN a() + A()
`, 3, "UDF lookup is case-sensitive"),
		S(`
FUNC length() => 100
RETURN LENGTH([1,2,3]) + length()
`, 103, "Lowercase UDF does not shadow differently-cased builtin call"),
		S(`
FUNC LENGTH(x) => 100
RETURN length([1,2,3])
`, 3, "Builtin survives differently-cased UDF declaration"),
		Array(`
FUNC f() => "outer-lower"
FUNC outer() (
  FUNC F() => "inner-upper"
  RETURN [f(), F()]
)
RETURN [outer(), f()]
`, []any{[]any{"outer-lower", "inner-upper"}, "outer-lower"}, "Nested UDF shadowing is exact-case"),
		Array(`
FUNC a0() => 0
FUNC a1(a) => a
FUNC a2(a, b) => a + b
FUNC a3(a, b, c) => a + b + c
FUNC a4(a, b, c, d) => a + b + c + d
FUNC a6(a, b, c, d, e, f) => a + b + c + d + e + f
RETURN [a0(), a1(1), a2(1, 2), a3(1, 2, 3), a4(1, 2, 3, 4), a6(1, 2, 3, 4, 5, 6)]
`, []any{0, 1, 3, 6, 10, 21}, "UDF arity coverage"),
		S(`
FUNC outer(a, b, c, d, e, f) (
  FUNC inner(x, y, z, w, u, v) (
    RETURN x + y + z + w + u + v
  )
  RETURN inner(a, b, c, d, e, f)
)
RETURN outer(1, 2, 3, 4, 5, 6)
`, 21, "Nested frame arity preservation"),
		S(`
FUNC sum(n, acc) (
  RETURN MATCH n (
    0 => acc,
    _ => sum(n - 1, acc + n),
  )
)
RETURN sum(10, 0)
`, 55, "Tail recursion semantics"),
		S(`
LET base = 10
FUNC outer(seed) (
  FUNC loop(n, acc) (
    RETURN MATCH n (
      0 => acc + base + seed,
      _ => loop(n - 1, acc + n),
    )
  )
  RETURN loop(4, 0)
)
RETURN outer(1)
`, 21, "Nested UDF captures survive tail-recursive paths"),
	})
}
