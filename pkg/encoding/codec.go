package encoding

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// Encoder writes a runtime value to a borrowed writer using the caller's
	// context. It must not close or retain the writer, or use it after returning.
	// Write buffers may be reused after Write returns, as required by io.Writer.
	// Success completes all codec-owned writes and flushing; the caller owns
	// flushing, finalizing, and closing its own wrappers. Failure may leave a
	// prefix in dst. Use EncodeBytes and publish only on success for atomic output.
	// Codecs may buffer internally and cannot interrupt arbitrary blocking I/O.
	// Nil contexts and writers must return an error wrapping runtime.ErrInvalidArgument.
	// Preserve non-nil I/O errors; incomplete writes without an error return
	// io.ErrShortWrite and stop without retrying.
	Encoder interface {
		Encode(ctx context.Context, dst io.Writer, value runtime.Value) error
		EncodeWith() EncoderConfigurer
	}

	// PreEncoderHook runs before each visited value, in registration order.
	// Returning an error skips that value's encoding and post-hooks.
	PreEncoderHook func(context.Context, runtime.Value) error
	// PostEncoderHook runs after each visited value, including on encoding failure.
	// Ancestor hooks receive descendant failures. Hooks stop on the first error.
	PostEncoderHook func(context.Context, runtime.Value, error) error

	EncoderConfigurer interface {
		PreHook(PreEncoderHook) EncoderConfigurer
		PostHook(PostEncoderHook) EncoderConfigurer
		Encoder() Encoder
	}

	// Decoder reads one complete document into a materialized runtime value.
	// It borrows src only until Decode returns and must not close or retain it,
	// or read asynchronously. It may buffer and read ahead. Failure may consume
	// part or all of src. The caller owns the reader's lifecycle. Context is
	// propagated, but cannot interrupt arbitrary blocking reader implementations.
	// Nil contexts and readers must return an error wrapping runtime.ErrInvalidArgument.
	// Decoders reject extra roots, malformed input, and truncated containers.
	Decoder interface {
		Decode(ctx context.Context, src io.Reader) (runtime.Value, error)
		DecodeWith() DecoderConfigurer
	}

	// PreDecoderHook runs once before decoding, without consuming the input.
	// Returning an error skips decoding and post-hooks.
	PreDecoderHook func(context.Context) error
	// PostDecoderHook runs once after decoding, with the decoded value on success
	// or runtime.None on failure. Hooks run in order and stop on the first error.
	PostDecoderHook func(context.Context, runtime.Value, error) error

	DecoderConfigurer interface {
		PreHook(PreDecoderHook) DecoderConfigurer
		PostHook(PostDecoderHook) DecoderConfigurer
		Decoder() Decoder
	}

	// Codec combines encoder and decoder capabilities.
	Codec interface {
		ContentType() string
		Encoder
		Decoder
	}
)
