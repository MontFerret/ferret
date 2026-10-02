package optimization

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestTailCallEliminationEligibility(t *testing.T) {
	r1, r2, r3 := bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(3)
	call := bytecode.NewInstruction(bytecode.OpCall, r3, r1, r2)
	ret := bytecode.NewInstruction(bytecode.OpReturn, r3)
	cases := []struct {
		edit func(*bytecode.Program)
		name string
		code []bytecode.Instruction
		want bool
	}{
		{name: "direct", code: []bytecode.Instruction{call, ret}, want: true},
		{name: "zero arguments", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpCall, r3), ret}, want: true},
		{name: "borrowed capture", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMove, r2, r1), call, ret}, want: true},
		{name: "different result", code: []bytecode.Instruction{call, bytecode.NewInstruction(bytecode.OpReturn, r1)}},
		{name: "distinct result", code: []bytecode.Instruction{call, bytecode.NewInstruction(bytecode.OpDistinct, r2, r3), bytecode.NewInstruction(bytecode.OpReturn, r2)}},
		{name: "required work", code: []bytecode.Instruction{call, bytecode.NewInstruction(bytecode.OpClose, r1), ret}},
		{name: "result copy", code: []bytecode.Instruction{call, bytecode.NewInstruction(bytecode.OpMoveTracked, r2, r3), bytecode.NewInstruction(bytecode.OpReturn, r2)}},
		{name: "protected call", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpProtectedCall, r3, r1, r2), ret}},
		{name: "host call", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpHCall, r3, r1, r2), ret}},
		{name: "main call", code: []bytecode.Instruction{call, ret}, edit: func(p *bytecode.Program) { p.Functions.UserDefined = nil }},
		{name: "function boundary", code: []bytecode.Instruction{call, ret}, edit: func(p *bytecode.Program) {
			p.Functions.UserDefined = append(p.Functions.UserDefined, bytecode.UDF{Entry: 2, Registers: 4})
		}},
		{name: "shared return", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(3), r1), call, ret}},
		{name: "match failure target", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, r1, r2, bytecode.NewConstant(0)), call, ret}, edit: func(p *bytecode.Program) {
			p.Metadata.MatchFailTargets = []int{-1, 3, -1, -1}
		}},
		{name: "catch covers call", code: []bytecode.Instruction{call, ret}, edit: func(p *bytecode.Program) { p.CatchTable = []bytecode.Catch{{1, 1, -1}} }},
		{name: "catch covers return", code: []bytecode.Instruction{call, ret}, edit: func(p *bytecode.Program) { p.CatchTable = []bytecode.Catch{{2, 2, -1}} }},
		{name: "catch handler return", code: []bytecode.Instruction{call, ret}, edit: func(p *bytecode.Program) { p.CatchTable = []bytecode.Catch{{0, 0, 2}} }},
		{name: "local cell", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMakeCell, r1, r2), call, ret}},
		{name: "local cell alias chain", code: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(4), r3),
			bytecode.NewInstruction(bytecode.OpMoveTracked, r2, r1),
			bytecode.NewInstruction(bytecode.OpMove, r1, bytecode.NewRegister(4)), call, ret,
		}},
		{name: "overwritten cell alias", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMakeCell, r1, r2), bytecode.NewInstruction(bytecode.OpLoadNone, r1), call, ret}},
		{name: "cell alias cycle", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMakeCell, r1, r3), bytecode.NewInstruction(bytecode.OpMove, r2, r1), bytecode.NewInstruction(bytecode.OpMoveTracked, r1, r2), call, ret}},
		{name: "local cell not forwarded", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(4), r2), call, ret}, want: true},
		{name: "loaded cell value", code: []bytecode.Instruction{bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(4), r2), bytecode.NewInstruction(bytecode.OpLoadCell, r1, bytecode.NewRegister(4)), call, ret}, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			program := &bytecode.Program{
				Bytecode:  append([]bytecode.Instruction{bytecode.NewInstruction(bytecode.OpReturn)}, tc.code...),
				Functions: bytecode.Functions{UserDefined: []bytecode.UDF{{Entry: 1, Registers: 5, Params: 2}}},
			}
			if tc.edit != nil {
				tc.edit(program)
			}

			cfg, err := NewBuilder(program).Build()
			if err != nil {
				t.Fatal(err)
			}

			result, err := NewTailCallEliminationPass().Run(&PassContext{Program: program, CFG: cfg})
			if err != nil {
				t.Fatal(err)
			}

			if result.Modified != tc.want {
				t.Fatalf("modified = %v, want %v", result.Modified, tc.want)
			}

			for i, instruction := range tc.code {
				if tc.want && instruction.Opcode == bytecode.OpCall {
					instruction.Opcode = bytecode.OpTailCall
				}

				if program.Bytecode[i+1] != instruction {
					t.Fatalf("unexpected instruction at %d: %+v, want %+v", i+1, program.Bytecode[i+1], instruction)
				}
			}
		})
	}
}

