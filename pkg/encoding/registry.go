package encoding

import (
	"fmt"
)

const (
	// ContentTypeJSON is the content type for JSON codec.
	ContentTypeJSON = "application/json"
)

// Registry stores codecs by content type.
type (
	CodecRegistrar interface {
		Register(codec Codec) error
	}

	Registry struct {
		entries map[string]Codec
	}
)

// NewRegistry creates a registry seeded with the provided codecs.
func NewRegistry(codecs ...Codec) *Registry {
	registry := NewEmptyRegistry()

	for _, codec := range codecs {
		_ = registry.Register(codec)
	}

	return registry
}

// NewEmptyRegistry creates a registry without predefined codecs.
func NewEmptyRegistry() *Registry {
	return &Registry{
		entries: make(map[string]Codec),
	}
}

// Register stores a full codec for the content type.
func (r *Registry) Register(codec Codec) error {
	if r == nil {
		return ErrNilRegistry
	}

	if codec == nil {
		return ErrNilCodec
	}

	normalized, err := normalizeContentType(codec.ContentType())
	if err != nil {
		return err
	}

	r.entries[normalized] = codec

	return nil
}

// Codec returns a codec for the content type.
func (r *Registry) Codec(contentType string) (Codec, error) {
	entry, err := r.lookup(contentType)

	if err != nil {
		return nil, err
	}

	return entry, nil
}

// Encoder returns an encoder for the content type.
func (r *Registry) Encoder(contentType string) (Encoder, error) {
	codec, err := r.lookup(contentType)

	if err != nil {
		return nil, err
	}

	return codec, nil
}

// Decoder returns a decoder for the content type.
func (r *Registry) Decoder(contentType string) (Decoder, error) {
	codec, err := r.lookup(contentType)

	if err != nil {
		return nil, err
	}

	return codec, nil
}

func (r *Registry) Clone() *Registry {
	if r == nil {
		return nil
	}

	clone := NewEmptyRegistry()

	for k, v := range r.entries {
		clone.entries[k] = v
	}

	return clone
}

func (r *Registry) lookup(contentType string) (Codec, error) {
	if r == nil {
		return nil, ErrNilRegistry
	}

	normalized, err := normalizeContentType(contentType)
	if err != nil {
		return nil, err
	}

	entry, exists := r.entries[normalized]

	if !exists || entry == nil {
		return nil, fmt.Errorf("%w: %s", ErrCodecNotFound, normalized)
	}

	return entry, nil
}
