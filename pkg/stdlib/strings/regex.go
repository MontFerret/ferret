package strings

import (
	"context"
	"regexp"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// REGEX_MATCH returns the matches in the given string text, using the regex.
// @param {String} str - The string to search in.
// @param {String} expression - A regular expression to use for matching the text.
// @param {Boolean} caseInsensitive - If set to true, the matching will be case-insensitive. The default is false.
// @return {Any[]} - An array of strings containing the matches.
func RegexMatch(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	text := args[0].String()
	exp := args[1].String()

	if len(args) > 2 {
		if args[2] == values.True {
			exp = "(?i)" + exp
		}
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return values.None, err
	}

	matches := reg.FindAllStringSubmatch(text, -1)
	res := values.NewArray(10)

	if len(matches) == 0 {
		return res, nil
	}

	for _, m := range matches[0] {
		res.Push(values.NewString(m))
	}

	return res, nil
}

// REGEX_SPLIT splits the given string text into a list of strings, using the separator.
// @param {String} str - The string to split.
// @param {String} expression - A regular expression to use for splitting the text.
// @param {Boolean} caseInsensitive - If set to true, the matching will be case-insensitive. The default is false.
// @param {Int} limit - Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @return {Any[]} - An array of strings splitted by the expression.
func RegexSplit(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	text := args[0].String()
	exp := args[1].String()
	limit := -1

	if len(args) > 2 {
		if args[2].Type() == types.Int {
			limit = int(args[2].(values.Int))
		}
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return values.None, err
	}

	matches := reg.Split(text, limit)
	res := values.NewArray(10)

	if len(matches) == 0 {
		return res, nil
	}

	for _, m := range matches {
		res.Push(values.NewString(m))
	}

	return res, nil
}

// REGEX_TEST test whether the regexp has at least one match in the given text.
// @param {String} str - The string to test.
// @param {String} expression - A regular expression to use for splitting the text.
// @param {Boolean} [caseInsensitive=False] - If set to true, the matching will be case-insensitive.
// @return {Boolean} - Returns true if the pattern is contained in text, and false otherwise.
func RegexTest(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	text := args[0].String()
	exp := args[1].String()

	if len(args) > 2 {
		if args[2] == values.True {
			exp = "(?i)" + exp
		}
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return values.None, err
	}

	matches := reg.MatchString(text)

	return values.NewBoolean(matches), nil
}

// REGEX_REPLACE replace every substring matched with the regexp with a given string.
// @param {String} str - The string to split.
// @param {String} expression - A regular expression search pattern.
// @param {String} replacement - The string to replace the search pattern with
// @param {Boolean} [caseInsensitive=False] - If set to true, the matching will be case-insensitive.
// @return {String} - Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
func RegexReplace(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	exp := args[1].String()
	repl := args[2].String()

	if len(args) > 3 {
		if args[3] == values.True {
			exp = "(?i)" + exp
		}
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return values.None, err
	}

	out := reg.ReplaceAllString(text, repl)

	return values.NewString(out), nil
}
