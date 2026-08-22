package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

func mustNewCompiler(t testing.TB, options ...compiler.Option) *compiler.Compiler {
	t.Helper()

	compilerInstance, err := compiler.New(options...)
	if err != nil {
		t.Fatalf("create compiler: %v", err)
	}

	return compilerInstance
}
