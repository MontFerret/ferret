package json

import (
	"bytes"
	"context"
	"errors"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/encoding"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	// ContentType is the content type for JSON codec.
	ContentType = "application/json"
)

// Codec encodes and decodes runtime values as JSON.
type Codec struct {
	encoder
	decoder
}

// Default is the default JSON codec.
var Default = Codec{}

func (Codec) ContentType() string {
	return ContentType
}

func (c Codec) Encode(value runtime.Value) ([]byte, error) {
	var buf bytes.Buffer

	if err := c.encodeValue(context.Background(), &buf, value); err != nil {
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

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	value, err := c.decodeValue(context.Background(), dec)

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
