package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestHostCallBindingNameConflicts(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		Failure(`
foo()
LET foo = "ff"
RETURN foo
`, expectedHostCallBindingCollision("foo"), "Host call before LET with the same name should fail during compilation"),
		Failure(`
LET foo = "ff"
RETURN foo()
`, expectedHostCallBindingCollision("foo"), "Host call after LET with the same name should fail during compilation"),
		Failure(`
FUNC f() (
  LET x = foo()
  LET foo = 1
  RETURN foo
)
RETURN f()
`, expectedHostCallBindingCollision("foo"), "Host call and LET collision inside UDF body should fail during compilation"),
		Failure(`
USE X::bar AS foo
foo()
LET foo = 1
RETURN foo
`, expectedHostCallBindingCollision("foo"), "Function alias source name should collide with same-scope binding"),
		Failure(`
FOR foo IN [1]
  RETURN foo()
`, expectedHostCallBindingCollision("foo"), "Loop value binding should collide with same-scope host call"),
		Failure(`
FOR item, foo IN [1]
  RETURN foo()
`, expectedHostCallBindingCollision("foo"), "Loop counter binding should collide with same-scope host call"),
		Failure(`
FOR item IN [1]
  COLLECT foo = item
  RETURN foo()
`, expectedHostCallBindingCollision("foo"), "COLLECT grouping output should collide with same-scope host call"),
		ProgramCheck(`
foo()
LET values = (
  FOR foo IN [1]
    RETURN foo
)
RETURN values
`, compileOnly, "Inner loop binding may shadow a host call name from an outer scope"),
		ProgramCheck(`
LET foo = 1
LET values = (
  FOR item IN [1]
    RETURN foo()
)
RETURN foo
`, compileOnly, "Nested host call should not collide with an outer binding"),
		ProgramCheck(`
USE X AS NS
NS::foo()
LET foo = 1
RETURN foo
`, compileOnly, "Direct namespaced host call should not collide with local binding"),
		ProgramCheck(`
foo()
RETURN 1
`, compileOnly, "Unresolved host call without binding collision should remain a runtime warmup concern"),
	}, compiler.O0, compiler.O1)
}

func expectedHostCallBindingCollision(name string) E {
	return E{
		Kind:    parserd.NameError,
		Message: fmt.Sprintf("Variable '%s' conflicts with function call '%s'", name, name),
		Hint:    fmt.Sprintf("Rename either the variable '%s' or the function call '%s' to make the reference unambiguous.", name, name),
	}
}

func compileOnly(*bytecode.Program) error {
	return nil
}
