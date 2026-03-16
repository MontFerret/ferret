package json

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const jsonHex = "0123456789abcdef"

type encoder struct {
	pre  []encoding.PreEncoderHook
	post []encoding.PostEncoderHook
}

func (enc encoder) encodeValue(ctx context.Context, buf *bytes.Buffer, value runtime.Value) error {
	if err := enc.runPreHooks(value); err != nil {
		return err
	}

	if value == nil || value == runtime.None {
		buf.WriteString("null")

		return enc.runPostHooks(runtime.None, nil)
	}

	if marshaler, ok := value.(json.Marshaler); ok {
		out, err := marshaler.MarshalJSON()

		if hookErr := enc.runPostHooks(value, err); hookErr != nil {
			return errors.Join(err, hookErr)
		}

		if err != nil {
			return err
		}

		buf.Write(out)

		return nil
	}

	var err error

	switch v := value.(type) {
	case runtime.Boolean:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case runtime.Int:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case runtime.Float:
		buf.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case runtime.String:
		enc.writeJSONString(buf, string(v))
	case runtime.Binary:
		enc.writeJSONString(buf, base64.StdEncoding.EncodeToString(v))
	case runtime.DateTime:
		enc.writeJSONString(buf, v.Time.Format(time.RFC3339Nano))
	case runtime.Map:
		err = enc.encodeMap(ctx, buf, v)
	case runtime.List:
		err = enc.encodeList(ctx, buf, v)
	case runtime.Iterable:
		list, e := runtime.ToList(ctx, value)
		if e != nil {
			err = e
			break
		}

		err = enc.encodeList(ctx, buf, list)
	case runtime.Unwrappable:
		err = enc.encodeAny(buf, v.Unwrap())
	default:
		err = enc.encodeAny(buf, v)
	}

	if hookErr := enc.runPostHooks(value, err); hookErr != nil {
		return errors.Join(err, hookErr)
	}

	return err
}

func (enc encoder) encodeAny(buf *bytes.Buffer, value any) error {
	switch v := value.(type) {
	case json.RawMessage:
		buf.Write(v)

		return nil
	case string:
		enc.writeJSONString(buf, v)

		return nil
	case []byte:
		enc.writeJSONString(buf, base64.StdEncoding.EncodeToString(v))

		return nil
	case time.Time:
		enc.writeJSONString(buf, v.Format(time.RFC3339Nano))

		return nil
	}

	out, err := json.Marshal(value)
	if err != nil {
		return err
	}

	buf.Write(out)

	return nil
}

func (enc encoder) encodeMap(ctx context.Context, buf *bytes.Buffer, value runtime.Map) error {
	buf.WriteString("{")
	first := true

	err := value.ForEach(ctx, func(_ context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		if !first {
			buf.WriteString(",")
		} else {
			first = false
		}

		enc.writeJSONString(buf, key.String())
		buf.WriteString(":")

		if err := enc.encodeValue(ctx, buf, val); err != nil {
			return false, err
		}
		return true, nil
	})

	if err != nil {
		return err
	}

	buf.WriteString("}")

	return nil
}

func (enc encoder) encodeList(ctx context.Context, buf *bytes.Buffer, value runtime.List) error {
	length, err := value.Length(ctx)
	if err != nil {
		return err
	}

	buf.WriteString("[")

	for i := runtime.Int(0); i < length; i++ {
		if i > 0 {
			buf.WriteString(",")
		}

		item, err := value.At(ctx, i)

		if err != nil {
			return err
		}

		if err := enc.encodeValue(ctx, buf, item); err != nil {
			return err
		}
	}

	buf.WriteString("]")

	return nil
}

func (enc encoder) writeJSONString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\', '"':
			buf.WriteByte('\\')
			buf.WriteRune(r)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				buf.WriteString(`\u00`)
				buf.WriteByte(jsonHex[r>>4])
				buf.WriteByte(jsonHex[r&0x0f])
				continue
			}
			buf.WriteRune(r)
		}
	}
	buf.WriteByte('"')
}

func (enc encoder) runPreHooks(value runtime.Value) error {
	if len(enc.pre) == 0 {
		return nil
	}

	for _, hook := range enc.pre {
		if err := hook(value); err != nil {
			return err
		}
	}

	return nil
}

func (enc encoder) runPostHooks(value runtime.Value, err error) error {
	if len(enc.post) == 0 {
		return nil
	}

	for _, hook := range enc.post {
		if err := hook(value, err); err != nil {
			return err
		}
	}

	return nil
}
