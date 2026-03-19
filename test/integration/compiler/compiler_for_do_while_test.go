package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func TestForDoWhileCompiles(t *testing.T) {
	expressions := []string{
		`
FOR DO WHILE false
	RETURN 1
`,
		`
FOR i DO WHILE false
	RETURN i
`,
		`
FOR _ DO WHILE false
	RETURN 1
`,
	}

	for _, expr := range expressions {
		for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
			_ = compileWithLevel(t, level, expr)
		}
	}
}
