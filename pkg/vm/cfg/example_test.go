package cfg_test

import (
	"fmt"
	"log"

	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/pkg/vm/cfg"
)

// Example demonstrates how to build and analyze a control flow graph
func Example() {
	// Create a simple program with a conditional branch
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadBool, 0, 1),    // Load condition
			vm.NewInstruction(vm.OpJumpIfFalse, 4, 0), // Conditional jump
			vm.NewInstruction(vm.OpLoadConst, 1, 0),   // Then branch
			vm.NewInstruction(vm.OpJump, 5),           // Jump to merge
			vm.NewInstruction(vm.OpLoadConst, 2, 0),   // Else branch
			vm.NewInstruction(vm.OpReturn, 0),         // Return
		},
	}

	// Build the control flow graph
	builder := cfg.NewBuilder(program)
	graph, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build CFG: %v", err)
	}

	// Print the CFG
	fmt.Println(graph.String())

	// Analyze the CFG
	analyzer := cfg.NewAnalyzer(graph)

	// Find unreachable blocks (dead code)
	unreachable := analyzer.FindUnreachableBlocks()
	fmt.Printf("Unreachable blocks: %d\n", len(unreachable))

	// Find loops (back edges)
	backEdges := analyzer.FindBackEdges()
	fmt.Printf("Loops (back edges): %d\n", len(backEdges))

	// Calculate dominators
	dominators := analyzer.CalculateDominators()
	fmt.Printf("Dominators calculated for %d blocks\n", len(dominators))

	// Generate DOT format for visualization
	dot := graph.ToDOT()
	fmt.Println("DOT format:")
	fmt.Println(dot)
}

// Example_loop demonstrates CFG analysis for a loop
func Example_loop() {
	// Create a program with a simple loop
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 0, 0),   // Loop header
			vm.NewInstruction(vm.OpJumpIfFalse, 4, 0), // Exit condition
			vm.NewInstruction(vm.OpLoadConst, 1, 0),   // Loop body
			vm.NewInstruction(vm.OpJump, 0),           // Back edge
			vm.NewInstruction(vm.OpReturn, 0),         // Exit
		},
	}

	// Build CFG
	builder := cfg.NewBuilder(program)
	graph, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build CFG: %v", err)
	}

	// Analyze for loops
	analyzer := cfg.NewAnalyzer(graph)
	backEdges := analyzer.FindBackEdges()

	fmt.Printf("Found %d loop(s)\n", len(backEdges))
	for i, edge := range backEdges {
		from, to := edge[0], edge[1]
		fmt.Printf("Loop %d: Block %d -> Block %d\n", i+1, from.ID, to.ID)
	}
}
