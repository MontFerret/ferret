package optimization

import (
	"errors"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestPipeline_ReusesFreshDependencyMetadata(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	dependencyRuns := 0
	consumerRuns := 0

	p := NewPipeline()
	p.Add(newPipelineTestPass("dep", nil, func(*PassContext) (*PassResult, error) {
		dependencyRuns++

		return &PassResult{
			Metadata: map[string]any{"runs": dependencyRuns},
		}, nil
	}))
	p.Add(newPipelineTestPass("consumer-a", []string{"dep"}, func(ctx *PassContext) (*PassResult, error) {
		consumerRuns++

		meta, ok := ctx.Metadata["dep"].(map[string]any)
		if !ok {
			t.Fatalf("expected dependency metadata, got %T", ctx.Metadata["dep"])
		}

		if got := meta["runs"]; got != 1 {
			t.Fatalf("expected first consumer to see fresh dependency metadata, got %v", got)
		}

		return &PassResult{}, nil
	}))
	p.Add(newPipelineTestPass("consumer-b", []string{"dep"}, func(ctx *PassContext) (*PassResult, error) {
		consumerRuns++

		meta, ok := ctx.Metadata["dep"].(map[string]any)
		if !ok {
			t.Fatalf("expected dependency metadata, got %T", ctx.Metadata["dep"])
		}

		if got := meta["runs"]; got != 1 {
			t.Fatalf("expected second consumer to reuse fresh dependency metadata, got %v", got)
		}

		return &PassResult{}, nil
	}))

	if _, err := p.Run(program); err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}

	if dependencyRuns != 1 {
		t.Fatalf("expected dependency to run once, got %d", dependencyRuns)
	}

	if consumerRuns != 2 {
		t.Fatalf("expected both consumers to run, got %d", consumerRuns)
	}
}

func TestPipeline_RefreshesStaleDependencyAfterMutation(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	dependencyRuns := 0
	seenRuns := 0

	p := NewPipeline()
	p.Add(newPipelineTestPass("dep", nil, func(*PassContext) (*PassResult, error) {
		dependencyRuns++

		return &PassResult{
			Metadata: map[string]any{"runs": dependencyRuns},
		}, nil
	}))
	p.Add(newPipelineTestPass("mutator", nil, func(ctx *PassContext) (*PassResult, error) {
		// Any mutation forces the pipeline to refresh analysis metadata for the next epoch.
		ctx.Program.Bytecode[0] = bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2))

		return &PassResult{
			Modified: true,
		}, nil
	}))
	p.Add(newPipelineTestPass("consumer", []string{"dep"}, func(ctx *PassContext) (*PassResult, error) {
		meta, ok := ctx.Metadata["dep"].(map[string]any)
		if !ok {
			t.Fatalf("expected dependency metadata, got %T", ctx.Metadata["dep"])
		}

		runs, ok := meta["runs"].(int)
		if !ok {
			t.Fatalf("expected dependency run count, got %T", meta["runs"])
		}

		seenRuns = runs

		return &PassResult{}, nil
	}))

	result, err := p.Run(program)
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}

	if !result.Modified {
		t.Fatalf("expected pipeline to report a modification")
	}

	if dependencyRuns != 2 {
		t.Fatalf("expected dependency to rerun after mutation, got %d runs", dependencyRuns)
	}

	if seenRuns != 2 {
		t.Fatalf("expected consumer to see refreshed dependency metadata, got %d", seenRuns)
	}
}

func TestPipeline_FailsWhenDependencyMissing(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	p := NewPipeline()
	p.Add(newPipelineTestPass("consumer", []string{"dep"}, func(*PassContext) (*PassResult, error) {
		return &PassResult{}, nil
	}))

	_, err := p.Run(program)
	if !errors.Is(err, ErrMissingDependency) {
		t.Fatalf("expected missing dependency error, got %v", err)
	}
}

func TestPipeline_FailsWhenRefreshedDependencyModifiesProgram(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	dependencyRuns := 0

	p := NewPipeline()
	p.Add(newPipelineTestPass("dep", nil, func(*PassContext) (*PassResult, error) {
		dependencyRuns++

		return &PassResult{
			Modified: dependencyRuns > 1,
			Metadata: map[string]any{"runs": dependencyRuns},
		}, nil
	}))
	p.Add(newPipelineTestPass("mutator", nil, func(ctx *PassContext) (*PassResult, error) {
		ctx.Program.Bytecode[0] = bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2))

		return &PassResult{
			Modified: true,
		}, nil
	}))
	p.Add(newPipelineTestPass("consumer", []string{"dep"}, func(*PassContext) (*PassResult, error) {
		return &PassResult{}, nil
	}))

	_, err := p.Run(program)
	if !errors.Is(err, ErrDependencyRefreshModified) {
		t.Fatalf("expected dependency refresh modify error, got %v", err)
	}
}

func TestPipeline_RefreshesLivenessBeforePeepholeAfterRegisterCoalescing(t *testing.T) {
	newProgram := func() *bytecode.Program {
		return &bytecode.Program{
			Registers: 7,
			Constants: []runtime.Value{
				runtime.Int(9),
				runtime.Boolean(true),
			},
			Bytecode: []bytecode.Instruction{
				bytecode.NewInstruction(bytecode.OpLoadBool, bytecode.NewRegister(6), bytecode.Operand(1)),
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(5), bytecode.NewConstant(0)),
				bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(4), bytecode.NewRegister(3), bytecode.NewRegister(5)),
				bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(5), bytecode.NewRegister(6)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(5)),
			},
		}
	}

	coalesced := newProgram()
	runCoalescing(t, coalesced)

	full := newProgram()
	p := NewPipeline()
	p.Add(NewLivenessAnalysisPass())
	p.Add(NewRegisterCoalescingPass())
	p.Add(NewPeepholePass())

	if _, err := p.Run(full); err != nil {
		t.Fatalf("full optimization pipeline failed: %v", err)
	}

	if !reflect.DeepEqual(full.Bytecode, coalesced.Bytecode) {
		t.Fatalf("expected peephole to keep post-coalescing bytecode unchanged when temp stays live across the branch:\ncoalesced: %#v\nfull: %#v", coalesced.Bytecode, full.Bytecode)
	}

	if len(full.Bytecode) != 6 {
		t.Fatalf("expected peephole to keep all instructions, got %d", len(full.Bytecode))
	}

	if full.Bytecode[1].Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected LOADC to remain after peephole, got %s", full.Bytecode[1].Opcode)
	}

	if full.Bytecode[2].Opcode != bytecode.OpAdd {
		t.Fatalf("expected ADD to remain after peephole, got %s", full.Bytecode[2].Opcode)
	}
}
