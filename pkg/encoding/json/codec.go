package json

import (
	"bytes"
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

type (
	// Codec encodes and decodes runtime values as JSON.
	Codec struct {
		config *codecConfig
	}

	// Published configuration is immutable and shared by copies of a codec.
	codecConfig struct {
		encoder
		decoder
	}
)

var (
	// Default is the default JSON codec.
	Default = Codec{}

	defaultEncoder encoder

	encodePool = sync.Pool{New: func() any { return &encodeState{} }}
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

	enc := &defaultEncoder
	if c.config != nil {
		enc = &c.config.encoder
	}

	// Preserve native text for the ordinary string writer when the shortcut
	// declines it, so values such as microsecond durations are formatted once.
	var rootText string
	var rootBuffer *bytes.Buffer
	switch v := value.(type) {
	case runtime.String:
		rootBuffer = enc.bufferForRootString(&op, dst)
		if rootBuffer != nil {
			rootText = string(v)
		}
	case runtime.Duration:
		rootBuffer = enc.bufferForRootString(&op, dst)
		if rootBuffer != nil {
			rootText = v.String()
		}
	case *runtime.Regexp:
		if v != nil {
			rootBuffer = enc.bufferForRootString(&op, dst)
			if rootBuffer != nil {
				rootText = v.String()
			}
		}
	}

	rootString := rootBuffer != nil
	// Complete eligibility scanning before changing the buffer. This bound
	// stays below the ordinary string writer's first payload checkpoint.
	if rootString && len(rootText) <= codecutil.ChunkSize-7 && jsonASCIIEnd(rootText, 0, len(rootText)) == len(rootText) {
		if err := op.Step(); err != nil {
			return err
		}

		rootBuffer.Grow(len(rootText) + 2)
		out := rootBuffer.AvailableBuffer()[:len(rootText)+2]
		out[0] = '"'
		copy(out[1:], rootText)
		out[len(out)-1] = '"'
		_, _ = rootBuffer.Write(out)

		return op.Check()
	}

	state := encodePool.Get().(*encodeState)
	state.writer = codecutil.Writer{Dst: dst, Op: op}

	defer func() {
		state.writer = codecutil.Writer{}
		encodePool.Put(state)
	}()

	if rootString {
		if err := state.writer.Op.Step(); err != nil {
			return err
		}

		if err := state.writeJSONString(rootText); err != nil {
			return err
		}
	} else if err := enc.encodeValue(ctx, state, value); err != nil {
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

	configured := decoder{}
	if c.config != nil {
		configured = c.config.decoder
	}

	if err := configured.runPreHooks(ctx); err != nil {
		return runtime.None, err
	}

	reader := codecutil.Reader{Src: src, Op: op}
	dec := json.NewDecoder(reader.Input())
	dec.UseNumber()

	value, err := configured.decodeValue(ctx, &reader.Op, dec)
	if readErr := reader.Error(); readErr != nil && !errors.Is(err, readErr) {
		err = errors.Join(err, readErr)
	}

	if err == nil {
		err = reader.Op.Check()
	}

	if err != nil {
		value = runtime.None
	}

	if hookErr := configured.runPostHooks(ctx, value, err); hookErr != nil {
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

func (enc *encoder) bufferForRootString(op *codecutil.Operation, dst io.Writer) *bytes.Buffer {
	if len(enc.pre) > 0 || len(enc.post) > 0 || op.HasValueAdopter() {
		return nil
	}

	buffer, _ := dst.(*bytes.Buffer)

	return buffer
}
