package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func mustNewCompiler(t testing.TB, options ...compiler.Option) *compiler.Compiler {
	t.Helper()

	compilerInstance, err := compiler.New(options...)
	if err != nil {
		t.Fatalf("create compiler: %v", err)
	}

	return compilerInstance
}

func compileWithLevel(t *testing.T, level compiler.OptimizationLevel, expr string) *bytecode.Program {
	t.Helper()

	c := mustNewCompiler(t, compiler.WithOptimizationLevel(level))
	prog, err := c.Compile(source.NewAnonymous(expr))
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
