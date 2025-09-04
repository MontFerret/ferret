package engine

import (
	"context"

	"github.com/MontFerret/ferret/pkg/vm"
)

type Session struct {
	vm  *vm.VM
	env *vm.Environment
}

func newSession(vmi *vm.VM, env *vm.Environment) *Session {
	return &Session{
		vm:  vmi,
		env: env,
	}
}

func (s *Session) Run(ctx context.Context) (Result, error) {
	out, err := s.vm.Run(ctx, s.env)

	if err != nil {
		return nil, err
	}

	return newResult(out), nil
}

func (s *Session) Close() error {
	return nil
}
