package file

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
