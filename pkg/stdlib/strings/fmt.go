package strings

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"

	"errors"
)

// FMT formats the template using these arguments.
// @param {String} template - template.
// @param {Any, repeated} args - template arguments.
// @return {String} - string formed by template using arguments.
func Fmt(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	arg0, err := runtime.CastString(args[0])

	if err != nil {
		return runtime.None, err
	}

	formatted, err := format(arg0.String(), args[1:])
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(formatted), nil
}

func format(template string, args []runtime.Value) (string, error) {
	rgx := regexp.MustCompile("{[0-9]*}")

	argsCount := len(args)
	emptyBracketsCount := strings.Count(template, "{}")

	if argsCount > emptyBracketsCount && emptyBracketsCount != 0 {
		return "", errors.New("there are arguments that have never been used")
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
				err = errors.New("not enough arguments")
				return ""
			}

			lastArgIdx++
			return args[lastArgIdx-1].String()
		}

		n, err = strconv.Atoi(betweenBrackets)
		if err != nil {
			err = fmt.Errorf("failed to parse int: %v", err)
			return ""
		}

		if n >= argsCount {
			err = fmt.Errorf("invalid reference to argument `%d`", n)
			return ""
		}

		return args[n].String()
	})

	return template, err
}
