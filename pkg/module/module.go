package module

// Module represents a self-contained unit of functionality that can be registered with the engine.
type Module interface {
	// Name returns the module name used for identification and diagnostics.
	Name() string
	// Register applies the module's registrations to the engine bootstrap context.
	Register(Bootstrap) error
}
