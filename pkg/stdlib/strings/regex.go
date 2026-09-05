package strings

import (
	"context"
	"regexp"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// regex_test reports whether a Go regular expression matches any part of text. Use inline flags such as (?i).
// @param text {String} The source string.
// @param pattern {String} The Go regular expression.
// @return {Boolean} Whether a match exists.
func RegexTest(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, expression, err := compileRegex(arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Boolean(expression.MatchString(text)), nil
}

// regex_find finds the first match of a Go regular expression. Use inline flags such as (?i). Captures exclude the full match, preserve declaration order, and use empty strings for unmatched groups. Named captures use the first declared group for duplicate names.
// @param text {String} The source string.
// @param pattern {String} The Go regular expression.
// @return {Object | None} An object with match (String), groups (String[]), and named (Object), or None when no match exists.
func RegexFind(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, expression, err := compileRegex(arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	match := expression.FindStringSubmatch(text)
	if match == nil {
		return runtime.None, nil
	}

	return regexMatchValue(expression, match), nil
}

// regex_find_all finds non-overlapping matches in source order using Go regular expression semantics. Empty matches adjacent to a preceding match are ignored. Use inline flags such as (?i). Capture fields follow regex_find.
// @param text {String} The source string.
// @param pattern {String} The Go regular expression.
// @return {Object[]} Objects with match (String), groups (String[]), and named (Object). No matches returns an empty array.
func RegexFindAll(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, expression, err := compileRegex(arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	matches := expression.FindAllStringSubmatch(text, -1)
	out := runtime.NewArray(len(matches))
	for _, match := range matches {
		_ = out.Append(ctx, regexMatchValue(expression, match))
	}

	return out, nil
}

// regex_replace replaces all non-overlapping matches of a Go regular expression. Use inline flags such as (?i).
// @param text {String} The source string.
// @param pattern {String} The Go regular expression.
// @param replacement {String} Go replacement template: $1 and ${name} expand captures, and $$ inserts a literal dollar sign.
// @return {String} The string with matches replaced.
func RegexReplace(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	text, expression, err := compileRegex(arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	replacement, err := runtime.CastArg[runtime.String](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(expression.ReplaceAllString(text, string(replacement))), nil
}

// regex_split splits text at matches of a Go regular expression. Use inline flags such as (?i). Empty matches follow Go regexp.Split semantics.
// @param text {String} The source string.
// @param pattern {String} The separator expression.
// @param limit {Int} Non-negative maximum result count. Zero returns an empty array; positive limits preserve the unsplit remainder. Omitted means unlimited.
// @return {String[]} The pieces between matches.
func RegexSplit(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return regexSplit2(ctx, args[0], args[1])
	}

	return regexSplit3(ctx, args[0], args[1], args[2])
}

// regex_split splits text at matches of a Go regular expression. Use inline flags such as (?i). Empty matches follow Go regexp.Split semantics.
// @param text {String} The source string.
// @param pattern {String} The separator expression.
// @return {String[]} All pieces between matches.
func regexSplit2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return regexSplit(ctx, arg1, arg2, -1)
}

// regex_split splits text at matches of a Go regular expression. Use inline flags such as (?i). Empty matches follow Go regexp.Split semantics.
// @param text {String} The source string.
// @param pattern {String} The separator expression.
// @param limit {Int} Non-negative maximum result count. Zero returns an empty array; positive limits preserve the unsplit remainder.
// @return {String[]} The pieces between matches.
func regexSplit3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	limit, err := nonNegativeLimit(arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	return regexSplit(ctx, arg1, arg2, limit)
}

func regexSplit(ctx context.Context, arg1, arg2 runtime.Value, limit int) (runtime.Value, error) {
	text, expression, err := compileRegex(arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	pieces := expression.Split(text, limit)
	out := runtime.NewArray(len(pieces))
	for _, piece := range pieces {
		_ = out.Append(ctx, runtime.String(piece))
	}

	return out, nil
}

func compileRegex(arg1, arg2 runtime.Value) (string, *regexp.Regexp, error) {
	text, pattern, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return "", nil, err
	}

	expression, err := regexp.Compile(string(pattern))
	if err != nil {
		return "", nil, runtime.ArgError(err, 1)
	}

	return string(text), expression, nil
}

func regexMatchValue(expression *regexp.Regexp, match []string) runtime.Value {
	groups := make([]runtime.Value, len(match)-1)
	for index, value := range match[1:] {
		groups[index] = runtime.String(value)
	}

	named := make(map[string]runtime.Value)
	for index, name := range expression.SubexpNames() {
		if name == "" {
			continue
		}

		if _, exists := named[name]; !exists {
			named[name] = runtime.String(match[index])
		}
	}

	return runtime.NewObjectWith(map[string]runtime.Value{
		"match":  runtime.String(match[0]),
		"groups": runtime.NewArrayWith(groups...),
		"named":  runtime.NewObjectWith(named),
	})
}
