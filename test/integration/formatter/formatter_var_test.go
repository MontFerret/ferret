package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterVarBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
	VAR    foo =      10
	 foo    =   foo +   1
	 RETURN foo
		`, `VAR foo = 10
foo = foo + 1
RETURN foo`),
		S(`
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
		S(`
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
		S(`
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
		S(`
LET    STEP =  10
RETURN STEP
`, `LET STEP = 10
RETURN STEP`),
		S(`
FOR item IN [ 1, 2 ]
VAR current = item
current=current+1
RETURN current
`, `FOR item IN [1, 2]
    VAR current = item
    current = current + 1
    RETURN current`),
		S(`
LET item={ deprecated: true, keep: true}
DELETE    item.deprecated
RETURN item
`, `LET item = { deprecated: true, keep: true }
DELETE item.deprecated
RETURN item`),
		S(`
FUNC clean(payload)(
DELETE  payload["debug"]
RETURN payload
)
RETURN clean({})
`, `FUNC clean(payload) (
    DELETE payload["debug"]
    RETURN payload
)
RETURN clean({})`),
		S(`
FOR item IN [{ stale: true }]
DELETE item.stale
RETURN item
`, `FOR item IN [{ stale: true }]
    DELETE item.stale
    RETURN item`),
	})
}
