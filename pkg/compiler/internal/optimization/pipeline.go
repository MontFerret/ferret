package optimization

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	Level int

	// Pipeline manages a sequence of passes to be executed on a program
	Pipeline struct {
		passes []Pass
	}

	// PipelineResult contains the results of running a pipeline
	PipelineResult struct {
		Modified bool
	}
)

const (
	LevelNone Level = iota
	LevelBasic
	LevelFull
	LevelAggressive
)

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
		Modified: false,
	}

	// Build CFG once for all passes
	builder := NewBuilder(program)
	cfg, err := builder.Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build CFG: %w", err)
	}

	var modified bool
	ctx := &PassContext{
		Program:  program,
		CFG:      cfg,
		Metadata: make(map[string]any),
	}

	// Run each pass
	for i, pass := range p.passes {
		// If previous pass modified the program, rebuild CFG
		// We do this before running the pass to ensure it has the latest CFG
		// And not after, to avoid unnecessary rebuilds
		if modified {
			result.Modified = true

			// Rebuild CFG if program was modified
			cfg, err = builder.Build()

			if err != nil {
				// Report which pass caused the failure
				// by looking at the previous pass in the pipeline
				name := p.passes[i-1].Name()
				return nil, fmt.Errorf("%w after %s: %w", ErrCFGBuildFailed, name, err)
			}
		}

		passResult, err := pass.Run(ctx)

		if err != nil {
			return nil, fmt.Errorf("%w -> %s: %w", ErrPassFailed, pass.Name(), err)
		}

		modified = passResult.Modified
		// Store pass metadata in context for future passes to use
		ctx.Metadata[pass.Name()] = passResult.Metadata
	}

	return result, nil
}

func Run(program *vm.Program, level Level) error {
	if level <= LevelNone {
		return nil
	}

	p := NewPipeline()
	p.Add(NewLivenessAnalysisPass())
	p.Add(NewRegisterCoalescingPass())

	//	//if level == "O2" {
	//	//	pm.Add(&ConstantPropagationPass{})
	//	//	pm.Add(&DCEPass{})
	//	//	pm.Add(&PeepholePass{})
	//	//	pm.Add(&RegisterCoalescingPass{})
	//	//}

	_, err := p.Run(program)

	return err
}
