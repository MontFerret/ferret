package msgpack

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	// ContentType is the content type for MessagePack codec.
	ContentType = "application/vnd.msgpack"
)

type (
	// Codec encodes and decodes runtime values as MessagePack.
	Codec struct {
		config *codecConfig
	}

	encodeState struct {
		encoder *vmmsgpack.Encoder
		writer  codecutil.Writer
	}
)

var (
	// Default is the default MessagePack codec.
	Default = Codec{}

	defaultEncoder encoder

	encodePool = sync.Pool{New: func() any {
		return &encodeState{encoder: vmmsgpack.NewEncoder(nil)}
	}}
)

func (Codec) ContentType() string {
	return ContentType
}

func (c Codec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	op, err := codecutil.NewOperation(ctx, dst, "writer")
	if err != nil {
		return err
	}

	// Register stored ownership before cancellation can abandon the root.
	op.CacheValueAdopter()
	op.Adopt(value, true)

	if err := op.Check(); err != nil {
		return err
	}

	state := encodePool.Get().(*encodeState)
	state.writer = codecutil.Writer{Dst: dst, Op: op}
	state.encoder.Reset(&state.writer)
	defer func() {
		// Clear borrowed I/O, contexts, callbacks, and custom encoder settings.
		state.encoder.Reset(nil)
		state.writer = codecutil.Writer{}
		encodePool.Put(state)
	}()

	if err := c.encoderConfig().encodeValue(ctx, &state.writer.Op, state.encoder, value); err != nil {
		return err
	}

	if err := state.writer.Error(); err != nil {
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

	if err := op.Check(); err != nil {
		return runtime.None, err
	}

	decoder := c.decoderConfig()
	if err := decoder.runPreHooks(ctx); err != nil {
		return runtime.None, err
	}

	dec := vmmsgpack.GetDecoder()
	reader := codecutil.Reader{Src: src, Op: op}
	dec.Reset(reader.Input())
	defer vmmsgpack.PutDecoder(dec)

	value, err := decoder.decodeValue(ctx, &reader.Op, dec)

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

	if hookErr := decoder.runPostHooks(ctx, value, err); hookErr != nil {
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

func (c Codec) encoderConfig() *encoder {
	if c.config == nil {
		return &defaultEncoder
	}

	return &c.config.encoder
}

func (c Codec) decoderConfig() decoder {
	if c.config == nil {
		return decoder{}
	}

	return c.config.decoder
}
