package optimization

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestBuildCFG_EmptyProgram(t *testing.T) {
	program := &vm.Program{
		Bytecode: []vm.Instruction{},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Entry != nil {
		t.Errorf("expected nil entry for empty program")
	}

	if cfg.Exit == nil {
		t.Errorf("expected exit block for empty program")
	}

	if len(cfg.Blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(cfg.Blocks))
	}
}

func TestBuildCFG_SingleBlock(t *testing.T) {
	// Simple program with no control flow
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, 0),
			vm.NewInstruction(vm.OpReturn, 0),
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Entry == nil {
		t.Fatalf("expected entry block")
	}

	// Should have 2 blocks: main block + exit block
	if len(cfg.Blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(cfg.Blocks))
	}

	// Entry block should have 2 instructions
	if len(cfg.Entry.Instructions) != 2 {
		t.Errorf("expected 2 instructions in entry block, got %d", len(cfg.Entry.Instructions))
	}

	// Entry block should have exit as successor
	if len(cfg.Entry.Successors) != 1 {
		t.Errorf("expected 1 successor, got %d", len(cfg.Entry.Successors))
	}

	if cfg.Entry.Successors[0] != cfg.Exit {
		t.Errorf("expected entry successor to be exit block")
	}
}

func TestBuildCFG_UnconditionalJump(t *testing.T) {
	// Program with unconditional jump
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, 0), // Block 0
			vm.NewInstruction(vm.OpJump, 3),         // Block 0
			vm.NewInstruction(vm.OpLoadConst, 1, 0), // Block 1 (unreachable)
			vm.NewInstruction(vm.OpLoadConst, 2, 0), // Block 2 (jump target)
			vm.NewInstruction(vm.OpReturn, 0),       // Block 2
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 4 blocks: 0, 1 (unreachable), 2, exit
	if len(cfg.Blocks) != 4 {
		t.Errorf("expected 4 blocks, got %d", len(cfg.Blocks))
	}

	// Block 0 should jump to block 2
	block0 := cfg.Blocks[0]
	if len(block0.Successors) != 1 {
		t.Errorf("expected 1 successor for block 0, got %d", len(block0.Successors))
	}

	// Check that block 0 jumps to the block starting at index 3
	if block0.Successors[0].Start != 3 {
		t.Errorf("expected block 0 to jump to block starting at 3, got %d", block0.Successors[0].Start)
	}
}

func TestBuildCFG_ConditionalJump(t *testing.T) {
	// Program with conditional jump (if-else structure)
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadBool, 0, 1),    // Block 0: condition
			vm.NewInstruction(vm.OpJumpIfFalse, 4, 0), // Block 0: if false, jump to 4
			vm.NewInstruction(vm.OpLoadConst, 1, 0),   // Block 1: then branch
			vm.NewInstruction(vm.OpJump, 5),           // Block 1: jump to merge
			vm.NewInstruction(vm.OpLoadConst, 2, 0),   // Block 2: else branch
			vm.NewInstruction(vm.OpReturn, 0),         // Block 3: merge
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 5 blocks: 0, 1, 2, 3, exit
	if len(cfg.Blocks) != 5 {
		t.Errorf("expected 5 blocks, got %d", len(cfg.Blocks))
	}

	// Block 0 (condition) should have 2 successors
	block0 := cfg.Blocks[0]
	if len(block0.Successors) != 2 {
		t.Errorf("expected 2 successors for block 0, got %d", len(block0.Successors))
	}

	// Check successors are at indices 2 (fall-through) and 4 (jump target)
	successorStarts := make(map[int]bool)
	for _, succ := range block0.Successors {
		successorStarts[succ.Start] = true
	}

	if !successorStarts[2] || !successorStarts[4] {
		t.Errorf("expected block 0 successors at indices 2 and 4")
	}
}

func TestBuildCFG_Loop(t *testing.T) {
	// Program with a loop (back edge)
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, 0),   // Block 0: loop header
			vm.NewInstruction(vm.OpJumpIfFalse, 4, 0), // Block 0: exit condition
			vm.NewInstruction(vm.OpLoadConst, 1, 0),   // Block 1: loop body
			vm.NewInstruction(vm.OpJump, 0),           // Block 1: back to loop header
			vm.NewInstruction(vm.OpReturn, 0),         // Block 2: exit
		},
	}

	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have 4 blocks: 0 (header), 1 (body), 2 (exit), virtual exit
	if len(cfg.Blocks) != 4 {
		t.Errorf("expected 4 blocks, got %d", len(cfg.Blocks))
	}

	// Block 0 should have 2 successors (loop body and exit)
	block0 := cfg.Blocks[0]
	if len(block0.Successors) != 2 {
		t.Errorf("expected 2 successors for block 0, got %d", len(block0.Successors))
	}

	// Block 1 (loop body) should jump back to block 0
	block1 := cfg.Blocks[1]
	if len(block1.Successors) != 1 {
		t.Errorf("expected 1 successor for block 1, got %d", len(block1.Successors))
	}

	if block1.Successors[0] != block0 {
		t.Errorf("expected block 1 to jump back to block 0")
	}

	// Block 0 should have block 1 as a predecessor (loop back edge)
	if len(block0.Predecessors) != 1 {
		t.Errorf("expected 1 predecessor for block 0, got %d", len(block0.Predecessors))
	}

	if block0.Predecessors[0] != block1 {
		t.Errorf("expected block 0 predecessor to be block 1")
	}
}

func TestBasicBlock_AddSuccessor(t *testing.T) {
	block1 := NewBasicBlock(1, 0)
	block2 := NewBasicBlock(2, 5)

	block1.AddSuccessor(block2)

	if len(block1.Successors) != 1 {
		t.Errorf("expected 1 successor, got %d", len(block1.Successors))
	}

	if len(block2.Predecessors) != 1 {
		t.Errorf("expected 1 predecessor, got %d", len(block2.Predecessors))
	}

	// Adding the same successor again should not duplicate
	block1.AddSuccessor(block2)

	if len(block1.Successors) != 1 {
		t.Errorf("expected 1 successor after duplicate add, got %d", len(block1.Successors))
	}
}

func TestBasicBlock_IsTerminator(t *testing.T) {
	tests := []struct {
		name     string
		opcode   vm.Opcode
		expected bool
	}{
		{"Return", vm.OpReturn, true},
		{"Jump", vm.OpJump, true},
		{"JumpIfFalse", vm.OpJumpIfFalse, true},
		{"JumpIfTrue", vm.OpJumpIfTrue, true},
		{"IterNext", vm.OpIterNext, true},
		{"Add", vm.OpAdd, false},
		{"LoadConst", vm.OpLoadConst, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := NewBasicBlock(1, 0)
			block.AddInstruction(vm.NewInstruction(tt.opcode))

			if block.IsTerminator() != tt.expected {
				t.Errorf("expected IsTerminator() = %v for %s", tt.expected, tt.name)
			}
		})
	}
}
