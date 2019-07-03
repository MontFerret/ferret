package common

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/gorilla/css/scanner"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func DeserializeStyles(input values.String) (*values.Object, error) {
	styles := values.NewObject()

	if input == values.EmptyString {
		return styles, nil
	}

	s := scanner.New(input.String())

	var name string
	var value bytes.Buffer
	var setValue = func() {
		styles.Set(values.NewString(strings.TrimSpace(name)), values.NewString(strings.TrimSpace(value.String())))
		name = ""
		value.Reset()
	}

	for {
		token := s.Next()

		if token == nil {
			break
		}

		if token.Type == scanner.TokenEOF {
			break
		}

		if name == "" && token.Type == scanner.TokenIdent {
			name = token.Value

			// skip : and white spaces
			for {
				token = s.Next()

				if token.Value != ":" && token.Type != scanner.TokenS {
					break
				}
			}
		}

		switch token.Type {
		case scanner.TokenChar:
			// end of style declaration
			if token.Value == ";" {
				if name != "" {
					setValue()
				}
			} else {
				value.WriteString(token.Value)
			}
		case scanner.TokenNumber:
			num, err := strconv.ParseFloat(token.Value, 64)

			if err == nil {
				styles.Set(values.NewString(name), values.NewFloat(num))
				// reset prop
				name = ""
				value.Reset()
			}
		default:
			value.WriteString(token.Value)
		}
	}

	if name != "" && value.Len() > 0 {
		setValue()
	}

	return styles, nil
}

func SerializeStyles(_ context.Context, styles *values.Object) values.String {
	if styles == nil {
		return values.EmptyString
	}

	var b bytes.Buffer

	styles.ForEach(func(value core.Value, key string) bool {
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(value.String())
		b.WriteString("; ")

		return true
	})

	return values.NewString(b.String())
}
