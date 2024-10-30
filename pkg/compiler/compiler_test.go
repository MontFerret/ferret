package compiler_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"testing"
)

func TestCompiler_Variables(t *testing.T) {
	RunAsmUseCases(t, []ByteCodeUseCase{
		{
			`LET i = NONE RETURN i`,
			ExpectedProgram{
				Disassembly: `
0: [1] LOADK R1 C0
1: [2] STOREG C1 R1
2: [3] LOADG R2 C1
3: [0] MOVE R0 R2
4: [58] RET
`,
				Constants: []core.Value{
					values.None,
					values.NewString("i"),
				},
			},
		},
	})
}
