// Package core provides v1-compatible runtime core types for the Ferret compatibility layer.
// It mirrors the github.com/MontFerret/ferret/pkg/runtime/core package from Ferret v1,
// wrapping v2 internals so that v1-style code can be migrated with minimal changes.
package core

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"io"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MaxArgs defines the maximum number of arguments that a function can accept.
const MaxArgs = runtime.MaxArgs

// CloseFunc is a function that closes a resource.
type CloseFunc func() error

// Value mirrors the v1 core.Value interface.
// All values returned by the compat layer implement this interface.
type Value interface {
	json.Marshaler
	Type() Type
	String() string
	Compare(other Value) int64
	Unwrap() interface{}
	Hash() uint64
	Copy() Value
}

// Type mirrors the v1 core.Type interface.
type Type interface {
	ID() int64
	String() string
	Equals(other Type) bool
}

// Function mirrors the v1 core.Function type.
type Function = func(ctx context.Context, args ...Value) (Value, error)

// Iterator mirrors the v1 core.Iterator interface.
type Iterator interface {
	Next(ctx context.Context) (value Value, key Value, err error)
}

// Iterable mirrors the v1 core.Iterable interface.
type Iterable interface {
	Iterate(ctx context.Context) (Iterator, error)
}

// --- Value adapter ---

// valueAdapter wraps a v2 runtime.Value and implements the compat Value interface.
type valueAdapter struct {
	inner runtime.Value
}

// WrapValue converts a v2 runtime.Value into a compat Value.
func WrapValue(v runtime.Value) Value {
	if v == nil {
		return nil
	}

	return &valueAdapter{inner: v}
}

// UnwrapValue converts a compat Value back to a v2 runtime.Value.
func UnwrapValue(v Value) runtime.Value {
	if v == nil {
		return runtime.None
	}

	if a, ok := v.(*valueAdapter); ok {
		return a.inner
	}

	// Custom core.Value that doesn't wrap a v2 value – wrap it for v2 consumption.
	return &coreValueAsRuntimeValue{inner: v}
}

func (a *valueAdapter) Type() Type {
	return WrapType(runtime.TypeOf(a.inner))
}

func (a *valueAdapter) String() string {
	return a.inner.String()
}

func (a *valueAdapter) Compare(other Value) int64 {
	otherRT := UnwrapValue(other)
	if cmp, ok := a.inner.(runtime.Comparable); ok {
		return int64(cmp.Compare(otherRT))
	}

	return int64(runtime.CompareValues(a.inner, otherRT))
}

func (a *valueAdapter) Unwrap() interface{} {
	if u, ok := a.inner.(runtime.Unwrappable); ok {
		return u.Unwrap()
	}

	// Fallback: return the underlying runtime.Value when it cannot be further unwrapped.
	return a.inner
}

func (a *valueAdapter) Hash() uint64 {
	return a.inner.Hash()
}

func (a *valueAdapter) Copy() Value {
	return WrapValue(a.inner.Copy())
}

func (a *valueAdapter) MarshalJSON() ([]byte, error) {
	return encodingjson.Default.Encode(a.inner)
}

// coreValueAsRuntimeValue wraps a compat Value so it satisfies runtime.Value.
// Used when a user-defined core.Value needs to be passed into v2 engine.
type coreValueAsRuntimeValue struct {
	inner Value
}

func (c *coreValueAsRuntimeValue) String() string {
	return c.inner.String()
}

func (c *coreValueAsRuntimeValue) Hash() uint64 {
	return c.inner.Hash()
}

func (c *coreValueAsRuntimeValue) Copy() runtime.Value {
	return UnwrapValue(c.inner.Copy())
}

// --- Type adapter ---

// typeAdapter wraps a v2 runtime.Type and implements the compat Type interface.
type typeAdapter struct {
	inner runtime.Type
}

// WrapType converts a v2 runtime.Type into a compat Type.
func WrapType(t runtime.Type) Type {
	if t == nil {
		return nil
	}

	return &typeAdapter{inner: t}
}

func (t *typeAdapter) ID() int64 {
	h := fnv.New64a()
	h.Write([]byte(t.inner.Name()))
	// truncate to int64 range by interpreting raw bytes
	return int64(h.Sum64())
}

func (t *typeAdapter) String() string {
	return t.inner.String()
}

func (t *typeAdapter) Equals(other Type) bool {
	if other == nil {
		return false
	}

	return t.String() == other.String()
}

// --- Function adapters ---

// WrapFunction converts a v2 runtime.Function into a compat Function.
// The returned function wraps v2 values into compat Values and unwraps them back.
func WrapFunction(fn runtime.Function) Function {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, args ...Value) (Value, error) {
		rtArgs := make([]runtime.Value, len(args))
		for i, a := range args {
			rtArgs[i] = UnwrapValue(a)
		}

		result, err := fn(ctx, rtArgs...)
		if err != nil {
			return nil, err
		}

		return WrapValue(result), nil
	}
}

// UnwrapFunction converts a compat Function into a v2 runtime.Function.
func UnwrapFunction(fn Function) runtime.Function {
	if fn == nil {
		return nil
	}

	return func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		coreArgs := make([]Value, len(args))
		for i, a := range args {
			coreArgs[i] = WrapValue(a)
		}

		result, err := fn(ctx, coreArgs...)
		if err != nil {
			return nil, err
		}

		return UnwrapValue(result), nil
	}
}

// --- Iterator adapter ---

// iteratorAdapter wraps a v2 runtime.Iterator and implements the compat Iterator interface.
type iteratorAdapter struct {
	inner runtime.Iterator
}

// WrapIterator converts a v2 runtime.Iterator into a compat Iterator.
func WrapIterator(it runtime.Iterator) Iterator {
	if it == nil {
		return nil
	}

	return &iteratorAdapter{inner: it}
}

func (it *iteratorAdapter) Next(ctx context.Context) (Value, Value, error) {
	val, key, err := it.inner.Next(ctx)
	if err != nil {
		if err == io.EOF {
			return nil, nil, err
		}

		return nil, nil, err
	}

	return WrapValue(val), WrapValue(key), nil
}

// iterableAdapter wraps a v2 runtime.Iterable and implements the compat Iterable interface.
type iterableAdapter struct {
	inner runtime.Iterable
}

// WrapIterable converts a v2 runtime.Iterable into a compat Iterable.
func WrapIterable(it runtime.Iterable) Iterable {
	if it == nil {
		return nil
	}

	return &iterableAdapter{inner: it}
}

func (it *iterableAdapter) Iterate(ctx context.Context) (Iterator, error) {
	iter, err := it.inner.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	return WrapIterator(iter), nil
}
