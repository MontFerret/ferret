package session

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/host"
)

// NewDebugSession acquires host resources and a dedicated retained execution.
// The debugger owns the VM; its services own only host resources and the permit.
// The native caller must close a successful session if publication is canceled.
func NewDebugSession(ctx context.Context, config Config, h *host.Host, hooks *host.SessionHooks, limiter *Limiter, program *bytecode.Program, planClosed <-chan struct{}) (_ *debugger.Session, err error) {
	a := acquisition{limiter: limiter}
	defer a.rollback(&err)

	if err := a.prepare(ctx, config, h, planClosed); err != nil {
		return nil, err
	}

	return a.newDebugSession(config, h, hooks, program)
}
