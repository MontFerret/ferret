// Package universal exposes a Native Ferret engine through the Universal API.
// New creates and owns a Native engine; Runtime.Close delegates its cleanup to
// Native. Wrap borrows an existing engine, and its Close is a no-op that leaves
// the adapter usable. Callers settle work and close directly created sessions
// and plans before their owning runtime or the borrowed Native engine.
//
// Universal session setters queue Native options; Native validates and converts
// them when applied. Plan options are translated explicitly. Native
// configuration, execution, resource ownership, and debugging remain with their
// owning packages. Source coordinates, encoded output, and debugger values already
// use portable types; Native diagnostics are projected while retaining their causes.
//
// Portable option callbacks run before delegation. Native receives caller
// contexts unchanged and owns context validation and cancellation. Plan and
// session close calls delegate to Native. Ordinary Run calls must settle before
// session cleanup; Native debug Close terminates and settles active commands.
package universal
