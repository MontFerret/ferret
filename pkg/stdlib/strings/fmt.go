package strings

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"
)

// FMT formats the template using these arguments.
// @param {String} template - template.
// @param {Any, repeated} args - template arguments.
// @return {String} - string formed by template using arguments.
func Fmt(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)
	if err != nil {
		return values.None, err
	}

	formatted, err := format(args[0].String(), args[1:])
	if err != nil {
		return values.None, err
	}

	return values.NewString(formatted), nil
}

func format(template string, args []core.Value) (string, error) {
	rgx := regexp.MustCompile("{[0-9]*}")

	argsCount := len(args)
	emptyBracketsCount := strings.Count(template, "{}")

	if argsCount > emptyBracketsCount && emptyBracketsCount != 0 {
		return "", errors.Errorf("there are arguments that have never been used")
	}

	var betweenBrackets string
	var n int
	// index of the last value
	// inserted into the template
	var lastArgIdx int
	var err error

	template = rgx.ReplaceAllStringFunc(template, func(s string) string {
		if err != nil {
			return ""
		}

		betweenBrackets = s[1 : len(s)-1]

		if betweenBrackets == "" {
			if argsCount <= lastArgIdx {
				err = errors.Errorf("not enough arguments")
				return ""
			}

			lastArgIdx++
			return args[lastArgIdx-1].String()
		}

		n, err = strconv.Atoi(betweenBrackets)
		if err != nil {
			err = errors.Errorf("failed to parse int: %v", err)
			return ""
		}

		if n >= argsCount {
			err = errors.Errorf("invalid reference to argument `%d`", n)
			return ""
		}

		return args[n].String()
	})

	return template, err
}
