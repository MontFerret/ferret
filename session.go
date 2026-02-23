package ferret

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Session struct {
	vm       *vm.VM
	env      *vm.Environment
	registry *encoding.Registry
}

func newSession(vmi *vm.VM, env *vm.Environment, registry *encoding.Registry) *Session {
	if registry == nil {
		registry = encoding.NewRegistry()
	}

	return &Session{
		vm:       vmi,
		env:      env,
		registry: registry,
	}
}

func (s *Session) Run(ctx context.Context) (Result, error) {
	ctx = encoding.WithRegistry(ctx, s.registry)

	out, err := s.vm.Run(ctx, s.env)

	if err != nil {
		return nil, err
	}

	return newResult(out), nil
}

func (s *Session) Close() error {
	return nil
}
