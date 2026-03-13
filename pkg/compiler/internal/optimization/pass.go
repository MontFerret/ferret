package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	// PassResult contains the result of running a pass
	PassResult struct {
		Metadata map[string]any
		Modified bool
	}

	PassContext struct {
		Program  *bytecode.Program
		CFG      *ControlFlowGraph
		Metadata map[string]any // Shared metadata between passes
	}

	Pass interface {
		Name() string
		Requires() []string
		Run(ctx *PassContext) (*PassResult, error)
	}
)
