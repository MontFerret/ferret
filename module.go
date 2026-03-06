package ferret

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (

	// ModuleRegistry provides an interface for modules to register their functions, encoding, and context decorators.
	ModuleRegistry struct {
		ns             runtime.Namespace
		encoding       *encoding.Registry
		onEngineInit   []EngineInitHook
		onEngineClose  []EngineCloseHook
		onSessionInit  []SessionInitHook
		onSessionClose []SessionCloseHook
	}

	// Module represents a self-contained unit of functionality that can be registered with the engine.
	Module interface {
		Name() string
		Register(registry *ModuleRegistry) error
	}

	EngineInitHook interface {
		OnEngineInit(ctx context.Context) error
	}

	EngineCloseHook interface {
		OnEngineClose(ctx context.Context) error
	}

	SessionInitHook interface {
		OnSessionInit(ctx context.Context) (context.Context, error)
	}

	SessionCloseHook interface {
		OnSessionClose(ctx context.Context) error
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

// OnEngineInit registers a hook that will be called when the engine is initialized.
func (mr *ModuleRegistry) OnEngineInit(hook EngineInitHook) {
	mr.onEngineInit = append(mr.onEngineInit, hook)
}

// OnEngineClose registers a hook that will be called when the engine is closed.
func (mr *ModuleRegistry) OnEngineClose(hook EngineCloseHook) {
	mr.onEngineClose = append(mr.onEngineClose, hook)
}

// OnSessionInit registers a hook that will be called when a session is initialized.
func (mr *ModuleRegistry) OnSessionInit(hook SessionInitHook) {
	mr.onSessionInit = append(mr.onSessionInit, hook)
}

// OnSessionClose registers a hook that will be called when a session is closed.
func (mr *ModuleRegistry) OnSessionClose(hook SessionCloseHook) {
	mr.onSessionClose = append(mr.onSessionClose, hook)
}
