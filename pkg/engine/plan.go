package engine

import "github.com/MontFerret/ferret/v2/pkg/vm"

type Plan struct {
	prog *vm.Program
	env  *vm.Environment
}

func newPlan(prog *vm.Program, env *vm.Environment) *Plan {
	return &Plan{
		prog: prog,
		env:  env,
	}
}

func (p *Plan) NewSession(setters ...SessionOption) *Session {
	env := vm.NewEnvironment(setters)

	return newSession(vm.New(p.prog), vm.MergeEnvironments(p.env, env))
}
