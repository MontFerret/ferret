package strings

import (
	"context"
	"strings"
	"unicode"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// trim returns the string value with whitespace stripped from the start and/or end.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars on both sides.
func Trim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)
	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return trim1(ctx, args[0])
	}

	return trim2(ctx, args[0], args[1])
}

// trim returns the string value with whitespace stripped from the start and/or end.
// @param str {String} The string.
// @return {String} The string without chars on both sides.
func trim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.TrimSpace(string(text))), nil
}

// trim returns the string value with whitespace stripped from the start and/or end.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars on both sides.
func trim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, cutset, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.Trim(string(text), string(cutset))), nil
}

// ltrim returns the string value with whitespace stripped from the start only.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars at the left-hand side.
func LTrim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)
	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return lTrim1(ctx, args[0])
	}

	return lTrim2(ctx, args[0], args[1])
}

// ltrim returns the string value with whitespace stripped from the start only.
// @param str {String} The string.
// @return {String} The string without chars at the left-hand side.
func lTrim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.TrimLeftFunc(string(text), unicode.IsSpace)), nil
}

// ltrim returns the string value with whitespace stripped from the start only.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars at the left-hand side.
func lTrim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, cutset, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.TrimLeft(string(text), string(cutset))), nil
}

// rtrim returns the string value with whitespace stripped from the end only.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars at the right-hand side.
func RTrim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)
	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return rTrim1(ctx, args[0])
	}

	return rTrim2(ctx, args[0], args[1])
}

// rtrim returns the string value with whitespace stripped from the end only.
// @param str {String} The string.
// @return {String} The string without chars at the right-hand side.
func rTrim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.TrimRightFunc(string(text), unicode.IsSpace)), nil
}

// rtrim returns the string value with whitespace stripped from the end only.
// @param str {String} The string.
// @param chars {String} Literal rune cutset. An empty cutset leaves the string unchanged. Omit to trim Unicode whitespace.
// @return {String} The string without chars at the right-hand side.
func rTrim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, cutset, err := runtime.CastArgs2[runtime.String, runtime.String](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(strings.TrimRight(string(text), string(cutset))), nil
}
