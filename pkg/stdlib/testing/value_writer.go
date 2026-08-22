package testing

import "strings"

type boundedValueWriter struct {
	builder   strings.Builder
	limit     int
	size      int
	truncated bool
}

func newBoundedValueWriter(limit int) *boundedValueWriter {
	if limit < 0 {
		limit = 0
	}

	return &boundedValueWriter{limit: limit}
}

func (w *boundedValueWriter) writeRune(value rune) {
	if w.truncated {
		return
	}

	if w.size >= w.limit {
		w.truncated = true

		return
	}

	w.builder.WriteRune(value)
	w.size++
}

func (w *boundedValueWriter) writeString(value string) {
	for _, char := range value {
		w.writeRune(char)
		if w.truncated {
			return
		}
	}
}

func (w *boundedValueWriter) remaining() int {
	return w.limit - w.size
}

func (w *boundedValueWriter) String() string {
	return w.builder.String()
}
