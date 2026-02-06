package optimization_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestRegisterCoalescing(t *testing.T) {
	RunUseCases(t, compiler.O1, []UseCase{
		ByteCodeCase(`
LET a = 10
LET b = a + 1
LET c = b * 2
LET d = c - 3
RETURN d
`, []vm.Instruction{}),
	})
}
