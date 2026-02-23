package encoding

import (
	"fmt"
	"sync"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
)

const (
	// ContentTypeJSON is the content type for JSON codec.
	ContentTypeJSON = "application/json"
)

// Registry stores codecs by content type.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Codec
}

// NewRegistry creates a registry with built-in codecs.
func NewRegistry() *Registry {
	registry := NewEmptyRegistry()
	_ = registry.Register(ContentTypeJSON, encodingjson.Default)

	return registry
}

// NewEmptyRegistry creates a registry without predefined codecs.
func NewEmptyRegistry() *Registry {
	return &Registry{
		entries: make(map[string]Codec),
	}
}

// Register stores a full codec for the content type.
func (r *Registry) Register(contentType string, codec Codec) error {
	if r == nil {
		return ErrNilRegistry
	}

	if codec == nil {
		return ErrNilCodec
	}

	normalized, err := normalizeContentType(contentType)
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.entries[normalized] = codec
	r.mu.Unlock()

	return nil
}

// Codec returns a codec for the content type.
func (r *Registry) Codec(contentType string) (Codec, error) {
	entry, normalized, err := r.lookup(contentType)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("%w: %s", ErrCodecNotFound, normalized)
	}

	return entry, nil
}

// Encoder returns an encoder for the content type.
func (r *Registry) Encoder(contentType string) (Encoder, error) {
	codec, normalized, err := r.lookup(contentType)
	if err != nil {
		return nil, err
	}

	if codec == nil {
		return nil, fmt.Errorf("%w: %s", ErrCodecNotFound, normalized)
	}

	return codec, nil
}

// Decoder returns a decoder for the content type.
func (r *Registry) Decoder(contentType string) (Decoder, error) {
	codec, normalized, err := r.lookup(contentType)
	if err != nil {
		return nil, err
	}

	if codec == nil {
		return nil, fmt.Errorf("%w: %s", ErrCodecNotFound, normalized)
	}

	return codec, nil
}

func (r *Registry) lookup(contentType string) (Codec, string, error) {
	if r == nil {
		return nil, "", ErrNilRegistry
	}

	normalized, err := normalizeContentType(contentType)
	if err != nil {
		return nil, "", err
	}

	r.mu.RLock()
	entry, exists := r.entries[normalized]
	r.mu.RUnlock()

	if !exists {
		return nil, normalized, fmt.Errorf("%w: %s", ErrCodecNotFound, normalized)
	}

	return entry, normalized, nil
}
