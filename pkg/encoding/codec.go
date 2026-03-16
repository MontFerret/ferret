package encoding

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type (
	// Encoder converts runtime values into bytes.
	Encoder interface {
		Encode(value runtime.Value) ([]byte, error)
		EncodeWith() EncoderConfigurer
	}

	PreEncoderHook  func(value runtime.Value) error
	PostEncoderHook func(value runtime.Value, err error) error

	EncoderConfigurer interface {
		PreHook(PreEncoderHook)
		PostHook(PostEncoderHook)
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
		PreHook(PreDecoderHook)
		PostHook(PostDecoderHook)
		Decoder() Decoder
	}

	// Codec combines encoder and decoder capabilities.
	Codec interface {
		ContentType() string
		Encoder
		Decoder
	}
)
