package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Encoder defines an interface for encoding a value of type T into a runtime.Value within a given context.
	Encoder[T any] interface {
		Encode(context.Context, T) (runtime.Value, error)
	}

	// EncoderFn defines a function type that encodes a value of type T into a Ferret runtime.Value within a given context.
	EncoderFn[T any] func(context.Context, T) (runtime.Value, error)

	// Decoder defines an interface for decoding a runtime.Value into a host representation of type T within a given context.
	Decoder[T any] interface {
		Decode(context.Context, runtime.Value) (T, error)
	}

	// DecoderFn defines a function type that decodes a runtime.Value into a host representation of type T.
	DecoderFn[T any] func(context.Context, runtime.Value) (T, error)

	// Codec represents a combined interface for encoding and decoding values of type T.
	Codec[T any] interface {
		Encoder[T]
		Decoder[T]
	}

	defaultCodec[T any] struct {
		encoder EncoderFn[T]
		decoder DecoderFn[T]
	}
)

// DefaultCodec creates and returns a default implementation of Codec for encoding and decoding values of type T.
func DefaultCodec[T any]() Codec[T] {
	return defaultCodec[T]{
		encoder: func(ctx context.Context, input T) (runtime.Value, error) {
			return Encode(ctx, input)
		},
		decoder: func(ctx context.Context, input runtime.Value) (T, error) {
			var output T
			err := Decode(ctx, input, &output)
			return output, err
		},
	}
}

// NewCodec creates a new Codec instance using the provided encoder and decoder functions for type T.
// If either the encoder or decoder is nil, it defaults to using the Encode and Decode functions for type T.
func NewCodec[T any](encoder EncoderFn[T], decoder DecoderFn[T]) Codec[T] {
	if encoder == nil {
		encoder = func(ctx context.Context, input T) (runtime.Value, error) {
			return Encode(ctx, input)
		}
	}

	if decoder == nil {
		decoder = func(ctx context.Context, input runtime.Value) (T, error) {
			var output T
			err := Decode(ctx, input, &output)
			return output, err
		}
	}

	return defaultCodec[T]{
		encoder: encoder,
		decoder: decoder,
	}
}

func (d defaultCodec[T]) Encode(ctx context.Context, input T) (runtime.Value, error) {
	return d.encoder(ctx, input)
}

func (d defaultCodec[T]) Decode(ctx context.Context, input runtime.Value) (T, error) {
	return d.decoder(ctx, input)
}
