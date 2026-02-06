package optimization

import "github.com/MontFerret/ferret/pkg/vm"

type (
	// PassResult contains the result of running a pass
	PassResult struct {
		Modified bool           // Whether the program was modified
		Metadata map[string]any // Pass-specific metadata
	}

	PassContext struct {
		Program  *vm.Program
		CFG      *ControlFlowGraph
		Metadata map[string]any // Shared metadata between passes
	}

	Pass interface {
		Name() string
		Requires() []string
		Run(ctx *PassContext) (*PassResult, error)
	}
)
