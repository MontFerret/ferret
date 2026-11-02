package internal

import (
	"strings"
	"unicode/utf8"
)

func formatStringLiteral(raw string, singleQuote bool) string {
	quote := byte('"')
	if singleQuote {
		quote = '\''
	}

	delimiter, width := utf8.DecodeRuneInString(raw)
	if delimiter == rune(quote) {
		return raw
	}

	if delimiter != '\'' && delimiter != '"' && delimiter != '´' {
		return raw
	}

	if len(raw) < 2*width || !strings.HasSuffix(raw, raw[:width]) {
		return raw
	}

	body := raw[width : len(raw)-width]
	// FQL preserves doubled delimiters and unknown backslash pairs. Escaping a
	// preferred quote would therefore change the value rather than encode it.
	if strings.IndexByte(body, quote) >= 0 {
		return raw
	}

	trailingBackslashes := 0
	for i := len(body) - 1; i >= 0 && body[i] == '\\'; i-- {
		trailingBackslashes++
	}

	// Acute-delimited strings can end in an unpaired backslash. The compiler
	// includes the closing acute accent in that escape pair; retain the token.
	if trailingBackslashes%2 != 0 {
		return raw
	}

	return string(quote) + body + string(quote)
}
