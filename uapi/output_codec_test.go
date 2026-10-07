package uapi_test

import (
	"bytes"
	"sync"

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

	// This custom codec intentionally reuses private storage, detaching inside
	// Encode under exclusive access, including on error and configured paths.
	detachedTestCodec struct {
		err error
		encodingjson.Codec
		data     []byte
		returned []byte
		mu       sync.Mutex
	}

	detachedTestEncoder struct {
		codec *detachedTestCodec
		encoding.Encoder
	}

	detachedTestConfigurer struct {
		codec *detachedTestCodec
		base  encoding.EncoderConfigurer
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

func (c *detachedTestCodec) ContentType() string { return "application/x-detached-test" }

func (c *detachedTestCodec) Encode(value runtime.Value) ([]byte, error) {
	return (detachedTestEncoder{codec: c, Encoder: c.Codec}).Encode(value)
}

func (c *detachedTestCodec) EncodeWith() encoding.EncoderConfigurer {
	return &detachedTestConfigurer{codec: c, base: c.Codec.EncodeWith()}
}

func (c *detachedTestConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
	c.base = c.base.PreHook(hook)

	return c
}

func (c *detachedTestConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
	c.base = c.base.PostHook(hook)

	return c
}

func (c *detachedTestConfigurer) Encoder() encoding.Encoder {
	return detachedTestEncoder{codec: c.codec, Encoder: c.base.Encoder()}
}

func (e detachedTestEncoder) Encode(value runtime.Value) ([]byte, error) {
	if _, err := e.Encoder.Encode(value); err != nil {
		return nil, err
	}

	e.codec.mu.Lock()
	defer e.codec.mu.Unlock()
	e.codec.returned = bytes.Clone(e.codec.data)

	return e.codec.returned, e.codec.err
}
