package compiler_test

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	"testing"
)

func Disassembly(instr []string, opcodes ...vm.Opcode) string {
	var disassembly string

	for i := 0; i < len(instr); i++ {
		disassembly += fmt.Sprintf("%d: [%d] %s\n", i, opcodes[i], instr[i])
	}

	return disassembly
}

func TestCompiler_Variables(t *testing.T) {
	RunAsmUseCases(t, []ByteCodeUseCase{
		{
			`LET i = NONE RETURN i`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADN R1
1: [%d] STOREG C0 R1
2: [%d] LOADG R2 C0
3: [%d] MOVE R0 R2
4: [%d] RET
`,
					vm.OpLoadNone,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("i"),
				},
			},
		},
		{
			`LET a = TRUE RETURN a`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADB R1 1
1: [%d] STOREG C0 R1
2: [%d] LOADG R2 C0
3: [%d] MOVE R0 R2
4: [%d] RET
`,
					vm.OpLoadBool,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("a"),
				},
			},
		},
		{
			`LET a = FALSE RETURN a`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADB R1 0
1: [%d] STOREG C0 R1
2: [%d] LOADG R2 C0
3: [%d] MOVE R0 R2
4: [%d] RET
`,
					vm.OpLoadBool,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("a"),
				},
			},
		},
		{
			`LET a = 1.1 RETURN a`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADC R1 C0
1: [%d] STOREG C1 R1
2: [%d] LOADG R2 C1
3: [%d] MOVE R0 R2
4: [%d] RET
`,
					vm.OpLoadConst,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewFloat(1.1),
					runtime.NewString("a"),
				},
			},
		},
		{
			`
LET a = 'foo'
LET b = a
RETURN a
`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADC R1 C0
1: [%d] STOREG C1 R1
2: [%d] LOADG R2 C1
3: [%d] STOREG C2 R2
4: [%d] LOADG R3 C2
5: [%d] MOVE R0 R3
6: [%d] RET
`,
					vm.OpLoadConst,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpStoreGlobal,
					vm.OpLoadGlobal,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("foo"),
					runtime.NewString("a"),
					runtime.NewString("b"),
				},
			},
		},
	})
}

func TestCompiler_FuncCall(t *testing.T) {
	RunAsmUseCases(t, []ByteCodeUseCase{
		{
			`RETURN FOO()`,
			ExpectedProgram{
				Disassembly: fmt.Sprintf(`
0: [%d] LOADC R1 C0
1: [%d] CALL R1
2: [%d] MOVE R0 R1
3: [%d] RET
`,
					vm.OpLoadConst,
					vm.OpCall,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("FOO"),
				},
			},
		},
		{
			`RETURN FOO("a", 1, TRUE)`,
			ExpectedProgram{
				Disassembly: Disassembly([]string{
					"LOADC R1 C0",
					"LOADC R2 C1",
					"LOADB R3 1",
					"CALL R1 R2 R3",
					"MOVE R0 R1",
					"RET",
				},
					vm.OpLoadConst,
					vm.OpLoadConst,
					vm.OpLoadBool,
					vm.OpCall,
					vm.OpMove,
					vm.OpReturn,
				),
				Constants: []runtime.Value{
					runtime.NewString("FOO"),
				},
			},
		},
	})
}
