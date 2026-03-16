package msgpack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
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

func (c Codec) Encode(value runtime.Value) ([]byte, error) {
	var buf bytes.Buffer

	enc := vmmsgpack.GetEncoder()
	enc.ResetWriter(&buf)
	defer vmmsgpack.PutEncoder(enc)

	if err := c.encodeValue(context.Background(), enc, value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c Codec) EncodeWith() encoding.EncoderConfigurer {
	return &encoderConfigurer{
		codec: c,
	}
}

func (c Codec) Decode(data []byte) (runtime.Value, error) {
	if err := c.decoder.runPreHooks(data); err != nil {
		return runtime.None, err
	}

	dec := vmmsgpack.GetDecoder()
	dec.ResetReader(bytes.NewReader(data))
	defer vmmsgpack.PutDecoder(dec)

	value, err := c.decodeValue(context.Background(), dec)

	if err == nil {
		if _, peekErr := dec.PeekCode(); peekErr == nil {
			err = fmt.Errorf("msgpack: multiple root values")
		} else if !errors.Is(peekErr, io.EOF) {
			err = peekErr
		}
	}

	if hookErr := c.decoder.runPostHooks(data, err); hookErr != nil {
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
