// Package sdk is the supported authoring layer for Ferret modules and host values.
//
// Module authors can construct callback-backed modules with NewModule, register
// declarative, arity-overloaded function definitions with RegisterFunctions,
// and adapt typed runtime functions with Bind0 through Bind4. Fixed arities take
// precedence over variadic fallbacks. Host-function names are case-insensitive
// in FQL and have one canonical lowercase terminal name; namespace casing is
// preserved. The binders intentionally operate
// on runtime.Value types; the SDK does not invoke arbitrary native Go functions
// through reflection.
//
// Encode, Decode, DecodeValue, DecodeArg, and Codec provide context-aware
// conversion at host boundaries. Decode options can require a root runtime
// type, restrict tagged root fields, reject unknown fields, and distinguish
// explicit None from omitted configuration. HostValue represents opaque
// wrapper identity with identity-only equality, while IterableValue,
// IteratorValue, SliceView, and MapView opt in to only the additional runtime
// capabilities they implement.
//
// The sdktest subpackage provides an Engine-backed black-box test harness for
// executing module functions through FQL.
package sdk
