package ferret

// Module represents a self-contained unit of functionality that can be registered with the engine.
type (
	Module interface {
		Name() string
		Register(Bootstrap) error
	}

	// Bootstrap defines an interface for configuring the host and registering lifecycle hooks with the runtime engine.
	Bootstrap interface {
		Host() HostConfigurer
		Hooks() HookRegistrar
	}

	bootstrap struct {
		host  *hostBuilder
		hooks *hookRegistry
	}
)

func newBootstrap(opts *options) *bootstrap {
	return &bootstrap{
		host:  newHostBuilder(opts),
		hooks: opts.hooks,
	}
}

func (b *bootstrap) Host() HostConfigurer {
	return b.host
}

func (b *bootstrap) Hooks() HookRegistrar {
	return b.hooks
}
