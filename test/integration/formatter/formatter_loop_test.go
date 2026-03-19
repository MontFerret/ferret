package formatter_test

import "testing"

func TestFormatterLoopBindings(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
FOR _ WHILE i < 2
RETURN i
`, `FOR WHILE i < 2
    RETURN i`),
		Case(`
FOR _ DO WHILE false
RETURN 1
`, `FOR DO WHILE FALSE
    RETURN 1`),
		Case(`
FOR n WHILE i < 2
RETURN n
`, `FOR n WHILE i < 2
    RETURN n`),
	})
}
