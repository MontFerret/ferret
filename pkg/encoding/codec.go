package encoding

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type (
	// Encoder converts runtime values into bytes. Every returned byte slice,
	// including bytes returned with an error, transfers ownership to the caller.
	// It must be detached from reusable encoder storage and borrowed runtime values.
	// Encoders using shared storage must detach while they still have exclusive
	// access, before releasing a lock or recycling the storage.
	Encoder interface {
		Encode(value runtime.Value) ([]byte, error)
		EncodeWith() EncoderConfigurer
	}

	PreEncoderHook  func(value runtime.Value) error
	PostEncoderHook func(value runtime.Value, err error) error

	EncoderConfigurer interface {
		PreHook(PreEncoderHook) EncoderConfigurer
		PostHook(PostEncoderHook) EncoderConfigurer
		Encoder() Encoder
	}

	// Decoder converts bytes into runtime values.
	Decoder interface {
		Decode(data []byte) (runtime.Value, error)
		DecodeWith() DecoderConfigurer
	}

	PreDecoderHook  func(data []byte) error
	PostDecoderHook func(data []byte, err error) error

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
