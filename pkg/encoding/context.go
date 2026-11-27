package encoding

import (
	"context"
)

type (
	registryContextKey struct{}

	readOnlyRegistry struct {
		registry *Registry
	}
)

var registryCtxKey = registryContextKey{}

// WithRegistry adds a caller-owned registry to the context.
func WithRegistry(ctx context.Context, registry *Registry) context.Context {
	return withRegistry(ctx, registry)
}

// WithReadOnlyRegistry adds a finalized registry to the context. Lookup helpers
// share the registered codecs, and [RegistryFrom] returns a shallow snapshot of
// the registrations. The caller must keep the original registrations fixed.
func WithReadOnlyRegistry(ctx context.Context, registry *Registry) context.Context {
	return withRegistry(ctx, readOnlyRegistry{registry: registry})
}

func withRegistry(ctx context.Context, registry any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, registryCtxKey, registry)
}

// RegistryFrom gets the registry from the context. A registry bound with
// [WithRegistry] is returned directly. A registry bound with [WithReadOnlyRegistry]
// is cloned, so edits to the returned registrations are local to the clone.
// Codec instances are shared in both cases.
func RegistryFrom(ctx context.Context) (*Registry, error) {
	registry, readOnly, err := lookupRegistry(ctx)
	if err != nil {
		return nil, err
	}

	if readOnly {
		return registry.Clone(), nil
	}

	return registry, nil
}

func lookupRegistry(ctx context.Context) (*Registry, bool, error) {
	if ctx != nil {
		switch binding := ctx.Value(registryCtxKey).(type) {
		case *Registry:
			if binding != nil {
				return binding, false, nil
			}
		case readOnlyRegistry:
			if binding.registry != nil {
				return binding.registry, true, nil
			}
		}
	}

	return nil, false, ErrRegistryNotFound
}

// CodecFrom resolves codec by content type from registry in context.
func CodecFrom(ctx context.Context, contentType string) (Codec, error) {
	registry, _, err := lookupRegistry(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Codec(contentType)
}

// EncoderFrom resolves encoder by content type from registry in context.
func EncoderFrom(ctx context.Context, contentType string) (Encoder, error) {
	registry, _, err := lookupRegistry(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Encoder(contentType)
}

// DecoderFrom resolves decoder by content type from registry in context.
func DecoderFrom(ctx context.Context, contentType string) (Decoder, error) {
	registry, _, err := lookupRegistry(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Decoder(contentType)
}
