package json

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/internal/codecutil"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const jsonHex = "0123456789abcdef"

type (
	encoder struct {
		pre  []encoding.PreEncoderHook
		post []encoding.PostEncoderHook
	}

	encodeState struct {
		scratch []byte
		writer  codecutil.Writer
		number  [32]byte
	}
)

func (enc *encoder) encodeValue(ctx context.Context, state *encodeState, value runtime.Value) error {
	state.writer.Op.Adopt(value, len(enc.pre) > 0 || len(enc.post) > 0)

	if err := state.writer.Op.Step(); err != nil {
		state.writer.Op.Adopt(value, true)

		return err
	}

	if len(enc.pre) > 0 {
		if err := enc.runPreHooks(ctx, value); err != nil {
			state.writer.Op.Adopt(value, true)

			return err
		}
	}

	var err error

	if value == nil || value == runtime.None {
		err = state.writeText("null")
	} else if marshaler, ok := value.(json.Marshaler); ok {
		err = state.writeMarshaler(marshaler)
	} else {
		switch v := value.(type) {
		case runtime.Boolean:
			if v {
				err = state.writeText("true")
			} else {
				err = state.writeText("false")
			}
		case runtime.Int:
			err = state.write(strconv.AppendInt(state.numberBuffer(), int64(v), 10))
		case runtime.Float:
			err = state.write(strconv.AppendFloat(state.numberBuffer(), float64(v), 'g', -1, 64))
		case runtime.Duration:
			err = state.writeJSONString(v.String())
		case runtime.String:
			err = state.writeJSONString(string(v))
		case *runtime.Regexp:
			err = state.writeJSONString(v.String())
		case runtime.Binary:
			err = state.writeBinary(v)
		case runtime.DateTime:
			err = state.writeJSONString(v.Time.Format(time.RFC3339Nano))
		case *runtime.Range:
			err = enc.encodeRange(ctx, state, v)
		case runtime.Map:
			err = enc.encodeMap(ctx, state, v)
		case runtime.List:
			err = enc.encodeList(ctx, state, v)
		case runtime.Iterable:
			err = enc.encodeIterable(ctx, state, v)
		case runtime.Unwrappable:
			err = enc.encodeAny(state, v.Unwrap())
		default:
			err = enc.encodeAny(state, v)
		}
	}

	if value == nil {
		value = runtime.None
	}

	if len(enc.post) > 0 {
		state.writer.Op.Adopt(value, true)

		if hookErr := enc.runPostHooks(ctx, value, err); hookErr != nil {
			err = errors.Join(err, hookErr)
		}
	}

	if err != nil || len(enc.pre) > 0 || len(enc.post) > 0 {
		state.writer.Op.Adopt(value, true)
	}

	return err
}

func (enc *encoder) encodeAny(state *encodeState, value any) error {
	switch v := value.(type) {
	case json.RawMessage:
		return state.write(v)
	case string:
		return state.writeJSONString(v)
	case []byte:
		return state.writeBinary(v)
	case time.Time:
		return state.writeJSONString(v.Format(time.RFC3339Nano))
	}

	if err := state.writer.Op.Check(); err != nil {
		return err
	}

	out, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := state.writer.Op.Check(); err != nil {
		return err
	}

	return state.write(out)
}

