package ferret

// Module represents a self-contained unit of functionality that can be registered with the engine.
type (
	Module interface {
		Name() string
		Register(Bootstrap) error
	}

	// Bootstrap defines an interface for configuring the host and registering lifecycle hooks with the runtime engine.
	Bootstrap interface {
		Host() HostContext
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
