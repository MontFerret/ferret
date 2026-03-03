package vm_test

import "testing"

func TestUDFProtectedUnwindDeepFrames(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
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
		CaseNil(`
FUNC dive(n) (
  RETURN MATCH n (
    0 => T::FAIL(),
    _ => dive(n - 1),
  )
)
RETURN dive(8)?
`, "Protected unwind across deep recursion"),
		CaseArray(`
FUNC bomb() => T::FAIL()
FUNC level2() => bomb()
FUNC level1() => level2()
FUNC once() => level1()?
RETURN [once(), once(), once()]
`, []any{nil, nil, nil}, "Repeated protected unwinds in one run"),
	})
}
