package exec

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	planOptions struct {
		env *vm.Environment
	}

	PlanOption func(opts *planOptions)
)

func (p *planOptions) Environment() *vm.Environment {
	return p.env
}

func newPlanOptions(setters []PlanOption) *planOptions {
	opts := &planOptions{
		env: vm.NewDefaultEnvironment(),
	}

	for _, setter := range setters {
		setter(opts)
	}

	return opts
}

func WithPlanEnvironment(env *vm.Environment) PlanOption {
	return func(o *planOptions) {
		o.env = vm.MergeEnvironments(o.env, env)
	}
}
