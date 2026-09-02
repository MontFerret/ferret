package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func Run(program *bytecode.Program, level Level) error {
	if level <= None {
		return nil
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
