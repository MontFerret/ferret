package formatter

import "strings"

func applyCase(mode CaseMode, val string) string {
	switch mode {
	case CaseModeUpper:
		return strings.ToUpper(val)
	case CaseModeLower:
		return strings.ToLower(val)
	default:
		return val
	}
}

func unquoteStringLiteral(raw string) string {
	if len(raw) < 2 {
		return raw
	}
	quote := raw[0]
	content := raw[1 : len(raw)-1]
	var b strings.Builder
	b.Grow(len(content))

	for i := 0; i < len(content); i++ {
		ch := content[i]
		if ch == '\\' && i+1 < len(content) {
			next := content[i+1]
			switch next {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			default:
				b.WriteByte('\\')
				b.WriteByte(next)
			}
			i++
			continue
		}
		if (quote == '\'' || quote == '"') && ch == quote && i+1 < len(content) && content[i+1] == quote {
			b.WriteByte(quote)
			i++
			continue
		}
		b.WriteByte(ch)
	}

	return b.String()
}

func quoteStringLiteral(val string, singleQuote bool) string {
	quote := byte('"')
	doubleQuote := "\"\""
	if singleQuote {
		quote = '\''
		doubleQuote = "''"
	}
	var b strings.Builder
	b.Grow(len(val) + 2)
	b.WriteByte(quote)
	for _, r := range val {
		switch r {
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		case '\\':
			b.WriteString("\\\\")
		default:
			if r == rune(quote) {
				b.WriteString(doubleQuote)
				continue
			}
			b.WriteRune(r)
		}
	}
	b.WriteByte(quote)
	return b.String()
}

func formatStringLiteral(raw string, singleQuote bool) string {
	return quoteStringLiteral(unquoteStringLiteral(raw), singleQuote)
}
