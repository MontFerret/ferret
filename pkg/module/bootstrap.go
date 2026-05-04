package module

// Bootstrap exposes the mutable engine bootstrap state passed to Module.Register.
type Bootstrap interface {
	// Host returns the host-scoped registries and services that a module can
	// populate during engine construction.
	Host() HostContext
	// Hooks returns the lifecycle hook registrars for the engine, plan, and
	// session stages created from this engine.
	Hooks() HookRegistrar
}
