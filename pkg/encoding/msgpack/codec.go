package msgpack

import (
	"context"
	"errors"
	"fmt"
	"io"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	// ContentType is the content type for MessagePack codec.
	ContentType = "application/vnd.msgpack"
)

// Codec encodes and decodes runtime values as MessagePack.
type Codec struct {
	encoder
	decoder
}

// Default is the default MessagePack codec.
var Default = Codec{}

func (Codec) ContentType() string {
	return ContentType
}

func (c Codec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	op, err := codecutil.NewOperation(ctx, dst, "writer")
	if err != nil {
		return err
	}

	op.CacheValueAdopter()
	writer := codecutil.Writer{Dst: dst, Op: op}

	enc := vmmsgpack.GetEncoder()
	enc.Reset(&writer)
	defer vmmsgpack.PutEncoder(enc)

	if err := c.encodeValue(ctx, &writer.Op, enc, value); err != nil {
		return err
	}

	if err := writer.Error(); err != nil {
		return err
	}

	return writer.Op.Check()
}

func (c Codec) EncodeWith() encoding.EncoderConfigurer {
	return &encoderConfigurer{
		codec: c,
	}
}

func (c Codec) Decode(ctx context.Context, src io.Reader) (runtime.Value, error) {
	op, err := codecutil.NewOperation(ctx, src, "reader")
	if err != nil {
		return runtime.None, err
	}

	if err := c.decoder.runPreHooks(ctx); err != nil {
		return runtime.None, err
	}

	dec := vmmsgpack.GetDecoder()
	reader := codecutil.Reader{Src: src, Op: op}
	dec.Reset(reader.Input())
	defer vmmsgpack.PutDecoder(dec)

	value, err := c.decodeValue(ctx, &reader.Op, dec)

	if err == nil {
		if _, peekErr := dec.PeekCode(); peekErr == nil {
			err = fmt.Errorf("msgpack: multiple root values")
		} else if !errors.Is(peekErr, io.EOF) {
			err = peekErr
		}
	}

	if readErr := reader.Error(); readErr != nil && !errors.Is(err, readErr) {
		err = errors.Join(err, readErr)
	}

	if err == nil {
		err = reader.Op.Check()
	}

	if err != nil {
		value = runtime.None
	}

	if hookErr := c.decoder.runPostHooks(ctx, value, err); hookErr != nil {
		return runtime.None, errors.Join(err, hookErr)
	}

	if err != nil {
		return runtime.None, err
	}

	return value, nil
}

func (c Codec) DecodeWith() encoding.DecoderConfigurer {
	return &decoderConfigurer{
		codec: c,
	}
}
