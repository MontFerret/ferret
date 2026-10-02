// Package uapi exposes Ferret Native through the Universal Ferret API.
//
// It is Ferret's official adapter for github.com/MontFerret/api, translating
// portable source, options, and results into Native engine operations while
// keeping Native and Universal APIs independently evolvable.
//
// New creates and owns a Native engine. Wrap borrows an existing engine, and
// closing the adapter leaves it usable. Callers settle work and close sessions
// and plans before their owning runtime or the borrowed Native engine.
//
// Plan.Params and Runtime.Version honor non-nil caller contexts even though
// metadata retrieval is local. Params preserves Native's detached snapshot and
// post-close access. Version returns the Ferret Core implementation version
// supplied to New or Wrap unchanged, including after Close. Callers supply the
// Core version independently of the host. Native Plan.Params remains context-free.
package uapi
