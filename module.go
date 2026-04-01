package ferret

// Module represents a self-contained unit of functionality that can be registered with the engine.
type (
	Module interface {
		// Name returns the module name used for identification and diagnostics.
		Name() string
		// Register applies the module's registrations to the engine bootstrap context.
		Register(Bootstrap) error
	}

	// Bootstrap defines an interface for configuring the host and registering lifecycle hooks with the runtime engine.
	Bootstrap interface {
		// Host returns access to host-level registration surfaces.
		Host() HostContext
		// Hooks returns access to engine, plan, and session hook registrars.
		Hooks() HookRegistrar
	}

	bootstrap struct {
		host  *hostContext
		hooks *hookRegistry
	}
)

func newBootstrap(opts *options) (*bootstrap, error) {
	hostCtx, err := newHostContext(opts)
	if err != nil {
		return nil, err
	}

	return &bootstrap{
		host:  hostCtx,
		hooks: opts.hooks,
	}, nil
}

func (b *bootstrap) Host() HostContext {
	return b.host
}

func (b *bootstrap) Hooks() HookRegistrar {
	return b.hooks
}
