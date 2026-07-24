package spec

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

type T interface {
	Helper()
	Fatalf(format string, args ...any)
}

func NewCompiler(t T, opts ...compiler.Option) *compiler.Compiler {
	t.Helper()

	c, err := compiler.New(opts...)
	if err != nil {
		t.Fatalf("failed to create compiler: %v", err)
	}

	return c
}
