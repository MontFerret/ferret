package cfg

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/vm"
)

// Pass represents a transformation or analysis pass over a program
type Pass interface {
	// Name returns the name of the pass
	Name() string

	// Run executes the pass on the given program
	Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error)
}

// PassResult contains the result of running a pass
type PassResult struct {
	Modified bool                   // Whether the program was modified
	Metadata map[string]interface{} // Pass-specific metadata
}

// Pipeline manages a sequence of passes to be executed on a program
type Pipeline struct {
	passes []Pass
}

// NewPipeline creates a new pass pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		passes: make([]Pass, 0),
	}
}

// Add adds a pass to the pipeline
func (p *Pipeline) Add(pass Pass) {
	p.passes = append(p.passes, pass)
}

// Run executes all passes in the pipeline
func (p *Pipeline) Run(program *vm.Program) (*PipelineResult, error) {
	result := &PipelineResult{
		PassResults: make(map[string]*PassResult),
		Modified:    false,
	}

	// Build CFG once for all passes
	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build CFG: %w", err)
	}

	// Run each pass
	for _, pass := range p.passes {
		passResult, err := pass.Run(program, cfg)
		if err != nil {
			return nil, fmt.Errorf("pass %s failed: %w", pass.Name(), err)
		}

		result.PassResults[pass.Name()] = passResult
		if passResult.Modified {
			result.Modified = true
			// Rebuild CFG if program was modified
			cfg, err = builder.Build()
			if err != nil {
				return nil, fmt.Errorf("failed to rebuild CFG after %s: %w", pass.Name(), err)
			}
		}
	}

	return result, nil
}

// PipelineResult contains the results of running a pipeline
type PipelineResult struct {
	PassResults map[string]*PassResult
	Modified    bool
}
