package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Plan struct {
	prog     *bytecode.Program
	env      *vm.Environment
	registry *encoding.Registry
}

func newPlan(prog *bytecode.Program, env *vm.Environment, registry *encoding.Registry) *Plan {
	if registry == nil {
		registry = encoding.NewRegistry()
	}

	return &Plan{
		prog:     prog,
		env:      env,
		registry: registry,
	}
}

func (p *Plan) NewSession(setters ...SessionOption) *Session {
	env := vm.NewEnvironment(setters)

	return newSession(vm.New(p.prog), vm.MergeEnvironments(p.env, env), p.registry)
}
