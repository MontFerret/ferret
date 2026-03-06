package ferret

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session represents a single execution of a compiled Ferret program.
// It holds the state of the execution, including the virtual machine, environment, and encoding registry.
// A Session is created from a Plan and can be run to obtain results.
// It is not thread-safe and should be used for a single execution.
// After running, it can be closed to release any resources if necessary, or it can be reused for multiple runs if the environment and registry are not modified.
type Session struct {
	vm       *vm.VM
	env      *vm.Environment
	encoding *encoding.Registry
	hooks    sessionHooks
}

func (s *Session) Run(c context.Context) (Result, error) {
	ctx, err := s.hooks.runBeforeRunHooks(c)
	if err != nil {
		return nil, fmt.Errorf("before run hooks: %w", err)
	}

	ctx = encoding.WithRegistry(ctx, s.encoding)
	out, err := s.vm.Run(ctx, s.env)

	if hookErr := s.hooks.runAfterRunHooks(ctx, err); hookErr != nil {
		return nil, errors.Join(err, fmt.Errorf("after run hooks: %w", hookErr))
	}

	if err != nil {
		return nil, err
	}

	return newResult(out), nil
}

func (s *Session) Close() error {
	if err := s.hooks.runCloseHooks(); err != nil {
		return fmt.Errorf("close hooks: %w", err)
	}

	return nil
}
