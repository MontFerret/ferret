package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func compileWithLevel(t *testing.T, level compiler.OptimizationLevel, expr string) *bytecode.Program {
	t.Helper()

	c := compiler.New(compiler.WithOptimizationLevel(level))
	prog, err := c.Compile(file.NewAnonymousSource(expr))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	return prog
}

func firstCompilationError(err error) *diagnostics.Diagnostic {
	switch e := err.(type) {
	case *diagnostics.Diagnostic:
		return e
	case *diagnostics.DiagnosticSet:
		if e.Size() == 0 {
			return nil
		}

		return e.First()
	default:
		return nil
	}
}
