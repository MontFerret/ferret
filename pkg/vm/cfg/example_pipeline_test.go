package cfg_test

import (
	"fmt"
	"log"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/pkg/vm/cfg"
)

// Example_pipeline demonstrates using the optimization pipeline
func Example_pipeline() {
	// Create a program with optimization opportunities
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1), // r1 = 10
			vm.NewInstruction(vm.OpLoadConst, 2, -2), // r2 = 5
			vm.NewInstruction(vm.OpMove, 3, 1),       // r3 = r1 (can be coalesced)
			vm.NewInstruction(vm.OpAdd, 4, 3, 2),     // r4 = r3 + r2
			vm.NewInstruction(vm.OpReturn, 4),
		},
		Constants: []runtime.Value{
			runtime.NewInt(10),
			runtime.NewInt(5),
		},
		Registers: 10,
	}

	// Create optimization pipeline
	pipeline := cfg.NewPipeline()

	// Add analysis passes
	pipeline.Add(cfg.NewLivenessAnalysisPass())
	pipeline.Add(cfg.NewLoopDetectionPass())

	// Add transformation passes
	pipeline.Add(cfg.NewConstantFoldingPass())
	pipeline.Add(cfg.NewRegisterCoalescingPass())

	// Run the pipeline
	result, err := pipeline.Run(program)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	// Display results
	fmt.Printf("Program modified: %v\n", result.Modified)
	fmt.Printf("Passes run: %d\n", len(result.PassResults))

	// Show results from each pass
	for passName, passResult := range result.PassResults {
		fmt.Printf("\nPass: %s\n", passName)
		fmt.Printf("  Modified: %v\n", passResult.Modified)

		// Display pass-specific metadata
		if passName == "loop-detection" {
			if loops, ok := passResult.Metadata["loops"].([]*cfg.Loop); ok {
				fmt.Printf("  Loops found: %d\n", len(loops))
			}
		}

		if passName == "register-coalescing" {
			if count, ok := passResult.Metadata["registers_coalesced"].(int); ok {
				fmt.Printf("  Registers coalesced: %d\n", count)
			}
		}
	}
}

// Example_livenessAnalysis demonstrates liveness analysis
func Example_livenessAnalysis() {
	program := &vm.Program{
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, 1, -1), // r1 = 42
			vm.NewInstruction(vm.OpLoadConst, 2, -2), // r2 = 10
			vm.NewInstruction(vm.OpAdd, 3, 1, 2),     // r3 = r1 + r2
			vm.NewInstruction(vm.OpReturn, 3),        // return r3
		},
		Constants: []runtime.Value{
			runtime.NewInt(42),
			runtime.NewInt(10),
		},
		Registers: 10,
	}

	// Build CFG
	builder := cfg.NewBuilder(program)
	graph, _ := builder.Build()

	// Run liveness analysis
	pass := cfg.NewLivenessAnalysisPass()
	result, _ := pass.Run(program, graph)

	fmt.Println("Liveness analysis completed")
	fmt.Printf("Modified program: %v\n", result.Modified)
}
