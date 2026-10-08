package json

import (
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

func (enc encoder) encodeValue(ctx context.Context, state *encodeState, value runtime.Value) error {
	if err := state.writer.Op.Step(); err != nil {
		return err
	}

	if err := enc.runPreHooks(ctx, value); err != nil {
		return err
	}

	err := enc.encodeBody(ctx, state, value)
	if value == nil {
		value = runtime.None
	}

	if hookErr := enc.runPostHooks(ctx, value, err); hookErr != nil {
		return errors.Join(err, hookErr)
	}

	return err
}

func (enc encoder) encodeBody(ctx context.Context, state *encodeState, value runtime.Value) error {
	if value == nil || value == runtime.None {
		return state.writeText("null")
	}

	if marshaler, ok := value.(json.Marshaler); ok {
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

	switch v := value.(type) {
	case runtime.Boolean:
		if v {
			return state.writeText("true")
		}

		return state.writeText("false")
	case runtime.Int:
		return state.write(strconv.AppendInt(state.number[:0], int64(v), 10))
	case runtime.Float:
		return state.write(strconv.AppendFloat(state.number[:0], float64(v), 'g', -1, 64))
	case runtime.Duration:
		return state.writeJSONString(v.String())
	case runtime.String:
		return state.writeJSONString(string(v))
	case *runtime.Regexp:
		return state.writeJSONString(v.String())
	case runtime.Binary:
		return state.writeBinary(v)
	case runtime.DateTime:
		return state.writeJSONString(v.Time.Format(time.RFC3339Nano))
	case *runtime.Range:
		return enc.encodeRange(ctx, state, v)
	case runtime.Map:
		return enc.encodeMap(ctx, state, v)
	case runtime.List:
		return enc.encodeList(ctx, state, v)
	case runtime.Iterable:
		return enc.encodeIterable(ctx, state, v)
	case runtime.Unwrappable:
		return enc.encodeAny(state, v.Unwrap())
	default:
		return enc.encodeAny(state, v)
	}
}

func (enc encoder) encodeAny(state *encodeState, value any) error {
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

func (enc encoder) encodeMap(ctx context.Context, state *encodeState, value runtime.Map) error {
	if err := state.writer.WriteByte('{'); err != nil {
		return err
	}

	first := true
	err := value.ForEach(ctx, func(ctx context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		if !first {
			if err := state.writer.WriteByte(','); err != nil {
				return false, err
			}
		}

		first = false

		if err := state.writeJSONString(key.String()); err != nil {
			return false, err
		}

		if err := state.writer.WriteByte(':'); err != nil {
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

func (enc encoder) encodeList(ctx context.Context, state *encodeState, value runtime.List) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := state.writer.WriteByte('['); err != nil {
		return err
	}

	for i := runtime.Int(0); i < length; i++ {
		if i > 0 {
			if err := state.writer.WriteByte(','); err != nil {
				return err
			}
		}

		item, err := value.At(ctx, i)
		if err != nil {
			return err
		}

		if err := enc.encodeValue(ctx, state, item); err != nil {
			return err
		}
	}

	return state.writer.WriteByte(']')
}

func (enc encoder) encodeIterable(ctx context.Context, state *encodeState, value runtime.Iterable) (err error) {
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
	err = runtime.ForEachIter(ctx, iter.Iterator, func(ctx context.Context, item, _ runtime.Value) (runtime.Boolean, error) {
		if !first {
			if err := state.writer.WriteByte(','); err != nil {
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

func (enc encoder) encodeRange(ctx context.Context, state *encodeState, value *runtime.Range) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	if err := state.writer.WriteByte('['); err != nil {
		return err
	}

	start := value.Start()
	ascending := start <= value.End()
	for offset := runtime.Int(0); offset < length; offset++ {
		if offset > 0 {
			if err := state.writer.WriteByte(','); err != nil {
				return err
			}
		}

		item := start + offset
		if !ascending {
			item = start - offset
		}

		if len(enc.pre) == 0 && len(enc.post) == 0 {
			if err := state.writer.Op.Step(); err != nil {
				return err
			}

			if err := state.write(strconv.AppendInt(state.number[:0], int64(item), 10)); err != nil {
				return err
			}

			continue
		}

		if err := enc.encodeValue(ctx, state, item); err != nil {
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
	return state.write([]byte(text))
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

	for _, r := range text {
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

func (enc encoder) runPreHooks(ctx context.Context, value runtime.Value) error {
	for _, hook := range enc.pre {
		if err := hook(ctx, value); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPostHooks(ctx context.Context, value runtime.Value, err error) error {
	for _, hook := range enc.post {
		if hookErr := hook(ctx, value, err); hookErr != nil {
			return hookErr
		}
	}

	return nil
}
