package module

// Module defines an engine extension that participates in bootstrap.
//
// Modules are passed to the engine during construction and can register host
// services, codecs, and lifecycle hooks before the engine is initialized.
type Module interface {
	// Name returns the stable module identifier used in diagnostics and error
	// reporting.
	Name() string
	// Register mutates the engine bootstrap state for this module instance.
	// Returning an error aborts engine construction.
	Register(Bootstrap) error
}
