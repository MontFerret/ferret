package strings

import (
	"context"
	"regexp"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// regex_match returns the matches in the given string text, using the regex.
// @param str {String} The string to search in.
// @param expression {String} A regular expression to use for matching the text.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Any[]} An array of strings containing the matches.
func RegexMatch(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return regexMatch2(ctx, args[0], args[1])
	}

	return regexMatch3(ctx, args[0], args[1], args[2])
}

// regex_match returns the matches in the given string text, using the regex.
// @param str {String} The string to search in.
// @param expression {String} A regular expression to use for matching the text.
// @return {Any[]} An array of strings containing the matches.
func regexMatch2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return regexMatch3(ctx, arg1, arg2, runtime.False)
}

// regex_match returns the matches in the given string text, using the regex.
// @param str {String} The string to search in.
// @param expression {String} A regular expression to use for matching the text.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Any[]} An array of strings containing the matches.
func regexMatch3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	text := arg1.String()
	exp := arg2.String()

	if arg3 == runtime.True {
		exp = "(?i)" + exp
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return runtime.None, err
	}

	matches := reg.FindAllStringSubmatch(text, -1)
	res := runtime.NewArray(10)

	if len(matches) == 0 {
		return res, nil
	}

	for _, m := range matches[0] {
		_ = res.Append(ctx, runtime.NewString(m))
	}

	return res, nil
}

// regex_split splits the given string text into a list of strings, using the separator.
// @param str {String} The string to split.
// @param expression {String} A regular expression to use for splitting the text.
// @param limit {Int} Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @param reserved {Any} Reserved compatibility argument; currently ignored.
// @return {Any[]} An array of strings splitted by the expression.
func RegexSplit(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 4)

	if err != nil {
		return runtime.None, err
	}

	switch len(args) {
	case 2:
		return regexSplit2(ctx, args[0], args[1])
	case 3:
		return regexSplit3(ctx, args[0], args[1], args[2])
	default:
		return regexSplit4(ctx, args[0], args[1], args[2], args[3])
	}
}

// regex_split splits the given string text into a list of strings, using the separator.
// @param str {String} The string to split.
// @param expression {String} A regular expression to use for splitting the text.
// @return {Any[]} An array of strings splitted by the expression.
func regexSplit2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return regexSplit(ctx, arg1, arg2, -1)
}

// regex_split splits the given string text into a list of strings, using the separator.
// @param str {String} The string to split.
// @param expression {String} A regular expression to use for splitting the text.
// @param limit {Int} Limit the number of split values in the result.
// @return {Any[]} An array of strings splitted by the expression.
func regexSplit3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	limit := runtime.CastOr[runtime.Int](arg3, runtime.Int(-1))
	return regexSplit(ctx, arg1, arg2, int(limit))
}

// regex_split splits the given string text into a list of strings, using the separator.
// @param str {String} The string to split.
// @param expression {String} A regular expression to use for splitting the text.
// @param limit {Int} Limit the number of split values in the result.
// @param reserved {Any} Reserved compatibility argument; currently ignored.
// @return {Any[]} An array of strings splitted by the expression.
func regexSplit4(ctx context.Context, arg1, arg2, arg3, _ runtime.Value) (runtime.Value, error) {
	return regexSplit3(ctx, arg1, arg2, arg3)
}

func regexSplit(ctx context.Context, arg1, arg2 runtime.Value, limit int) (runtime.Value, error) {
	text := arg1.String()
	exp := arg2.String()

	reg, err := regexp.Compile(exp)

	if err != nil {
		return runtime.None, err
	}

	matches := reg.Split(text, limit)
	res := runtime.NewArray(10)

	if len(matches) == 0 {
		return res, nil
	}

	for _, m := range matches {
		_ = res.Append(ctx, runtime.NewString(m))
	}

	return res, nil
}

// regex_test test whether the regexp has at least one match in the given text.
// @param str {String} The string to test.
// @param expression {String} A regular expression to use for splitting the text.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
func RegexTest(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 3)

	if err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return regexTest2(ctx, args[0], args[1])
	}

	return regexTest3(ctx, args[0], args[1], args[2])
}

// regex_test test whether the regexp has at least one match in the given text.
// @param str {String} The string to test.
// @param expression {String} A regular expression to use for splitting the text.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
func regexTest2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return regexTest3(ctx, arg1, arg2, runtime.False)
}

// regex_test test whether the regexp has at least one match in the given text.
// @param str {String} The string to test.
// @param expression {String} A regular expression to use for splitting the text.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
func regexTest3(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	text := arg1.String()
	exp := arg2.String()

	if arg3 == runtime.True {
		exp = "(?i)" + exp
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return runtime.None, err
	}

	matches := reg.MatchString(text)

	return runtime.NewBoolean(matches), nil
}

// regex_replace replace every substring matched with the regexp with a given string.
// @param str {String} The string to split.
// @param expression {String} A regular expression search pattern.
// @param replacement {String} The string to replace the search pattern with
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive.
// @return {String} Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
func RegexReplace(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 3, 4)

	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 3 {
		return regexReplace3(ctx, args[0], args[1], args[2])
	}

	return regexReplace4(ctx, args[0], args[1], args[2], args[3])
}

// regex_replace replace every substring matched with the regexp with a given string.
// @param str {String} The string to split.
// @param expression {String} A regular expression search pattern.
// @param replacement {String} The string to replace the search pattern with
// @return {String} Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
func regexReplace3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return regexReplace4(ctx, arg1, arg2, arg3, runtime.False)
}

// regex_replace replace every substring matched with the regexp with a given string.
// @param str {String} The string to split.
// @param expression {String} A regular expression search pattern.
// @param replacement {String} The string to replace the search pattern with
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive.
// @return {String} Returns the string text with the search regex pattern replaced with the replacement string wherever the pattern exists in text
func regexReplace4(_ context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	text := arg1.String()
	exp := arg2.String()
	repl := arg3.String()

	if arg4 == runtime.True {
		exp = "(?i)" + exp
	}

	reg, err := regexp.Compile(exp)

	if err != nil {
		return runtime.None, err
	}

	out := reg.ReplaceAllString(text, repl)

	return runtime.NewString(out), nil
}
