package json

import "unicode/utf8"

// jsonASCIIEnd locates the end of an ASCII run that needs no JSON escaping.
func jsonASCIIEnd(text string, start, end int) int {
	for start < end {
		value := text[start]
		if value < 0x20 || value >= utf8.RuneSelf || value == '"' || value == '\\' {
			break
		}

		start++
	}

	return start
}
