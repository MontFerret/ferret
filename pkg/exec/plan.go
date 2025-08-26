package exec

import "github.com/MontFerret/ferret/pkg/vm"

type Plan struct {
	prog    *vm.Program
	options *planOptions
}

func newPlan(prog *vm.Program, setters []PlanOption) *Plan {
	return &Plan{
		prog:    prog,
		options: newPlanOptions(setters),
	}
}

func (p *Plan) NewSession(setters ...SessionOption) *Session {
	env := vm.NewEnvironment(setters)

	return newSession(vm.New(p.prog), vm.MergeEnvironments(p.options.env, env))
}
