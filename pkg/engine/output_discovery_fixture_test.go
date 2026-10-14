package engine

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// hookOwnershipCodec disables the delegated encoder's private callback to
	// isolate custom-codec hook compatibility. It keeps the built-in content type.
	hookOwnershipCodec      struct{ encoding.Codec }
	hookOwnershipEncoder    struct{ encoding.Encoder }
	hookOwnershipConfigurer struct{ encoding.EncoderConfigurer }
)

func (c hookOwnershipCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return c.Codec.Encode(encodingownership.WithValueAdopter(ctx, nil), dst, value)
}

func (c hookOwnershipCodec) EncodeWith() encoding.EncoderConfigurer {
	return &hookOwnershipConfigurer{EncoderConfigurer: c.Codec.EncodeWith()}
}

func (e hookOwnershipEncoder) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return e.Encoder.Encode(encodingownership.WithValueAdopter(ctx, nil), dst, value)
}

func (c *hookOwnershipConfigurer) PreHook(hook encoding.PreEncoderHook) encoding.EncoderConfigurer {
	c.EncoderConfigurer.PreHook(hook)
	return c
}

func (c *hookOwnershipConfigurer) PostHook(hook encoding.PostEncoderHook) encoding.EncoderConfigurer {
	c.EncoderConfigurer.PostHook(hook)
	return c
}

func (c *hookOwnershipConfigurer) Encoder() encoding.Encoder {
	return hookOwnershipEncoder{Encoder: c.EncoderConfigurer.Encoder()}
}
