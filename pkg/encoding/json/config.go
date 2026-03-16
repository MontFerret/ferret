package json

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

func (d *encoderConfigurer) PreHook(hook encoding.PreEncoderHook) {
	if hook == nil {
		return
	}

	d.pre = append(d.pre, hook)
}

func (d *encoderConfigurer) PostHook(hook encoding.PostEncoderHook) {
	if hook == nil {
		return
	}

	d.post = append(d.post, hook)
}

func (d *encoderConfigurer) Encoder() encoding.Encoder {
	enc := d.codec.encoder

	if len(d.pre) > 0 {
		enc.pre = append([]encoding.PreEncoderHook(nil), enc.pre...)
		enc.pre = append(enc.pre, d.pre...)
	}

	if len(d.post) > 0 {
		enc.post = append([]encoding.PostEncoderHook(nil), enc.post...)
		enc.post = append(enc.post, d.post...)
	}

	c := d.codec
	c.encoder = enc

	return c
}

func (d *decoderConfigurer) PreHook(hook encoding.PreDecoderHook) {
	if hook == nil {
		return
	}

	d.pre = append(d.pre, hook)
}

func (d *decoderConfigurer) PostHook(hook encoding.PostDecoderHook) {
	if hook == nil {
		return
	}

	d.post = append(d.post, hook)
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
