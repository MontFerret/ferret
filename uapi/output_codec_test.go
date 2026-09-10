package uapi_test

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	emptyOutputCodec struct {
		encodingjson.Codec
	}

	emptyOutputEncoder struct {
		encoding.Encoder
	}

	emptyOutputConfigurer struct {
		encoding.EncoderConfigurer
	}
)

func (c emptyOutputCodec) Encode(value runtime.Value) ([]byte, error) {
	return emptyOutputEncoder{Encoder: c.Codec}.Encode(value)
}

func (c emptyOutputCodec) EncodeWith() encoding.EncoderConfigurer {
	return &emptyOutputConfigurer{EncoderConfigurer: c.Codec.EncodeWith()}
}

func (e emptyOutputEncoder) Encode(value runtime.Value) ([]byte, error) {
	// Exercise the encoder and its ownership hooks, then produce zero bytes.
	_, err := e.Encoder.Encode(value)

	return nil, err
}

func (c *emptyOutputConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
	c.EncoderConfigurer = c.EncoderConfigurer.PreHook(hook)

	return c
}

func (c *emptyOutputConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
	c.EncoderConfigurer = c.EncoderConfigurer.PostHook(hook)

	return c
}

func (c *emptyOutputConfigurer) Encoder() encoding.Encoder {
	return emptyOutputEncoder{Encoder: c.EncoderConfigurer.Encoder()}
}
