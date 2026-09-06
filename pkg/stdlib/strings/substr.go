package strings

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// substring returns a substring of value.
// @param str {String} The source string.
// @param offset {Int} Zero-based rune offset. Negative or out-of-range offsets return an empty string.
// @param length {Int} Non-negative maximum rune count; omitted means the rest of the string.
// @return {String} A substring of value.
func Substring(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 2 {
		return substring2(ctx, args[0], args[1])
	}

	return substring3(ctx, args[0], args[1], args[2])
}

// substring returns a substring of value.
// @param str {String} The source string.
// @param offset {Int} Zero-based rune offset. Negative or out-of-range offsets return an empty string.
// @return {String} A substring of value.
func substring2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return substring(ctx, arg1, arg2, runtime.None, false)
}

// substring returns a substring of value.
// @param str {String} The source string.
// @param offset {Int} Zero-based rune offset. Negative or out-of-range offsets return an empty string.
// @param length {Int} Non-negative maximum rune count; omitted means the rest of the string.
// @return {String} A substring of value.
func substring3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return substring(ctx, arg1, arg2, arg3, true)
}

func substring(_ context.Context, arg1, arg2, arg3 runtime.Value, hasLength bool) (runtime.Value, error) {
	text, offset, err := runtime.CastArgs2[runtime.String, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	runes := []rune(text)
	length := runtime.Int(len(runes))
	if hasLength {
		length, err = runtime.CastArg[runtime.Int](arg3, 2)
		if err != nil {
			return runtime.None, err
		}

		if length < 0 {
			return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "length must be non-negative"), 2)
		}
	}

	if offset < 0 || offset >= runtime.Int(len(runes)) {
		return runtime.EmptyString, nil
	}

	// Clamp before addition and conversion, including on 32-bit platforms.
	length = min(length, runtime.Int(len(runes))-offset)

	return runtime.NewStringFromRunes(runes[int(offset):int(offset+length)]), nil
}

// left returns the leftmost characters of the string value by index.
// @param str {String} The source string.
// @param length {Int} Non-negative rune count, clamped to the string length.
// @return {String} The leftmost characters of the string value by index.
func Left(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, length, err := runtime.CastArgs2[runtime.String, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	if length < 0 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "length must be non-negative"), 1)
	}

	runes := []rune(text)
	if length >= runtime.Int(len(runes)) {
		return text, nil
	}

	return runtime.NewStringFromRunes(runes[:int(length)]), nil
}

// right returns the rightmost characters of the string value.
// @param str {String} The source string.
// @param length {Int} Non-negative rune count, clamped to the string length.
// @return {String} The rightmost characters of the string value.
func Right(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, length, err := runtime.CastArgs2[runtime.String, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	if length < 0 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "length must be non-negative"), 1)
	}

	runes := []rune(text)
	if length >= runtime.Int(len(runes)) {
		return text, nil
	}

	return runtime.NewStringFromRunes(runes[len(runes)-int(length):]), nil
}
