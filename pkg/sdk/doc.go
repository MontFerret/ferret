// Package sdk is the supported authoring layer for Ferret modules and host values.
//
// Module authors can construct callback-backed modules with NewModule, register
// declarative function definitions with RegisterFunctions, and adapt typed
// runtime functions with Bind0 through Bind4. The binders intentionally operate
// on runtime.Value types; the SDK does not invoke arbitrary native Go functions
// through reflection.
//
// Encode, Decode, DecodeValue, DecodeArg, and Codec provide context-aware
// conversion at host boundaries. Decode options can require a root runtime
// type, restrict tagged root fields, reject unknown fields, and distinguish
// explicit None from omitted configuration. HostValue represents opaque
// identity, while IterableValue, IteratorValue, SliceView, and MapView opt in
// to only the runtime capabilities they implement.
//
// The sdktest subpackage provides an Engine-backed black-box test harness for
// executing module functions through FQL.
package sdk
