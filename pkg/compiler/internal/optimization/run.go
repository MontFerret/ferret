package optimization

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Run applies the optimizer pipeline selected by level. Basic runs constant
// propagation, liveness analysis, and peephole optimization; Full additionally
// runs register coalescing.
func Run(program *bytecode.Program, level Level) error {
	switch level {
	case None:
		return nil
	case Basic, Full:
	default:
		return fmt.Errorf("unsupported optimization level %d", level)
	}

	p := NewPipeline()
	p.Add(NewConstantPropagationPass())
	p.Add(NewLivenessAnalysisPass())

	if level == Full {
		p.Add(NewRegisterCoalescingPass())
	}

	p.Add(NewPeepholePass())

	_, err := p.Run(program)

	return err
}
