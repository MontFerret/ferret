package optimization

import "github.com/MontFerret/ferret/pkg/vm"

type (
	// PassResult contains the result of running a pass
	PassResult struct {
		Modified bool                   // Whether the program was modified
		Metadata map[string]interface{} // Pass-specific metadata
	}

	Pass interface {
		Name() string
		Requires() []string
		Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error)
	}
)
