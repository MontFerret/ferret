package strings

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SUBSTRING returns a substring of value.
// @param {String} str - The source string.
// @param {Int} offset - Start at offset, offsets start at position 0.
// @param {Int} [length] - At most length characters, omit to get the substring from offset to the end of the string.
// @return {String} - A substring of value.
func Substring(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 2 {
		return substring2(ctx, args[0], args[1])
	}

	return substring3(ctx, args[0], args[1], args[2])
}

func substring2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return substring(ctx, arg1, arg2, runtime.None, false)
}

func substring3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return substring(ctx, arg1, arg2, arg3, true)
}

func substring(_ context.Context, arg1, arg2, arg3 runtime.Value, hasLength bool) (runtime.Value, error) {
	offsetArg, err := runtime.CastArg[runtime.Int](arg2, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := arg1.String()
	runes := []rune(text)
	size := len(runes)
	offset := int(offsetArg)
	length := size

	if hasLength {
		length = int(runtime.CastOr[runtime.Int](arg3, runtime.Int(size)))
	}

	// Handle edge cases for bounds checking
	if offset < 0 || offset >= size {
		return runtime.NewString(""), nil
	}

	if length <= 0 {
		return runtime.NewString(""), nil
	}

	var substr []rune
	end := offset + length

	// Ensure end doesn't exceed the string size
	if end > size {
		end = size
	}

	substr = runes[offset:end]

	return runtime.NewStringFromRunes(substr), nil
}

// LEFT returns the leftmost characters of the string value by index.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The leftmost characters of the string value by index.
func Left(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)

	var pos int
	length, ok := arg2.(runtime.Int)

	if ok {
		pos = int(length)
	}

	if len(text) < pos {
		return runtime.NewString(text), nil
	}

	return runtime.NewStringFromRunes(runes[0:pos]), nil
}

// RIGHT returns the rightmost characters of the string value.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The rightmost characters of the string value.
func Right(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)
	size := len(runes)
	pos := size

	length, ok := arg2.(runtime.Int)

	if ok {
		pos = int(length)
	}

	if len(text) < pos {
		return runtime.NewString(text), nil
	}

	return runtime.NewStringFromRunes(runes[size-pos : size]), nil
}
