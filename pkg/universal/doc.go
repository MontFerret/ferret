// Package universal exposes a Native Ferret engine through the Universal API.
// New borrows its engine. Callers close sessions, plans, and the adapter before
// closing the engine; closing the adapter never closes the engine or descendants.
//
// Universal options are translated into Native options at this boundary. Native
// configuration, execution, resource ownership, and debugging remain with their
// owning packages. Source coordinates, encoded output, and debugger values already
// use portable types; Native diagnostics are projected while retaining their causes.
//
// Runtime.Close waits for admitted calls without canceling them. Plan.Close
// cancels pending session construction and waits before closing the Native plan.
// Published children remain caller-owned. Ordinary Run calls must be settled
// before closing their session; debug Close terminates and settles commands.
// Close must not be called synchronously from an option or hook inside an
// admitted operation on the same parent, or recursively from its close hook.
package universal
