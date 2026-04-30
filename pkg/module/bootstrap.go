package module

// Bootstrap defines an interface for configuring the host and registering lifecycle hooks with the runtime engine.
type Bootstrap interface {
	// Host returns access to host-level registration surfaces.
	Host() HostContext
	// Hooks returns access to engine, plan, and session hook registrars.
	Hooks() HookRegistrar
}
