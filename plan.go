package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Plan struct {
	prog     *bytecode.Program
	env      *vm.Environment
	encoding *encoding.Registry
}

func newPlan(prog *bytecode.Program, env *vm.Environment, enc *encoding.Registry) *Plan {
	if enc == nil {
		enc = encoding.NewRegistry()
	}

	return &Plan{
		prog:     prog,
		env:      env,
		encoding: enc,
	}
}

func (p *Plan) NewSession(setters ...SessionOption) (*Session, error) {
	env, err := vm.NewEnvironment(setters)

	if err != nil {
		return nil, err
	}

	mergedEnv, err := vm.MergeEnvironments(p.env, env)

	if err != nil {
		return nil, err
	}

	return newSession(vm.New(p.prog), mergedEnv, p.encoding), nil
}
