package debugger

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func formatValue(value runtime.Value, description vm.DebugValueDescription, maxBytes int) Value {
	typeName := formatTypeName(description, maxBytes)
	if description.Hint.Display != "" {
		return Value{Type: typeName, Display: boundedText(description.Hint.Display, maxBytes)}
	}

	var display string
	switch v := value.(type) {
	case runtime.String:
		display, _ = renderText(string(v), maxBytes, true)

		return Value{Type: typeName, Display: display}
	case runtime.Boolean:
		display = strconv.FormatBool(bool(v))
	case runtime.Int:
		display = strconv.FormatInt(int64(v), 10)
	case runtime.Float:
		display = strconv.FormatFloat(float64(v), 'g', -1, 64)
	case runtime.Duration:
		display = v.String()
	case runtime.DateTime:
		display = v.UTC().Format(time.RFC3339Nano)
	default:
		switch {
		case value == nil || value == runtime.None:
			display = "NONE"
		case description.Length >= 0:
			display = typeName + "(" + strconv.Itoa(description.Length) + ")"
		default:
			display = "<" + typeName + ">"
		}
	}

	return Value{Type: typeName, Display: boundedText(display, maxBytes)}
}

func formatTypeName(description vm.DebugValueDescription, maxBytes int) string {
	if hint := description.Hint.TypeName; hint != "" {
		if name, valid := renderText(hint, maxBytes, false); valid {
			return name
		}
	}

	return boundedText(description.TypeName, maxBytes)
}

func boundedText(value string, maxBytes int) string {
	text, _ := renderText(value, maxBytes, false)

	return text
}

// renderText counts output bytes, including delimiters and the truncation
// marker. It remembers the last whole token that leaves room for the marker,
// so an exact fit can still use the entire budget.
func renderText(value string, maxBytes int, quoted bool) (string, bool) {
	const marker = "..."
	if maxBytes <= 0 {
		return "", true
	}

	if !quoted && len(value) <= maxBytes && utf8.ValidString(value) {
		return value, true
	}

	var scratch [128]byte
	out := scratch[:0]
	limit := maxBytes
	if quoted {
		if maxBytes < 2 {
			return marker[:maxBytes], true
		}

		out = append(out, '"')
		limit-- // closing quote
	}

	markerAt := len(out)
	valid := true
	markerValid := true
	for len(value) > 0 {
		r, size := utf8.DecodeRuneInString(value)
		unit := value[:size]
		if r == utf8.RuneError && size == 1 {
			valid = false
			unit = "\uFFFD"
		}

		var escape [12]byte
		if quoted && (r < ' ' || r >= 0x7f || r == '"' || r == '\\') {
			token := strconv.AppendQuote(escape[:0], unit)
			unit = string(token[1 : len(token)-1])
		}

		if len(unit) > limit-len(out) {
			if maxBytes < len(marker) || (quoted && maxBytes < len(marker)+2) {
				return marker[:min(maxBytes, len(marker))], true
			}

			valid = markerValid
			out = append(out[:markerAt], marker...)
			break
		}

		out = append(out, unit...)
		if len(out) <= limit-len(marker) {
			markerAt = len(out)
			markerValid = valid
		}

		value = value[size:]
	}

	if quoted {
		out = append(out, '"')
	}

	return string(out), valid
}
