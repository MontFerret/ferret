package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
)

func TestRegisterCoalescing(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		AtMostRegistersCase(`
LET a = 10
LET b = a + 1
LET c = b * 2
LET d = c - 3
RETURN d
`, 3),
		AtMostRegistersCase(`
LET a = 10
LET b = a + 1
LET c = b * 2
LET d = c - 3
RETURN d
`, 3),
		AtMostRegistersCase(`
LET a = 10
LET b = a
LET c = b + 1
RETURN c
`, 3),
		AtMostRegistersCase(`
LET a = 1
LET b = a
LET c = b
RETURN c
`, 1),
		AtMostRegistersCase(`
LET a = 1
LET b = 2
LET c = a + b
RETURN c
`, 3),
		AtMostRegistersCase(`
LET a = 10
LET arr = [a, a + 1, a + 2, a + 3]
RETURN arr
`, 3, "Flat array literal with expression elems"),
	})
}
