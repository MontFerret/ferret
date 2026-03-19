package optimization

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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

func resolvePassDependencies(passes []Pass, passIndex int, ctx *PassContext, metadataEpoch map[string]int, epoch int, resolving map[string]bool) error {
	pass := passes[passIndex]

	for _, dep := range pass.Requires() {
		if metadataIsFresh(ctx.Metadata, metadataEpoch, dep, epoch) {
			continue
		}

		depIndex := findEarlierPassIndex(passes, passIndex, dep)
		if depIndex < 0 {
			return fmt.Errorf("%w: pass %q requires %q which has not been executed", ErrMissingDependency, pass.Name(), dep)
		}

		if resolving == nil {
			resolving = make(map[string]bool)
		}

		if resolving[dep] {
			return fmt.Errorf("%w: cyclic dependency involving %q", ErrMissingDependency, dep)
		}

		resolving[dep] = true
		if err := resolvePassDependencies(passes, depIndex, ctx, metadataEpoch, epoch, resolving); err != nil {
			delete(resolving, dep)
			return err
		}
		delete(resolving, dep)

		if metadataIsFresh(ctx.Metadata, metadataEpoch, dep, epoch) {
			continue
		}

		dependency := passes[depIndex]
		passResult, err := dependency.Run(ctx)
		if err != nil {
			return fmt.Errorf("%w -> %s (dependency for %s): %w", ErrPassFailed, dependency.Name(), pass.Name(), err)
		}

		if passResult == nil {
			passResult = &PassResult{}
		}

		if passResult.Modified {
			return fmt.Errorf("%w: pass %q refreshed for %q", ErrDependencyRefreshModified, dependency.Name(), pass.Name())
		}

		ctx.Metadata[dependency.Name()] = passResult.Metadata
		metadataEpoch[dependency.Name()] = epoch
	}

	return nil
}

func metadataIsFresh(metadata map[string]any, metadataEpoch map[string]int, name string, epoch int) bool {
	if _, ok := metadata[name]; !ok {
		return false
	}

	return metadataEpoch[name] == epoch
}

func findEarlierPassIndex(passes []Pass, passIndex int, name string) int {
	for i := passIndex - 1; i >= 0; i-- {
		if passes[i].Name() == name {
			return i
		}
	}

	return -1
}

func Run(program *bytecode.Program, level Level) error {
	if level <= LevelNone {
		return nil
	}

	p := NewPipeline()
	p.Add(NewConstantPropagationPass())
	p.Add(NewLivenessAnalysisPass())
	p.Add(NewRegisterCoalescingPass())
	p.Add(NewPeepholePass())

	_, err := p.Run(program)

	return err
}
