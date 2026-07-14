package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Encoder converts host values of type T into Ferret runtime values.
	Encoder[T any] interface {
		Encode(context.Context, T) (runtime.Value, error)
	}

	// EncoderFn encodes a host value into a Ferret runtime value.
	EncoderFn[T any] func(context.Context, T) (runtime.Value, error)

	// Decoder converts Ferret runtime values into host values of type T.
	Decoder[T any] interface {
		Decode(context.Context, runtime.Value) (T, error)
	}

	// DecoderFn decodes a Ferret runtime value into a host value.
	DecoderFn[T any] func(context.Context, runtime.Value) (T, error)

	// Codec converts host values to and from their Ferret runtime representation.
	Codec[T any] interface {
		Encoder[T]
		Decoder[T]
	}

	defaultCodec[T any] struct{}
)

// DefaultCodec returns a codec backed by Encode and Decode.
func DefaultCodec[T any]() Codec[T] {
	return defaultCodec[T]{}
}

func (defaultCodec[T]) Encode(ctx context.Context, input T) (runtime.Value, error) {
	return Encode(ctx, input)
}

func (defaultCodec[T]) Decode(ctx context.Context, input runtime.Value) (T, error) {
	var output T
	err := Decode(ctx, input, &output)

	return output, err
}
