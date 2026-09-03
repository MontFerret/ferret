package optimization

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	// Pipeline manages a sequence of passes to be executed on a program
	Pipeline struct {
		passes []Pass
	}

	// PipelineResult contains the results of running a pipeline
	PipelineResult struct {
		Modified bool
	}
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
func (p *Pipeline) Run(program *bytecode.Program) (*PipelineResult, error) {
	result := &PipelineResult{
		Modified: false,
	}

	// Build CFG once for all passes
	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build CFG: %w", err)
	}

	ctx := &PassContext{
		Program:  program,
		CFG:      cfg,
		Metadata: make(map[string]any),
	}
	metadataEpoch := make(map[string]int)
	epoch := 0
	rebuildRequired := false

	// Run each pass
	for i, pass := range p.passes {
		// If the previous pass modified the program, rebuild the CFG before
		// resolving dependencies so refreshed analyses observe the current state.
		if rebuildRequired {
			rebuildRequired = false

			// Rebuild CFG if program was modified
			cfg, err = builder.Build()

			if err != nil {
				// Report which pass caused the failure
				// by looking at the previous pass in the pipeline
				name := p.passes[i-1].Name()

				return nil, fmt.Errorf("%w after %s: %w", ErrCFGBuildFailed, name, err)
			}

			ctx.CFG = cfg
		}

		if err := resolvePassDependencies(p.passes, i, ctx, metadataEpoch, epoch, nil); err != nil {
			return nil, err
		}

		passResult, err := pass.Run(ctx)
		if err != nil {
			return nil, fmt.Errorf("%w -> %s: %w", ErrPassFailed, pass.Name(), err)
		}

		if passResult == nil {
			passResult = &PassResult{}
		}

		// Store pass metadata in context for future passes to use
		ctx.Metadata[pass.Name()] = passResult.Metadata
		if passResult.Modified {
			result.Modified = true
			epoch++
			rebuildRequired = true
		}

		metadataEpoch[pass.Name()] = epoch
	}

	return result, nil
}
