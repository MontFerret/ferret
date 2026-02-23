package encoding

import (
	"context"
)

type registryContextKey struct{}

var registryCtxKey = registryContextKey{}

// WithRegistry adds registry to context.
func WithRegistry(ctx context.Context, registry *Registry) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, registryCtxKey, registry)
}

// RegistryFrom gets registry from context.
func RegistryFrom(ctx context.Context) (*Registry, error) {
	if ctx == nil {
		return nil, ErrRegistryNotFound
	}

	registry, ok := ctx.Value(registryCtxKey).(*Registry)
	if !ok || registry == nil {
		return nil, ErrRegistryNotFound
	}

	return registry, nil
}

// CodecFrom resolves codec by content type from registry in context.
func CodecFrom(ctx context.Context, contentType string) (Codec, error) {
	registry, err := RegistryFrom(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Codec(contentType)
}

// EncoderFrom resolves encoder by content type from registry in context.
func EncoderFrom(ctx context.Context, contentType string) (Encoder, error) {
	registry, err := RegistryFrom(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Encoder(contentType)
}

// DecoderFrom resolves decoder by content type from registry in context.
func DecoderFrom(ctx context.Context, contentType string) (Decoder, error) {
	registry, err := RegistryFrom(ctx)
	if err != nil {
		return nil, err
	}

	return registry.Decoder(contentType)
}
