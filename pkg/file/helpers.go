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
