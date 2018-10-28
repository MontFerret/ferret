package strings

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Fmt(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.StringType)
	if err != nil {
		return values.None, err
	}

	formatted, err := format(args[0].String(), args[1:]...)
	if err != nil {
		return values.None, err
	}

	return values.NewString(formatted), nil
}

func format(template string, args ...core.Value) (string, error) {
	argsMap := map[int]string{}
	maxArgs := len(args)

	template = replaceEmptyBrackets(template, strings.Count(template, "{}"))

	for idx, arg := range args {
		argsMap[idx] = arg.String()
	}

	rgx, err := regexp.Compile("{[0-9]+}")
	if err != nil {
		return "", fmt.Errorf("failed to build regexp: %v", err)
	}

	n := 0

	template = rgx.ReplaceAllStringFunc(template, func(s string) string {
		if err != nil {
			return ""
		}

		n, err = strconv.Atoi(s[1 : len(s)-1])
		if err != nil {
			err = fmt.Errorf("failed to parse int: %v", err)
			return ""
		}

		if n > maxArgs {
			err = fmt.Errorf("there is no arg {%d}", n)
			return ""
		}

		argm, ok := argsMap[n]
		if !ok {
			err = fmt.Errorf("value {%d} doesn't exists", n)
		}

		return argm
	})

	return template, err
}

func replaceEmptyBrackets(s string, n int) string {
	for i := 0; i < n; i++ {
		s = strings.Replace(s, "{}", "{"+strconv.Itoa(i)+"}", 1)
	}
	return s
}
