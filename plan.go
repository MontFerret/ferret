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

	return &Session{
		// TODO: create a VM pool and get a VM from it instead of creating a new one for each session
		vm:       vm.New(p.prog),
		env:      env,
		encoding: p.host.encoding,
		hooks:    p.sessionHooks,
	}, nil
}

func (p *Plan) Close() error {
	if err := p.hooks.runCloseHooks(); err != nil {
		return fmt.Errorf("close hooks: %w", err)
	}

	return nil
}
