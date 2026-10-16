package engine

import (
	"context"
	"io"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type aliasCodec struct {
	base        ferretencoding.Codec
	pre         ferretencoding.PreEncoderHook
	contentType string
}

func (c aliasCodec) ContentType() string {
	return c.contentType
}

func (c aliasCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return c.base.Encode(ctx, dst, value)
}

func (c aliasCodec) EncodeWith() ferretencoding.EncoderConfigurer {
	return c.base.EncodeWith().PreHook(c.pre)
}

func (c aliasCodec) Decode(ctx context.Context, src io.Reader) (runtime.Value, error) {
	return c.base.Decode(ctx, src)
}

func (c aliasCodec) DecodeWith() ferretencoding.DecoderConfigurer {
	return c.base.DecodeWith()
}