func TestTailCallEliminationCellAliasesBelongToEachUDF(t *testing.T) {
	r1, r2, r3 := bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(3)
	call := bytecode.NewInstruction(bytecode.OpCall, r3, r1, r2)
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn),
			bytecode.NewInstruction(bytecode.OpMakeCell, r1, r2),
			call,
			bytecode.NewInstruction(bytecode.OpReturn, r3),
			bytecode.NewInstruction(bytecode.OpMoveTracked, r2, r1),
			call,
			bytecode.NewInstruction(bytecode.OpReturn, r3),
		},
		// Function metadata order does not define bytecode ownership. These two
		// frames reuse register numbers, but only the first creates a cell.
		Functions: bytecode.Functions{UserDefined: []bytecode.UDF{{Entry: 4, Registers: 4}, {Entry: 1, Registers: 4}}},
	}
	cfg, err := NewBuilder(program).Build()
	if err != nil {
		t.Fatal(err)
	}

	result, err := NewTailCallEliminationPass().Run(&PassContext{Program: program, CFG: cfg})
	if err != nil {
		t.Fatal(err)
	}

	if !result.Modified || program.Bytecode[2] != call {
		t.Fatal("the local cell owner must retain its frame")
	}

	call.Opcode = bytecode.OpTailCall
	if program.Bytecode[5] != call {
		t.Fatal("the borrowed cell must remain eligible in its own frame")
	}
}

func TestTailCallEliminationPreservesProgramAndIsIdempotent(t *testing.T) {
	code := []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpReturn),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
	}
	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Bytecode:   code,
		Constants:  []runtime.Value{runtime.NewInt(1), runtime.NewInt(1)},
		Registers:  4,
		Functions: bytecode.Functions{UserDefined: []bytecode.UDF{
			{Name: "forward", Entry: 1, Registers: 4, Params: 2},
			{Name: "recursive", Entry: 4, Registers: 4, Params: 2},
		}},
		CatchTable: []bytecode.Catch{{0, 0, -1}},
		Metadata: bytecode.Metadata{
			Labels:                 map[int]string{1: "forward", 4: "recursive"},
			DebugPoints:            []bytecode.DebugPoint{{PC: 2, FunctionID: 0, Bindings: []bytecode.DebugBinding{{Name: "value", Register: bytecode.NewRegister(1)}}}},
			DebugSpans:             make([]source.Span, len(code)),
			CallArgumentSpans:      make([][]source.Span, len(code)),
			AggregateSelectorSlots: []int{-1, -1, -1, -1, -1, -1, -1},
			MatchFailTargets:       []int{-1, -1, -1, -1, -1, -1, -1},
		},
	}
	program.Metadata.DebugSpans[2] = source.Span{Start: 2, End: 8}
	program.Metadata.CallArgumentSpans[2] = []source.Span{{Start: 3, End: 4}, {Start: 5, End: 6}}
	want := *program
	want.Bytecode = append([]bytecode.Instruction(nil), code...)
	want.Bytecode[2].Opcode = bytecode.OpTailCall
	want.Bytecode[5].Opcode = bytecode.OpTailCall
	wantJSON, err := want.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	for iteration := range 2 {
		cfg, err := NewBuilder(program).Build()
		if err != nil {
			t.Fatal(err)
		}

		result, err := NewTailCallEliminationPass().Run(&PassContext{Program: program, CFG: cfg})
		if err != nil {
			t.Fatal(err)
		}

		if result.Modified != (iteration == 0) {
			t.Fatalf("iteration %d modified = %v", iteration, result.Modified)
		}

		got, err := program.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, wantJSON) {
			t.Fatalf("program changed beyond call opcodes:\n%s\nwant:\n%s", got, wantJSON)
		}
	}
}

func TestTailCallEliminationPipelineRebuildsControlFlow(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn),
			bytecode.NewInstruction(bytecode.OpCall, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Functions: bytecode.Functions{UserDefined: []bytecode.UDF{{Entry: 1, Registers: 2}}},
	}
	pipeline := NewPipeline()
	pipeline.Add(NewTailCallEliminationPass())
	pipeline.Add(NewLivenessAnalysisPass())
	pipeline.Add(newPipelineTestPass("inspect", []string{LivenessAnalysisPassName}, func(ctx *PassContext) (*PassResult, error) {
		blocks := buildInstructionBlockMap(ctx.CFG, len(program.Bytecode))
		if blocks[1] == blocks[2] || blocks[1].Instructions[0].Opcode != bytecode.OpTailCall {
			t.Fatal("subsequent analysis saw stale control flow")
		}

		uses, defs := instructionUseDef(ctx.Program.Bytecode[1])
		if len(uses) != 1 || uses[0] != 1 || len(defs) != 0 {
			t.Fatalf("unexpected tail-call use/def: %v / %v", uses, defs)
		}

		return &PassResult{}, nil
	}))
	if _, err := pipeline.Run(program); err != nil {
		t.Fatal(err)
	}
}
