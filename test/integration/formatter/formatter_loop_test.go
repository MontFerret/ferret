package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterLoopBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
FOR _ WHILE i < 2
RETURN i
`, `FOR WHILE i < 2
    RETURN i`),
		S(`
FOR _ DO WHILE false
RETURN 1
`, `FOR DO WHILE FALSE
    RETURN 1`),
		S(`
FOR n WHILE i < 2
RETURN n
`, `FOR n WHILE i < 2
    RETURN n`),
	})
}
