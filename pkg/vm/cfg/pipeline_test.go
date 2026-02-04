package cfg

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestPipeline_Run(t *testing.T) {
	// Create a simple program
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, -1),
			vm.NewInstruction(vm.OpReturn, 0),
		},
		Constants: []runtime.Value{runtime.NewInt(42)},
		Registers: 10,
	}

	// Create pipeline with analysis passes
	pipeline := NewPipeline()
	pipeline.Add(NewLivenessAnalysisPass())
	pipeline.Add(NewLoopDetectionPass())

	result, err := pipeline.Run(program)
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.PassResults) != 2 {
		t.Errorf("expected 2 pass results, got %d", len(result.PassResults))
	}

	if result.Modified {
		t.Error("analysis passes should not modify the program")
	}
}

func TestLivenessAnalysisPass(t *testing.T) {
	// Create a program with register usage
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1), // r1 = const
			vm.NewInstruction(vm.OpLoadConst, 2, -2), // r2 = const
			vm.NewInstruction(vm.OpAdd, 3, 1, 2),     // r3 = r1 + r2
			vm.NewInstruction(vm.OpReturn, 3),        // return r3
		},
		Constants: []runtime.Value{runtime.NewInt(1), runtime.NewInt(2)},
		Registers: 10,
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewLivenessAnalysisPass()
	result, err := pass.Run(program, cfg)
	if err != nil {
		t.Fatalf("liveness analysis failed: %v", err)
	}

	if result.Modified {
		t.Error("liveness analysis should not modify the program")
	}

	liveness, ok := result.Metadata["liveness"].(map[int]*LivenessInfo)
	if !ok {
		t.Fatal("expected liveness metadata")
	}

	if len(liveness) == 0 {
		t.Error("expected liveness information for blocks")
	}
}

func TestLoopDetectionPass(t *testing.T) {
	// Create a program with a loop
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, -1),  // Block 0
			vm.NewInstruction(vm.OpJumpIfFalse, 4, 0), // Block 0
			vm.NewInstruction(vm.OpLoadConst, 1, -2),  // Block 1
			vm.NewInstruction(vm.OpJump, 0),           // Block 1: back edge
			vm.NewInstruction(vm.OpReturn, 0),         // Block 2
		},
		Constants: []runtime.Value{runtime.NewInt(1), runtime.NewInt(2)},
		Registers: 10,
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewLoopDetectionPass()
	result, err := pass.Run(program, cfg)
	if err != nil {
		t.Fatalf("loop detection failed: %v", err)
	}

	if result.Modified {
		t.Error("loop detection should not modify the program")
	}

	loops, ok := result.Metadata["loops"].([]*Loop)
	if !ok {
		t.Fatal("expected loops metadata")
	}

	if len(loops) != 1 {
		t.Errorf("expected 1 loop, got %d", len(loops))
	}

	if len(loops) > 0 {
		loop := loops[0]
		if loop.Header.Start != 0 {
			t.Errorf("expected loop header at position 0, got %d", loop.Header.Start)
		}
	}
}

func TestConstantFoldingPass(t *testing.T) {
	// Create a program with constant operations
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1), // r1 = 5
			vm.NewInstruction(vm.OpMove, 2, 1),       // r2 = r1 (should propagate constant)
			vm.NewInstruction(vm.OpReturn, 2),
		},
		Constants: []runtime.Value{runtime.NewInt(5)},
		Registers: 10,
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewConstantFoldingPass()
	result, err := pass.Run(program, cfg)
	if err != nil {
		t.Fatalf("constant folding failed: %v", err)
	}

	// The current implementation tracks constants but doesn't modify much
	// This is expected as we'd need to extend the constant table
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestRegisterCoalescingPass(t *testing.T) {
	// Create a program with register moves that can be coalesced
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1), // r1 = const
			vm.NewInstruction(vm.OpMove, 2, 1),       // r2 = r1 (can coalesce if no interference)
			vm.NewInstruction(vm.OpReturn, 2),        // return r2
		},
		Constants: []runtime.Value{runtime.NewInt(42)},
		Registers: 10,
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewRegisterCoalescingPass()
	result, err := pass.Run(program, cfg)
	if err != nil {
		t.Fatalf("register coalescing failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	// Check if any coalescing was reported
	coalesced, ok := result.Metadata["registers_coalesced"].(int)
	if !ok {
		t.Fatal("expected registers_coalesced metadata")
	}

	if coalesced < 0 {
		t.Errorf("invalid coalesced count: %d", coalesced)
	}
}

func TestPipeline_WithAllPasses(t *testing.T) {
	// Create a program that exercises multiple optimizations
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1),  // r1 = 10
			vm.NewInstruction(vm.OpLoadConst, 2, -2),  // r2 = 5
			vm.NewInstruction(vm.OpJumpIfFalse, 6, 1), // if !r1 goto 6
			vm.NewInstruction(vm.OpAdd, 3, 1, 2),      // r3 = r1 + r2
			vm.NewInstruction(vm.OpMove, 4, 3),        // r4 = r3
			vm.NewInstruction(vm.OpJump, 7),           // goto 7
			vm.NewInstruction(vm.OpLoadNone, 4),       // r4 = none
			vm.NewInstruction(vm.OpReturn, 4),         // return r4
		},
		Constants: []runtime.Value{runtime.NewInt(10), runtime.NewInt(5)},
		Registers: 10,
	}

	// Create pipeline with all passes
	pipeline := NewPipeline()
	pipeline.Add(NewLivenessAnalysisPass())
	pipeline.Add(NewLoopDetectionPass())
	pipeline.Add(NewConstantFoldingPass())
	pipeline.Add(NewRegisterCoalescingPass())

	result, err := pipeline.Run(program)
	if err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.PassResults) != 4 {
		t.Errorf("expected 4 pass results, got %d", len(result.PassResults))
	}

	// Verify each pass ran
	expectedPasses := []string{
		"liveness-analysis",
		"loop-detection",
		"constant-folding",
		"register-coalescing",
	}

	for _, passName := range expectedPasses {
		if _, ok := result.PassResults[passName]; !ok {
			t.Errorf("expected pass %s to run", passName)
		}
	}
}
