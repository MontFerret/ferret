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
`, "FOR n WHILE i < 2\n    RETURN n"),
		S(`
FOR WHILE i < 1
	DISPATCH "click" IN @d
	RETURN i
	`, "FOR WHILE i < 1\n    DISPATCH \"click\" IN @d\n    RETURN i"),
	})
}