func (enc *encoder) encodeMap(ctx context.Context, state *encodeState, value runtime.Map) error {
	if err := state.writer.WriteByte('{'); err != nil {
		return err
	}

	_, stored := value.(*runtime.Object)
	first := true
	err := value.ForEach(ctx, func(ctx context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		state.writer.Op.Adopt(val, !stored)
		state.writer.Op.Adopt(key, !stored)

		if !first {
			if err := state.writer.WriteByte(','); err != nil {
				state.writer.Op.Adopt(val, true)
				state.writer.Op.Adopt(key, true)

				return false, err
			}
		}

		first = false

		if err := state.writeJSONString(key.String()); err != nil {
			state.writer.Op.Adopt(val, true)
			state.writer.Op.Adopt(key, true)

			return false, err
		}

		if err := state.writer.WriteByte(':'); err != nil {
			state.writer.Op.Adopt(val, true)
			state.writer.Op.Adopt(key, true)

			return false, err
		}

		if err := enc.encodeValue(ctx, state, val); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	return state.writer.WriteByte('}')
}

func (enc *encoder) encodeList(ctx context.Context, state *encodeState, value runtime.List) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := state.writer.WriteByte('['); err != nil {
		return err
	}

	_, stored := value.(*runtime.Array)
	for i := runtime.Int(0); i < length; i++ {
		item, err := value.At(ctx, i)
		if err != nil {
			return err
		}

		state.writer.Op.Adopt(item, !stored)

		if i > 0 {
			if err := state.writer.WriteByte(','); err != nil {
				state.writer.Op.Adopt(item, true)

				return err
			}
		}

		if err := enc.encodeValue(ctx, state, item); err != nil {
			return err
		}
	}

	return state.writer.WriteByte(']')
}

func (enc *encoder) encodeIterable(ctx context.Context, state *encodeState, value runtime.Iterable) (err error) {
	if err := state.writer.WriteByte('['); err != nil {
		return err
	}

	iter, err := codecutil.NewIterator(ctx, value)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := iter.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	first := true
	err = runtime.ForEachIter(ctx, iter.Iterator, func(ctx context.Context, item, key runtime.Value) (runtime.Boolean, error) {
		state.writer.Op.Adopt(item, true)
		state.writer.Op.Adopt(key, true)

		if !first {
			if err := state.writer.WriteByte(','); err != nil {
				state.writer.Op.Adopt(item, true)
				state.writer.Op.Adopt(key, true)

				return false, err
			}
		}

		first = false

		if err := enc.encodeValue(ctx, state, item); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return err
	}

	return state.writer.WriteByte(']')
}

func (enc *encoder) encodeRange(ctx context.Context, state *encodeState, value *runtime.Range) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := state.writer.WriteByte('['); err != nil {
		return err
	}

	start := value.Start()
	ascending := start <= value.End()
	hooked := len(enc.pre) > 0 || len(enc.post) > 0
	_, buffered := state.writer.Dst.(*bytes.Buffer)
	chunked := buffered && !hooked
	var chunk []byte
	if chunked {
		// Concrete buffers have no write callbacks or errors whose timing could change.
		chunk = state.scratchFor(codecutil.ChunkSize)
	}

	for offset := runtime.Int(0); offset < length; offset++ {
		if offset > 0 && !chunked {
			if err := state.writer.WriteByte(','); err != nil {
				return err
			}
		}

		item := start + offset
		if !ascending {
			item = start - offset
		}

		if !hooked {
			if err := state.writer.Op.Step(); err != nil {
				return err
			}

			output := state.number[:0]
			if chunked {
				output = chunk
				if offset > 0 {
					output = append(output, ',')
				}
			}

			output = strconv.AppendInt(output, int64(item), 10)
			if chunked && len(output) <= codecutil.ChunkSize-len(state.number) {
				chunk = output

				continue
			}

			if err := state.write(output); err != nil {
				return err
			}

			chunk = chunk[:0]

			continue
		}

		if err := enc.encodeValue(ctx, state, item); err != nil {
			return err
		}
	}

	if len(chunk) > 0 {
		if err := state.write(chunk); err != nil {
			return err
		}
	}

	return state.writer.WriteByte(']')
}

func (state *encodeState) write(data []byte) error {
	_, err := state.writer.Write(data)

	return err
}

func (state *encodeState) writeText(text string) error {
	return state.write(append(state.number[:0], text...))
}

func (state *encodeState) writeMarshaler(marshaler json.Marshaler) error {
	if err := state.writer.Op.Check(); err != nil {
		return err
	}

	out, err := marshaler.MarshalJSON()
	if err != nil {
		return err
	}

	if err := state.writer.Op.Check(); err != nil {
		return err
	}

	return state.write(out)
}

func (state *encodeState) numberBuffer() []byte {
	if buffer, ok := state.writer.Dst.(*bytes.Buffer); ok {
		available := buffer.AvailableBuffer()
		if cap(available) >= len(state.number) {
			return available
		}
	}

	return state.number[:0]
}

func (state *encodeState) scratchFor(size int) []byte {
	if cap(state.scratch) < size {
		state.scratch = make([]byte, 0, size)
	}

	return state.scratch[:0]
}

func (state *encodeState) writeJSONString(text string) error {
	size := codecutil.ChunkSize
	if len(text) < (codecutil.ChunkSize-2)/6 {
		size = 6*len(text) + 2
	}

	scratch := append(state.scratchFor(size), '"')
	// Reserve run scanning for strings large enough to amortize its setup.
	asciiRuns := len(text) >= 64

	for offset := 0; offset < len(text); {
		// Reserve the longest escape and the closing quote within one chunk.
		if len(scratch) > codecutil.ChunkSize-7 {
			if err := state.writer.Op.Check(); err != nil {
				return err
			}

			if err := state.write(scratch); err != nil {
				return err
			}

			scratch = scratch[:0]
		}

		if asciiRuns {
			length := min(len(text)-offset, codecutil.ChunkSize-6-len(scratch))
			end := jsonASCIIEnd(text, offset, offset+length)
			if end > offset {
				scratch = append(scratch, text[offset:end]...)
				offset = end

				continue
			}
		}

		r, size := rune(text[offset]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRuneInString(text[offset:])
		}

		offset += size

		switch r {
		case '\\', '"':
			scratch = append(scratch, '\\', byte(r))
		case '\b':
			scratch = append(scratch, '\\', 'b')
		case '\f':
			scratch = append(scratch, '\\', 'f')
		case '\n':
			scratch = append(scratch, '\\', 'n')
		case '\r':
			scratch = append(scratch, '\\', 'r')
		case '\t':
			scratch = append(scratch, '\\', 't')
		default:
			if r < 0x20 {
				scratch = append(scratch, '\\', 'u', '0', '0', jsonHex[r>>4], jsonHex[r&0x0f])
			} else {
				scratch = utf8.AppendRune(scratch, r)
			}
		}
	}

	return state.write(append(scratch, '"'))
}

func (state *encodeState) writeBinary(data []byte) error {
	if err := state.writer.WriteByte('"'); err != nil {
		return err
	}

	const inputChunk = codecutil.ChunkSize / 4 * 3
	scratch := state.scratchFor(base64.StdEncoding.EncodedLen(min(len(data), inputChunk)))

	for len(data) > 0 {
		if err := state.writer.Op.Check(); err != nil {
			return err
		}

		size := min(len(data), inputChunk)
		output := scratch[:base64.StdEncoding.EncodedLen(size)]
		base64.StdEncoding.Encode(output, data[:size])

		if err := state.write(output); err != nil {
			return err
		}

		data = data[size:]
	}

	return state.writer.WriteByte('"')
}

func (enc *encoder) runPreHooks(ctx context.Context, value runtime.Value) error {
	for _, hook := range enc.pre {
		if err := hook(ctx, value); err != nil {
			return err
		}
	}

	return nil
}

func (enc *encoder) runPostHooks(ctx context.Context, value runtime.Value, err error) error {
	for _, hook := range enc.post {
		if hookErr := hook(ctx, value, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
