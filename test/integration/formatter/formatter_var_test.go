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
	VAR total=10
	total+=1
	total-=2
	total*=3
	total/=3
	RETURN total
		`, `VAR total = 10
total += 1
total -= 2
total *= 3
total /= 3
RETURN total`),
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
	FUNC run()(
	VAR total=10
	total+=1
	total-=2
	total*=3
	total/=3
	RETURN total
	)
	RETURN run()
		`, `FUNC run() (
    VAR total = 10
    total += 1
    total -= 2
    total *= 3
    total /= 3
    RETURN total
)
RETURN run()`),
		Case(`
LET    STEP =  10
RETURN STEP
`, `LET STEP = 10
RETURN STEP`),
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
