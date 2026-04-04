package internal

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CompileStringLiteral processes a string literal from the FQL AST and converts it into a runtime string.
// It handles escape sequences like \n and \t, and properly extracts the string content without quotes.
// Parameters:
//   - ctx: The string literal context from the AST
//
// Returns:
//   - An operand representing the compiled string constant
func parseStringLiteral(ctx fql.IStringLiteralContext) runtime.String {
	if ctx == nil || ctx.StringLiteral() == nil {
		return runtime.EmptyString
	}

	var b strings.Builder

	// Process each child node in the string literal
	for _, child := range ctx.GetChildren() {
		tree := child.(antlr.TerminalNode)
		sym := tree.GetSymbol()
		input := sym.GetInputStream()

		if input == nil {
			continue
		}

		size := input.Size()
		// Skip the opening and closing quotes
		start := sym.GetStart() + 1
		stop := sym.GetStop() - 1

		// Ensure we don't go beyond the input size
		if stop >= size {
			stop = size - 1
		}

		if start < size && stop < size {
			// Process each character in the string
			for i := start; i <= stop; i++ {
				ch := input.GetText(i, i)

				switch ch {
				case "\\":
					// Handle escape sequences
					c2 := input.GetText(i, i+1)

					switch c2 {
					case "\\n":
						b.WriteString("\n")
					case "\\t":
						b.WriteString("\t")
					default:
						b.WriteString(c2)
					}

					// Skip the next character as it's part of the escape sequence
					i++
				default:
					// Add regular characters as-is
					b.WriteString(ch)
				}
			}
		}
	}

	return runtime.NewString(b.String())
}

func parseTemplateChunk(text string) string {
	if text == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(text))

	for i := 0; i < len(text); i++ {
		ch := text[i]
		if ch == '\\' && i+1 < len(text) {
			next := text[i+1]
			switch next {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '`':
				b.WriteByte('`')
			case '$':
				b.WriteByte('$')
			case '\\':
				b.WriteByte('\\')
			default:
				b.WriteByte('\\')
				b.WriteByte(next)
			}
			i++
			continue
		}
		b.WriteByte(ch)
	}

	return b.String()
}

func parseTemplateLiteralConst(ctx fql.ITemplateLiteralContext) (runtime.String, bool) {
	if ctx == nil {
		return runtime.EmptyString, false
	}

	var b strings.Builder

	for _, el := range ctx.AllTemplateElement() {
		if el == nil {
			continue
		}
		if expr := el.Expression(); expr != nil {
			if val, ok := constStringFromExpression(expr); ok {
				b.WriteString(val)
				continue
			}
			return runtime.EmptyString, false
		}
		if chunk := el.TemplateChars(); chunk != nil {
			b.WriteString(parseTemplateChunk(chunk.GetText()))
		}
	}

	return runtime.NewString(b.String()), true
}

func parseStringLiteralConst(ctx fql.IStringLiteralContext) (runtime.String, bool) {
	if ctx == nil {
		return runtime.EmptyString, false
	}

	if ctx.StringLiteral() != nil {
		return parseStringLiteral(ctx), true
	}

	if tmpl := ctx.TemplateLiteral(); tmpl != nil {
		return parseTemplateLiteralConst(tmpl)
	}

	return runtime.EmptyString, false
}

func constStringFromExpression(expr fql.IExpressionContext) (string, bool) {
	val, ok := tryConstConcatStringFromExpression(expr)
	if !ok {
		return "", false
	}

	return val.String(), true
}

func invalidNumericLiteralDetails(kind string, err error) (string, string) {
	switch kind {
	case "integer":
		if errors.Is(err, strconv.ErrRange) {
			return "Integer literal is out of range", "Use an integer value that fits within the supported range."
		}

		return "Invalid integer literal", "Use a valid integer value, e.g. 42."
	case "float":
		if errors.Is(err, strconv.ErrRange) {
			return "Float literal is out of range", "Use a finite float value within the supported range."
		}

		return "Invalid float literal", "Use a valid float value, e.g. 1.5."
	default:
		return "Invalid numeric literal", "Use a valid numeric value."
	}
}

func literalValueOf(ctx fql.ILiteralContext) (runtime.Value, bool) {
	if ctx == nil {
		return nil, false
	}

	switch {
	case ctx.NoneLiteral() != nil:
		return runtime.None, true
	case ctx.StringLiteral() != nil:
		return parseStringLiteralConst(ctx.StringLiteral())
	case ctx.IntegerLiteral() != nil:
		val, err := strconv.Atoi(ctx.IntegerLiteral().GetText())
		if err != nil {
			return nil, false
		}
		return runtime.NewInt(val), true
	case ctx.FloatLiteral() != nil:
		val, err := strconv.ParseFloat(ctx.FloatLiteral().GetText(), 64)
		if err != nil {
			return nil, false
		}
		return runtime.NewFloat(val), true
	case ctx.BooleanLiteral() != nil:
		switch strings.ToLower(ctx.BooleanLiteral().GetText()) {
		case "true":
			return runtime.True, true
		case "false":
			return runtime.False, true
		}
	case ctx.ArrayLiteral() != nil:
		return runtime.NewArray(0), true
	case ctx.ObjectLiteral() != nil:
		return runtime.NewObject(), true
	}

	return nil, false
}
