package runtime

import (
	"encoding/base64"
	"sort"
	"strconv"
	"strings"
	"time"
)

const jsonHex = "0123456789abcdef"

func writeString(buf *strings.Builder, s string) {
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

func writeValue(buf *strings.Builder, value Value) {
	if value == nil || value == None {
		buf.WriteString("null")

		return
	}

	switch v := value.(type) {
	case Boolean:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case Int:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case Float:
		buf.WriteString(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case String:
		writeString(buf, string(v))
	case Binary:
		writeString(buf, base64.StdEncoding.EncodeToString([]byte(v)))
	case DateTime:
		writeString(buf, v.Time.Format(time.RFC3339Nano))
	case *Array:
		writeArray(buf, v)
	case *Object:
		writeObject(buf, v)
	default:
		buf.WriteString(v.String())
	}
}

func writeArray(buf *strings.Builder, arr *Array) {
	buf.WriteByte('[')

	for i, el := range arr.data {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeValue(buf, el)
	}

	buf.WriteByte(']')
}

func writeObject(buf *strings.Builder, obj *Object) {
	buf.WriteByte('{')

	if len(obj.data) == 0 {
		buf.WriteByte('}')
		return
	}

	keys := make([]string, 0, len(obj.data))
	for k := range obj.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		writeString(buf, k)
		buf.WriteByte(':')
		writeValue(buf, obj.data[k])
	}
	buf.WriteByte('}')
}
