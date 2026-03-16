package msgpack

import "github.com/MontFerret/ferret/v2/pkg/encoding"

type (
	encoderConfigurer struct {
		codec Codec
		pre   []encoding.PreEncoderHook
		post  []encoding.PostEncoderHook
	}

	decoderConfigurer struct {
		codec Codec
		pre   []encoding.PreDecoderHook
		post  []encoding.PostDecoderHook
	}
)

func (e *encoderConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
	if hook != nil {
		e.pre = append(e.pre, hook)
	}

	return e
}

func (e *encoderConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
	if hook == nil {
		e.post = append(e.post, hook)
	}

	return e
}

func (e *encoderConfigurer) Encoder() encoding.Encoder {
	enc := e.codec.encoder

	if len(e.pre) > 0 {
		enc.pre = append([]encoding.PreEncoderHook(nil), enc.pre...)
		enc.pre = append(enc.pre, e.pre...)
	}

	if len(e.post) > 0 {
		enc.post = append([]encoding.PostEncoderHook(nil), enc.post...)
		enc.post = append(enc.post, e.post...)
	}

	c := e.codec
	c.encoder = enc

	return c
}

func (d *decoderConfigurer) PreHook(hook encoding.PreDecoderHook) encoding.DecoderConfigurer {
	if hook != nil {
		d.pre = append(d.pre, hook)
	}

	return d
}

func (d *decoderConfigurer) PostHook(hook encoding.PostDecoderHook) encoding.DecoderConfigurer {
	if hook == nil {
		d.post = append(d.post, hook)
	}

	return d
}

func (d *decoderConfigurer) Decoder() encoding.Decoder {
	dec := d.codec.decoder

	if len(d.pre) > 0 {
		dec.pre = append([]encoding.PreDecoderHook(nil), dec.pre...)
		dec.pre = append(dec.pre, d.pre...)
	}

	if len(d.post) > 0 {
		dec.post = append([]encoding.PostDecoderHook(nil), dec.post...)
		dec.post = append(dec.post, d.post...)
	}

	c := d.codec
	c.decoder = dec

	return c
}
