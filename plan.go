package ferret

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Plan struct {
	prog         *bytecode.Program
	host         *host
	hooks        planHooks
	sessionHooks sessionHooks
}

func (p *Plan) NewSession(setters ...SessionOption) (*Session, error) {
	env, err := vm.ExtendEnvironment(&vm.Environment{
		Functions: p.host.functions,
		Params:    p.host.params,
		Logging:   p.host.logging,
	}, setters)

	if err != nil {
		return nil, err
	}

	instance, err := vm.New(p.prog)
	if err != nil {
		return nil, err
	}

	return &Session{
		// TODO: create a VM pool and get a VM from it instead of creating a new one for each session
		vm:       instance,
		env:      env,
		encoding: p.host.encoding,
		hooks:    p.sessionHooks,
	}, nil
}

func (p *Plan) Close() error {
	// Plan close hooks follow the hook registry close semantics (LIFO with error aggregation).
	if err := p.hooks.runCloseHooks(); err != nil {
		return fmt.Errorf("close hooks: %w", err)
	}

	return nil
}
