package formatter_test

import "testing"

func TestFormatterVarBindings(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
VAR    foo =      10
 foo    =   foo +   1
 RETURN foo
`, `VAR foo = 10
foo = foo + 1
RETURN foo`),
		Case(`
FUNC   run( )(
VAR total= 1
 total   =total+2
 RETURN total
)
RETURN run()
`, `FUNC run() (
    VAR total = 1
    total = total + 2
    RETURN total
)
RETURN run()`),
		Case(`
FOR item IN [ 1, 2 ]
VAR current = item
current=current+1
RETURN current
`, `FOR item IN [1, 2]
    VAR current = item
    current = current + 1
    RETURN current`),
	})
}
