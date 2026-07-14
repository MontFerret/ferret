package sdk

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type callbackCodec[T any] struct {
	encoder EncoderFn[T]
	decoder DecoderFn[T]
}

// NewCodec returns a codec backed by the supplied callbacks.
// A nil callback falls back to Encode or Decode independently.
func NewCodec[T any](encoder EncoderFn[T], decoder DecoderFn[T]) Codec[T] {
	return callbackCodec[T]{
		encoder: encoder,
		decoder: decoder,
	}
}

func (codec callbackCodec[T]) Encode(ctx context.Context, input T) (runtime.Value, error) {
	if codec.encoder == nil {
		return Encode(ctx, input)
	}

	return codec.encoder(ctx, input)
}

func (codec callbackCodec[T]) Decode(ctx context.Context, input runtime.Value) (T, error) {
	if codec.decoder == nil {
		var output T
		err := Decode(ctx, input, &output)

		return output, err
	}

	return codec.decoder(ctx, input)
}
