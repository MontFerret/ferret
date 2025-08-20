package ferret

import (
	"context"

	"github.com/MontFerret/ferret/pkg/vm"
)

type Session struct {
	vm *vm.VM
}

func (s *Session) Run(ctx context.Context, opts []SessionOption) (Result, error) {
	out, err := s.vm.Run(ctx, opts)

	if err != nil {
		return nil, err
	}

	return newResult(out), nil
}

func (s *Session) Close() error {
	return nil
}
