package uapi_test

import (
	"context"
	"io"
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

	// This custom codec writes reusable private storage synchronously while
	// holding exclusive access, including on error and configured paths.
	detachedTestCodec struct {
		err error
		encodingjson.Codec
		data    []byte
		written bool
		mu      sync.Mutex
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

func (c emptyOutputCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return emptyOutputEncoder{Encoder: c.Codec}.Encode(ctx, dst, value)
}

func (c emptyOutputCodec) EncodeWith() encoding.EncoderConfigurer {
	return &emptyOutputConfigurer{EncoderConfigurer: c.Codec.EncodeWith()}
}

func (e emptyOutputEncoder) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	// Exercise the encoder and its ownership hooks, then produce zero bytes.
	return e.Encoder.Encode(ctx, io.Discard, value)
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

func (c *detachedTestCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return (detachedTestEncoder{codec: c, Encoder: c.Codec}).Encode(ctx, dst, value)
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

func (e detachedTestEncoder) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	if err := e.Encoder.Encode(ctx, io.Discard, value); err != nil {
		return err
	}

	e.codec.mu.Lock()
	defer e.codec.mu.Unlock()
	e.codec.written = true
	n, err := dst.Write(e.codec.data)
	if err == nil && n != len(e.codec.data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return err
	}
	return e.codec.err
}
