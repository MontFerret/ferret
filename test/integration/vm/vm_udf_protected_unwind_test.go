package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUDFProtectedUnwindDeepFrames(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
FUNC ok() => 9
FUNC bomb() => T::FAIL()
FUNC level2() => bomb()
FUNC level1() => level2()
FUNC guard() (
  LET failed = level1()?
  LET safe = ok()
  RETURN [failed, safe]
)
RETURN guard()
`, []any{nil, 9}, "Protected unwind drops nested frames and resumes caller"),
		Nil(`
FUNC dive(n) (
  RETURN MATCH n (
    0 => T::FAIL(),
    _ => dive(n - 1),
  )
)
RETURN dive(8)?
`, "Protected unwind across deep recursion"),
		Array(`
FUNC bomb() => T::FAIL()
FUNC level2() => bomb()
FUNC level1() => level2()
FUNC once() => level1()?
RETURN [once(), once(), once()]
`, []any{nil, nil, nil}, "Repeated protected unwinds in one run"),
		Array(`
FUNC outer() (
  VAR total = 1
  FUNC wrapper(v) (
    FUNC setAndFail(next) (
      total = next
      RETURN T::FAIL()
    )
    RETURN setAndFail(v)
  )

  LET failed = wrapper(9)?
  RETURN [failed, total]
)
RETURN outer()
`, []any{nil, 9}, "Protected unwind preserves updated promoted VAR state"),
	})
}
