package ferret

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// ContextDecorator is a function that takes a context and returns a new context with additional values or modifications.
	// It can be used to extend the context with custom data or behavior before executing a module's functions.
	ContextDecorator func(ctx context.Context) (context.Context, error)

	// ModuleRegistry provides an interface for modules to register their functions, encoding, and context decorators.
	ModuleRegistry struct {
		ns         runtime.Namespace
		encoding   *encoding.Registry
		decorators []ContextDecorator
		closers    []io.Closer
	}

	// Module represents a self-contained unit of functionality that can be registered with the engine.
	Module interface {
		Name() string
		Register(registry *ModuleRegistry) error
	}
)

// Functions returns the runtime.Namespace instance associated with the ModuleRegistry.
func (mr *ModuleRegistry) Functions() runtime.Namespace {
	return mr.ns
}

// Encoding retrieves the encoding.Registry instance associated with the ModuleRegistry, which stores codecs by content type.
func (mr *ModuleRegistry) Encoding() *encoding.Registry {
	return mr.encoding
}

// WithContext adds a ContextDecorator to the list of decorators in the ModuleRegistry.
func (mr *ModuleRegistry) WithContext(extender ContextDecorator) {
	mr.decorators = append(mr.decorators, extender)
}

// WithCloser registers an io.Closer instance to be invoked during the module's cleanup process.
func (mr *ModuleRegistry) WithCloser(closer io.Closer) {
	mr.closers = append(mr.closers, closer)
}
