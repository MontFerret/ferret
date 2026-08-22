package testing

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	// maxFormattedValueRunes bounds each expected/actual field independently.
	maxFormattedValueRunes = 256
	maxRenderedBinaryBytes = 32
	// Large objects are summarized before key collection so a bounded message
	// does not require sorting an unbounded property set.
	maxRenderedObjectProperties = 32
	truncatedValueMarker        = "... <truncated>"
)

func formatValue(ctx context.Context, value runtime.Value) string {
	typeName := runtime.TypeOf(value).String()
	prefix := typeName + " '"
	suffix := "'"
	contentLimit := maxFormattedValueRunes - utf8.RuneCountInString(prefix) - utf8.RuneCountInString(suffix)
	writer := newBoundedValueWriter(contentLimit)

	writeTopLevelValue(ctx, writer, value)

	content := writer.String()
	if writer.truncated {
		available := contentLimit - utf8.RuneCountInString(truncatedValueMarker)
		content = truncateRunes(content, available) + truncatedValueMarker
	}

	formatted := prefix + content + suffix
	if utf8.RuneCountInString(formatted) <= maxFormattedValueRunes {
		return formatted
	}

	return truncateRunes(formatted, maxFormattedValueRunes-utf8.RuneCountInString(truncatedValueMarker)) + truncatedValueMarker
}

func writeTopLevelValue(ctx context.Context, writer *boundedValueWriter, value runtime.Value) {
	if value == nil || value == runtime.None {
		writeEscapedRepresentation(writer, "none")

		return
	}

	switch typed := value.(type) {
	case *runtime.Array:
		writeArrayValue(ctx, writer, typed)
	case runtime.ObjectLike:
		writeObjectValue(ctx, writer, typed)
	case runtime.Binary:
		if len(typed) > maxRenderedBinaryBytes {
			writeEscapedRepresentation(writer, fmt.Sprintf("<%d bytes; %s>", len(typed), truncatedValueMarker))
		} else {
			writeEscapedRepresentation(writer, string(typed))
		}
	case runtime.String:
		writeEscapedRepresentation(writer, string(typed))
	default:
		writeEscapedRepresentation(writer, value.String())
	}
}

func writeSerializedValue(ctx context.Context, writer *boundedValueWriter, value runtime.Value) {
	if value == nil || value == runtime.None {
		writeEscapedRepresentation(writer, "null")

		return
	}

	switch typed := value.(type) {
	case runtime.Boolean, runtime.Int, runtime.Float:
		writeEscapedRepresentation(writer, typed.String())
	case runtime.Duration:
		writeJSONString(writer, typed.String())
	case runtime.String:
		writeJSONString(writer, string(typed))
	case runtime.Binary:
		writeBinaryJSON(writer, typed)
	case runtime.DateTime:
		writeJSONString(writer, typed.Time.Format(time.RFC3339Nano))
	case *runtime.Array:
		writeArrayValue(ctx, writer, typed)
	case runtime.ObjectLike:
		writeObjectValue(ctx, writer, typed)
	default:
		writeEscapedRepresentation(writer, value.String())
	}
}

func writeArrayValue(ctx context.Context, writer *boundedValueWriter, array *runtime.Array) {
	length, err := array.Length(ctx)
	if err != nil {
		writeEscapedRepresentation(writer, array.String())

		return
	}

	writeEscapedRepresentation(writer, "[")
	for index := runtime.ZeroInt; index < length && !writer.truncated; index++ {
		if index > 0 {
			writeEscapedRepresentation(writer, ",")
		}

		value, found, err := array.LookupAt(ctx, index)
		if err != nil || !found {
			writer.truncated = true

			return
		}

		writeSerializedValue(ctx, writer, value)
	}

	writeEscapedRepresentation(writer, "]")
}

func writeObjectValue(ctx context.Context, writer *boundedValueWriter, object runtime.ObjectLike) {
	length, err := object.Length(ctx)
	if err != nil {
		writeEscapedRepresentation(writer, object.String())

		return
	}

	if length > maxRenderedObjectProperties {
		writeEscapedRepresentation(writer, fmt.Sprintf("{%d properties; %s}", length, truncatedValueMarker))

		return
	}

	values := make(map[string]runtime.Value, length)
	keys := make([]string, 0, length)
	err = object.ForEach(ctx, func(_ context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		stringKey, ok := key.(runtime.String)
		if !ok {
			return runtime.False, runtime.TypeErrorOf(key, runtime.TypeString)
		}

		keyText := string(stringKey)
		keys = append(keys, keyText)
		values[keyText] = value

		return runtime.True, nil
	})
	if err != nil {
		writeEscapedRepresentation(writer, object.String())

		return
	}

	sort.Strings(keys)
	writeEscapedRepresentation(writer, "{")
	for index, key := range keys {
		if index > 0 {
			writeEscapedRepresentation(writer, ",")
		}

		writeJSONString(writer, key)
		writeEscapedRepresentation(writer, ":")
		writeSerializedValue(ctx, writer, values[key])
		if writer.truncated {
			return
		}
	}

	writeEscapedRepresentation(writer, "}")
}

func writeBinaryJSON(writer *boundedValueWriter, value runtime.Binary) {
	writeEscapedRepresentation(writer, `"`)
	if writer.truncated {
		return
	}

	remaining := writer.remaining()
	byteLimit := (remaining / 4) * 3
	if byteLimit > len(value) {
		byteLimit = len(value)
	}

	encoded := base64.StdEncoding.EncodeToString(value[:byteLimit])
	writeEscapedRepresentation(writer, encoded)
	if byteLimit < len(value) {
		writer.truncated = true

		return
	}

	writeEscapedRepresentation(writer, `"`)
}

func writeJSONString(writer *boundedValueWriter, value string) {
	writeEscapedRepresentation(writer, `"`)
	for _, char := range value {
		if writer.truncated {
			return
		}

		switch char {
		case '\\', '"':
			writeEscapedRepresentation(writer, `\`)
			writeEscapedRune(writer, char)
		case '\b':
			writeEscapedRepresentation(writer, `\b`)
		case '\f':
			writeEscapedRepresentation(writer, `\f`)
		case '\n':
			writeEscapedRepresentation(writer, `\n`)
		case '\r':
			writeEscapedRepresentation(writer, `\r`)
		case '\t':
			writeEscapedRepresentation(writer, `\t`)
		default:
			if char < 0x20 {
				writeEscapedRepresentation(writer, fmt.Sprintf(`\u%04x`, char))

				continue
			}

			writeEscapedRune(writer, char)
		}
	}

	writeEscapedRepresentation(writer, `"`)
}

func writeEscapedRepresentation(writer *boundedValueWriter, value string) {
	for _, char := range value {
		writeEscapedRune(writer, char)
		if writer.truncated {
			return
		}
	}
}

func writeEscapedRune(writer *boundedValueWriter, value rune) {
	switch value {
	case '\\':
		writer.writeString(`\\`)
	case '\'':
		writer.writeString(`\'`)
	default:
		writer.writeRune(value)
	}
}

func truncateRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
	}

	if utf8.RuneCountInString(value) <= limit {
		return value
	}

	var builder strings.Builder
	for _, char := range value {
		if limit == 0 {
			break
		}

		builder.WriteRune(char)
		limit--
	}

	return builder.String()
}
