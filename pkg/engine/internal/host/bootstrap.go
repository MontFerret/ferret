package host

import (
	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	"github.com/MontFerret/ferret/v2/pkg/module"
)

// Bootstrap exposes mutable host and hook registration during engine construction.
// The engine's bootstrap component owns module iteration, hook snapshotting,
// initialization, and rollback.
type Bootstrap struct {
	host  *hostContext
	hooks *Hooks
}

// NewBootstrap constructs services and records their ownership in resources.
// The caller must close that manager if construction or later bootstrap work fails.
func NewBootstrap(config Config, hooks *Hooks, resources *resource.Manager) (*Bootstrap, error) {
	hostCtx, err := newHostContext(config, resources)
	if err != nil {
		return nil, err
	}

	return &Bootstrap{host: hostCtx, hooks: hooks}, nil
}

// Host exposes the registries and services available to module registration.
func (b *Bootstrap) Host() module.HostContext {
	return b.host
}

// Hooks exposes the lifecycle registrars available to module registration.
func (b *Bootstrap) Hooks() module.HookRegistrar {
	return b.hooks
}

// Build finalizes the function library and snapshots parameters and codecs.
// It does not run initialization hooks or transfer resource ownership.
func (b *Bootstrap) Build() (*Host, error) {
	return b.host.build()
}
