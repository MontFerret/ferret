package spec

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/asm"
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func PrintDebug(t *testing.T, name string, prog *bytecode.Program) {
	t.Helper()
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Suite: %s\n", name))
	b.WriteString(fmt.Sprintf("Expression: %s\n", prog.Source.Content()))
	b.WriteString("Bytecode:")

	out, e := asm.Disassemble(prog, asm.WithDebug())

	if e == nil {
		b.WriteString(out)
	}

	t.Log(b.String())
}

func PrintError(t *testing.T, err error) {
	t.Helper()

	t.Log("\n" + diagnostics.Format(err))
}
