// Package uapi exposes Ferret Native through the Universal Ferret API.
//
// It is Ferret's official adapter for github.com/MontFerret/api, translating
// portable source, options, and results into Native engine operations while
// keeping Native and Universal APIs independently evolvable.
//
// New creates and owns a Native engine. Wrap borrows an existing engine, and
// closing the adapter leaves it usable. Callers settle work and close sessions
// and plans before their owning runtime or the borrowed Native engine.
package uapi
