package json

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
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

var encodePool = sync.Pool{New: func() any { return &encodeState{} }}

func (Codec) ContentType() string {
	return ContentType
}

func (c Codec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	op, err := codecutil.NewOperation(ctx, dst, "writer")
	if err != nil {
		return err
	}

	state := encodePool.Get().(*encodeState)
	state.writer = codecutil.Writer{Dst: dst, Op: op}

	defer func() {
		state.writer = codecutil.Writer{}
		encodePool.Put(state)
	}()

	if err := c.encodeValue(ctx, state, value); err != nil {
		return err
	}

	return state.writer.Op.Check()
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

	reader := codecutil.Reader{Src: src, Op: op}
	dec := json.NewDecoder(reader.Input())
	dec.UseNumber()

	value, err := c.decodeValue(ctx, &reader.Op, dec)
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
