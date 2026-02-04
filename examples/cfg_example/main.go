package main

import (
	"fmt"
	"log"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/vm/cfg"
)

func main() {
	// Compile a simple FQL query with control flow
	query := `
		LET items = [1, 2, 3, 4, 5]
		
		FOR item IN items
			FILTER item > 2
			RETURN item * 2
	`

	// Create source and compile to bytecode
	src := file.NewAnonymousSource(query)
	c := compiler.New()
	program, err := c.Compile(src)
	if err != nil {
		log.Fatalf("Failed to compile query: %v", err)
	}

	// Build the control flow graph
	builder := cfg.NewBuilder(program)
	graph, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build CFG: %v", err)
	}

	// Display basic information
	fmt.Println("=== Control Flow Graph ===")
	fmt.Printf("Total blocks: %d\n", len(graph.Blocks))
	fmt.Printf("Entry block: %d\n", graph.Entry.ID)
	fmt.Printf("Exit block: %d\n", graph.Exit.ID)
	fmt.Println()

	// Print detailed CFG structure
	fmt.Println(graph.String())

	// Perform analysis
	analyzer := cfg.NewAnalyzer(graph)

	// Check for unreachable code
	unreachable := analyzer.FindUnreachableBlocks()
	if len(unreachable) > 0 {
		fmt.Printf("Warning: Found %d unreachable block(s)\n", len(unreachable))
		for _, block := range unreachable {
			fmt.Printf("  - Block %d at instructions [%d:%d]\n", block.ID, block.Start, block.End)
		}
	} else {
		fmt.Println("✓ No unreachable code detected")
	}
	fmt.Println()

	// Check for loops
	backEdges := analyzer.FindBackEdges()
	if len(backEdges) > 0 {
		fmt.Printf("Found %d loop(s):\n", len(backEdges))
		for i, edge := range backEdges {
			from, to := edge[0], edge[1]
			fmt.Printf("  - Loop %d: Block %d -> Block %d\n", i+1, from.ID, to.ID)
		}
	} else {
		fmt.Println("✓ No loops detected")
	}
	fmt.Println()

	// Run optimization pipeline
	fmt.Println("=== Optimization Pipeline ===")
	pipeline := cfg.NewPipeline()
	pipeline.Add(cfg.NewLivenessAnalysisPass())
	pipeline.Add(cfg.NewLoopDetectionPass())
	pipeline.Add(cfg.NewConstantFoldingPass())
	pipeline.Add(cfg.NewRegisterCoalescingPass())

	pipelineResult, err := pipeline.Run(program)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	fmt.Printf("Program modified: %v\n", pipelineResult.Modified)
	fmt.Printf("Passes executed: %d\n", len(pipelineResult.PassResults))
	
	for passName, passResult := range pipelineResult.PassResults {
		status := "✓"
		if passResult.Modified {
			status = "✓ (modified)"
		}
		fmt.Printf("  %s %s\n", status, passName)
	}
	fmt.Println()

	// Generate DOT format for visualization
	fmt.Println("=== Graphviz DOT Format ===")
	fmt.Println(graph.ToDOT())
	fmt.Println("\nTo visualize, save the DOT output to a file and run:")
	fmt.Println("  dot -Tpng cfg.dot -o cfg.png")
}
