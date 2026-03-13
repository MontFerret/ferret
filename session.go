package ferret

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// Session represents the execution of a compiled Ferret program.
// It holds the state of the execution, including the virtual machine, environment, and encoding registry.
// A Session is created from a Plan and can be run to obtain results.
//
// Session is not safe for concurrent use by multiple goroutines.
// It is typically used for a single logical execution. When a Session is created
// directly via Plan.NewSession, it may be reused for multiple sequential runs as
// long as the environment and encoding registry are not modified between runs.
// Helper APIs such as Engine.Run may take ownership of the Session and close it
// after a single execution, in which case the caller must not attempt to reuse it.
type Session struct {
	vm       *vm.VM
	env      *vm.Environment
	encoding *encoding.Registry
	hooks    sessionHooks
}

func (s *Session) Run(c context.Context) (*Result, error) {
	// Before-run hooks can replace the context used for the rest of execution.
	ctx, err := s.hooks.runBeforeRunHooks(c)
	if err != nil {
		return nil, fmt.Errorf("before run hooks: %w", err)
	}

	ctx = encoding.WithRegistry(ctx, s.encoding)
	out, err := s.vm.Run(ctx, s.env)

	// After-run hooks always run and receive the VM run error (if any).
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
