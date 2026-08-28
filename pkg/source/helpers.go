package source

func SkipWhitespaceForward(content string, offset int) int {
	for offset < len(content) {
		ch := content[offset]

		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			break
		}

		offset++
	}

	return offset
}

func SkipHorizontalWhitespaceForward(content string, offset int) int {
	for offset < len(content) {
		ch := content[offset]

		// Skip spaces and tabs only; do NOT cross line breaks
		if ch != ' ' && ch != '\t' {
			break
		}

		offset++
	}

	return offset
}

func computeVisualOffset(line string, column int) int {
	runes := []rune(line)
	offset := 0
	tabWidth := 4

	for i := 0; i < column-1 && i < len(runes); i++ {
		if runes[i] == '\t' {
			offset += tabWidth - (offset % tabWidth)
		} else {
			offset += 1
		}
	}

	return offset
}
