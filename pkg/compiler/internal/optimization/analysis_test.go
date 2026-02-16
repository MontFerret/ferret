package optimization

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestAnalyzer_FindReachableBlocks(t *testing.T) {
	// Create a program with unreachable code
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),
			bytecode.NewInstruction(bytecode.OpJump, 4),         // Jump over unreachable code
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, 0), // Unreachable
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, 0), // Unreachable
			bytecode.NewInstruction(bytecode.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyzer := NewAnalyzer(cfg)
	reachable := analyzer.FindReachableBlocks()

	// Should have entry block, return block, and exit block reachable
	// But not the unreachable block
	if len(reachable) != 3 {
		t.Errorf("expected 3 reachable blocks, got %d", len(reachable))
	}
}

func TestAnalyzer_FindUnreachableBlocks(t *testing.T) {
	// Create a program with unreachable code
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),
			bytecode.NewInstruction(bytecode.OpJump, 4),         // Jump over unreachable code
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, 0), // Unreachable
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, 0), // Unreachable
			bytecode.NewInstruction(bytecode.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyzer := NewAnalyzer(cfg)
	unreachable := analyzer.FindUnreachableBlocks()

	// Should have one unreachable block (the one we jump over)
	if len(unreachable) != 1 {
		t.Errorf("expected 1 unreachable block, got %d", len(unreachable))
	}

	if len(unreachable) > 0 && unreachable[0].Start != 2 {
		t.Errorf("expected unreachable block to start at index 2, got %d", unreachable[0].Start)
	}
}

func TestAnalyzer_FindBackEdges(t *testing.T) {
	// Create a program with a loop
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),   // Block 0: loop header
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, 4, 0), // Block 0: exit condition
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, 0),   // Block 1: loop body
			bytecode.NewInstruction(bytecode.OpJump, 0),           // Block 1: back to loop header (back edge)
			bytecode.NewInstruction(bytecode.OpReturn, 0),         // Block 2: exit
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyzer := NewAnalyzer(cfg)
	backEdges := analyzer.FindBackEdges()

	// Should have one back edge (from loop body back to loop header)
	if len(backEdges) != 1 {
		t.Errorf("expected 1 back edge, got %d", len(backEdges))
	}

	if len(backEdges) > 0 {
		from, to := backEdges[0][0], backEdges[0][1]
		if from.Start != 2 || to.Start != 0 {
			t.Errorf("expected back edge from block at 2 to block at 0, got %d to %d", from.Start, to.Start)
		}
	}
}

func TestAnalyzer_FindBackEdges_NoLoop(t *testing.T) {
	// Create a simple linear program
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),
			bytecode.NewInstruction(bytecode.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyzer := NewAnalyzer(cfg)
	backEdges := analyzer.FindBackEdges()

	// Should have no back edges
	if len(backEdges) != 0 {
		t.Errorf("expected 0 back edges, got %d", len(backEdges))
	}
}

func TestAnalyzer_CalculateDominators(t *testing.T) {
	// Create a simple if-else program
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadBool, 0, 1),    // Block 0: condition
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, 4, 0), // Block 0
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, 0),   // Block 1: then
			bytecode.NewInstruction(bytecode.OpJump, 5),           // Block 1
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, 0),   // Block 2: else
			bytecode.NewInstruction(bytecode.OpReturn, 0),         // Block 3: merge
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	analyzer := NewAnalyzer(cfg)
	dominators := analyzer.CalculateDominators()

	// Entry block should have no immediate dominator
	if _, ok := dominators[cfg.Entry.ID]; ok {
		t.Errorf("entry block should not have an immediate dominator")
	}

	// All other blocks should be dominated by entry
	for _, block := range cfg.Blocks {
		if block == cfg.Entry || block == cfg.Exit {
			continue
		}
		if dom, ok := dominators[block.ID]; !ok || dom == nil {
			t.Errorf("block %d should have an immediate dominator", block.ID)
		}
	}
}

func TestCFG_ToDOT(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),
			bytecode.NewInstruction(bytecode.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dot := cfg.ToDOT()

	// Check that DOT output contains expected elements
	if !strings.Contains(dot, "digraph CFG") {
		t.Errorf("DOT output should contain 'digraph CFG'")
	}

	if !strings.Contains(dot, "block0") {
		t.Errorf("DOT output should contain entry block")
	}

	if !strings.Contains(dot, "->") {
		t.Errorf("DOT output should contain edges")
	}
}

func TestCFG_String(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 0, 0),
			bytecode.NewInstruction(bytecode.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	str := cfg.String()

	// Check that string output contains expected elements
	if !strings.Contains(str, "Control Flow Graph") {
		t.Errorf("String output should contain 'Control Flow Graph'")
	}

	if !strings.Contains(str, "Block") {
		t.Errorf("String output should contain block information")
	}
}
