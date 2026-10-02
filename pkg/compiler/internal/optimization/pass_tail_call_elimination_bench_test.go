package optimization

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkTailCallEliminationPass(b *testing.B) {
	for _, candidates := range []int{32, 128, 512} {
		b.Run(fmt.Sprintf("Plain/UDFs%d", candidates), func(b *testing.B) {
			benchmarkTailCallEliminationPass(b, candidates, 0, false)
		})
		b.Run(fmt.Sprintf("Cell8Aliases/UDFs%d", candidates), func(b *testing.B) {
			benchmarkTailCallEliminationPass(b, candidates, 0, true)
		})
	}

	for _, catches := range []int{0, 128, 512, 2048} {
		b.Run(fmt.Sprintf("CatchScaling/Candidates128/Catches%d", catches), func(b *testing.B) {
			benchmarkTailCallEliminationPass(b, 128, catches, false)
		})
	}

	for _, candidates := range []int{32, 128, 512} {
		b.Run(fmt.Sprintf("CandidateScaling/Catches1024/Candidates%d", candidates), func(b *testing.B) {
			benchmarkTailCallEliminationPass(b, candidates, 1024, false)
		})
	}
}

func benchmarkTailCallEliminationPass(b *testing.B, candidates, catches int, cells bool) {
	b.Helper()
	b.StopTimer()

	program := newTailCallEliminationBenchmarkProgram(candidates, catches, cells)
	if err := bytecode.ValidateProgram(program); err != nil {
		b.Fatal(err)
	}

	cfg, err := NewBuilder(program).Build()
	if err != nil {
		b.Fatal(err)
	}

	pristine := append([]bytecode.Instruction(nil), program.Bytecode...)
	pass := NewTailCallEliminationPass()
	ctx := &PassContext{Program: program, CFG: cfg}
	result, err := pass.Run(ctx)
	if err != nil {
		b.Fatal(err)
	}

	if !result.Modified {
		b.Fatal("benchmark setup has no eligible tail calls")
	}

	assertTailCallEliminationBenchmarkResult(b, program, pristine, candidates)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		// Restore ordinary call/return input outside timing. The original CFG
		// still describes that input; only the program opcodes were changed.
		copy(program.Bytecode, pristine)
		b.StartTimer()
		result, err := pass.Run(ctx)
		if err != nil {
			b.Fatal(err)
		}

		if !result.Modified {
			b.Fatal("benchmark invocation did not transform fresh input")
		}

		b.StopTimer()
	}

	assertTailCallEliminationBenchmarkResult(b, program, pristine, candidates)
}

func assertTailCallEliminationBenchmarkResult(b *testing.B, program *bytecode.Program, pristine []bytecode.Instruction, candidates int) {
	b.Helper()

	if err := bytecode.ValidateProgram(program); err != nil {
		b.Fatal(err)
	}

	tails := 0
	for pc, instruction := range program.Bytecode {
		want := pristine[pc]
		if instruction.Opcode == bytecode.OpTailCall {
			tails++
			if want.Opcode != bytecode.OpCall {
				b.Fatalf("tail call at %d did not originate from an ordinary call", pc)
			}

			if pc+1 >= len(pristine) || pristine[pc+1].Opcode != bytecode.OpReturn || pristine[pc+1].Operands[0] != want.Operands[0] {
				b.Fatalf("tail call at %d has no retained matching return", pc)
			}

			want.Opcode = bytecode.OpTailCall
		}

		if instruction != want {
			b.Fatalf("unexpected benchmark transformation at %d: got %+v, want %+v", pc, instruction, want)
		}
	}

	if tails != candidates {
		b.Fatalf("tail calls = %d, want %d", tails, candidates)
	}
}

func newTailCallEliminationBenchmarkProgram(candidates, catches int, cells bool) *bytecode.Program {
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Registers:  13,
		Constants:  make([]runtime.Value, candidates+1),
		Functions: bytecode.Functions{
			Host: []bytecode.HostFunction{{Name: "STEP", ArgCount: 0}},
		},
	}
	for i := range program.Constants {
		program.Constants[i] = runtime.NewInt(i)
	}

	// Recovery regions belong to the main caller, outside every candidate.
	// Each candidate therefore scans the complete catch table without finding
	// coverage, exposing the per-candidate cost separately from CFG building.
	for range catches {
		program.Bytecode = append(program.Bytecode,
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpProtectedHCall, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
		)
		pc := len(program.Bytecode) - 2
		program.CatchTable = append(program.CatchTable, bytecode.Catch{pc, pc, pc + 1})
	}

	program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(2)))
	for i := range candidates {
		program.Bytecode = append(program.Bytecode,
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(i+1)),
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(3), bytecode.NewRegister(2), bytecode.NewRegister(2)),
		)
	}

	program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)))
	program.Functions.UserDefined = append(program.Functions.UserDefined, bytecode.UDF{
		Name: "leaf", Entry: len(program.Bytecode), Registers: 2, Params: 1,
	})
	program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)))

	for i := range candidates {
		registers := 4
		if cells {
			registers = 13
		}

		program.Functions.UserDefined = append(program.Functions.UserDefined, bytecode.UDF{
			Name: fmt.Sprintf("forward%d", i), Entry: len(program.Bytecode), Registers: registers, Params: 1,
		})
		argument := bytecode.NewRegister(2)
		function := bytecode.NewRegister(3)
		if cells {
			program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(2), bytecode.NewRegister(1)))
			for alias := 3; alias <= 10; alias++ {
				opcode := bytecode.OpMove
				if alias%2 == 0 {
					opcode = bytecode.OpMoveTracked
				}

				program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(opcode, bytecode.NewRegister(alias), bytecode.NewRegister(alias-1)))
			}

			// Forward the cell value, preserving eligibility while exercising
			// eight possible aliases of a locally owned cell handle.
			argument = bytecode.NewRegister(11)
			function = bytecode.NewRegister(12)
			program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpLoadCell, argument, bytecode.NewRegister(10)))
		} else {
			program.Bytecode = append(program.Bytecode, bytecode.NewInstruction(bytecode.OpMoveTracked, argument, bytecode.NewRegister(1)))
		}

		program.Bytecode = append(program.Bytecode,
			bytecode.NewInstruction(bytecode.OpLoadConst, function, bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpCall, function, argument, argument),
			bytecode.NewInstruction(bytecode.OpReturn, function),
		)
	}

	return program
}
