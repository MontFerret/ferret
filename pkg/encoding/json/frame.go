package json

import (
	"bytes"
	"context"
	"encoding/base64"
	"sort"
	"strconv"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/goccy/go-json"
)

type frameKind uint8

const (
	frameObject frameKind = iota + 1
	frameArray
)

type decodeFrame struct {
	kind      frameKind
	obj       *runtime.Object
	arr       *runtime.Array
	expectKey bool
	key       string
}

const jsonHex = "0123456789abcdef"

func writeJSONString(buf *bytes.Buffer, s string) {
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

func encodeValue(ctx context.Context, buf *bytes.Buffer, value runtime.Value) error {
	if value == nil || value == runtime.None {
		buf.WriteString("null")

		return nil
	}

	if marshaler, ok := value.(json.Marshaler); ok {
		out, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}

		buf.Write(out)

		return nil
	}

	switch v := value.(type) {
	case runtime.Boolean:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

		return nil
	case runtime.Int:
		buf.WriteString(strconv.FormatInt(int64(v), 10))

		return nil
	case runtime.Float:
		buf.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 64))

		return nil
	case runtime.String:
		writeJSONString(buf, string(v))

		return nil
	case runtime.Binary:
		writeJSONString(buf, base64.StdEncoding.EncodeToString([]byte(v)))
		return nil
	case runtime.DateTime:
		writeJSONString(buf, v.Time.Format(time.RFC3339Nano))

		return nil
	case runtime.Query:
		buf.WriteByte('{')
		writeJSONString(buf, "kind")
		buf.WriteByte(':')
		writeJSONString(buf, v.Kind.String())
		buf.WriteByte(',')
		writeJSONString(buf, "payload")
		buf.WriteByte(':')
		writeJSONString(buf, v.Payload.String())
		buf.WriteByte(',')
		writeJSONString(buf, "params")
		buf.WriteByte(':')
		if err := encodeValue(ctx, buf, v.Params); err != nil {
			return err
		}
		buf.WriteByte('}')

		return nil
	case runtime.Map:
		return encodeMap(ctx, buf, v)
	case runtime.List:
		return encodeList(ctx, buf, v)
	case runtime.Iterable:
		list, err := runtime.ToList(ctx, value)
		if err != nil {
			return err
		}
		return encodeList(ctx, buf, list)
	case runtime.Unwrappable:
		return encodeAny(buf, v.Unwrap())
	default:
		return encodeAny(buf, v)
	}
}

func encodeAny(buf *bytes.Buffer, value any) error {
	switch v := value.(type) {
	case json.RawMessage:
		buf.Write(v)

		return nil
	case string:
		writeJSONString(buf, v)

		return nil
	case []byte:
		writeJSONString(buf, base64.StdEncoding.EncodeToString(v))

		return nil
	case time.Time:
		writeJSONString(buf, v.Format(time.RFC3339Nano))

		return nil
	}

	out, err := json.Marshal(value)
	if err != nil {
		return err
	}

	buf.Write(out)

	return nil
}

func encodeMap(ctx context.Context, buf *bytes.Buffer, value runtime.Map) error {
	type entry struct {
		key string
		val runtime.Value
	}

	entries := make([]entry, 0)
	err := value.ForEach(ctx, func(_ context.Context, val, key runtime.Value) (runtime.Boolean, error) {
		entries = append(entries, entry{key: key.String(), val: val})

		return true, nil
	})

	if err != nil {
		return err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})

	buf.WriteString("{")

	for i, item := range entries {
		if i > 0 {
			buf.WriteString(",")
		}

		writeJSONString(buf, item.key)

		buf.WriteString(":")

		if err := encodeValue(ctx, buf, item.val); err != nil {
			return err
		}
	}

	buf.WriteString("}")

	return nil
}

func encodeList(ctx context.Context, buf *bytes.Buffer, value runtime.List) error {
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

		if err := encodeValue(ctx, buf, item); err != nil {
			return err
		}
	}

	buf.WriteString("]")

	return nil
}
