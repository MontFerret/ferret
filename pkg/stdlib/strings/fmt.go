package strings

import (
	"context"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// fmt substitutes values using either automatic {} placeholders or zero-based {n} indexes. Modes cannot be mixed. Repeated indexes are allowed. Escape literal braces with {{ and }}. Every supplied value must be referenced.
// @param template {String} The format string. Indexes contain decimal digits only.
// @param values {Any, repeated} Values substituted using their string representations.
// @return {String} The formatted string. Malformed braces, missing values, and unused values are errors.
func Fmt(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	template, err := runtime.CastArg[runtime.String](args[0], 0)
	if err != nil {
		return runtime.None, err
	}

	values := args[1:]
	used := make([]bool, len(values))
	var out strings.Builder
	var automatic, indexed bool
	next := 0
	for pos := 0; pos < len(template); {
		ch := template[pos]
		if ch != '{' && ch != '}' {
			start := pos
			for pos < len(template) && template[pos] != '{' && template[pos] != '}' {
				pos++
			}

			out.WriteString(string(template[start:pos]))

			continue
		}

		if pos+1 < len(template) && template[pos+1] == ch {
			out.WriteByte(ch)
			pos += 2

			continue
		}

		if ch == '}' {
			return runtime.None, formatError("unmatched closing brace")
		}

		end := pos + 1
		for end < len(template) && template[end] >= '0' && template[end] <= '9' {
			end++
		}

		if end == len(template) || template[end] != '}' {
			return runtime.None, formatError("malformed placeholder")
		}

		index := next
		if end == pos+1 {
			automatic = true
			next++
		} else {
			indexed = true
			index, err = strconv.Atoi(string(template[pos+1 : end]))
			if err != nil {
				return runtime.None, formatError("placeholder index is out of range")
			}
		}

		if automatic && indexed {
			return runtime.None, formatError("automatic and indexed placeholders cannot be mixed")
		}

		if index >= len(values) {
			return runtime.None, formatError("placeholder has no corresponding value")
		}

		used[index] = true
		out.WriteString(values[index].String())
		pos = end + 1
	}

	for index, referenced := range used {
		if !referenced {
			return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "unused formatting value"), index+1)
		}
	}

	return runtime.String(out.String()), nil
}

func formatError(message string) error {
	return runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, message), 0)
}
