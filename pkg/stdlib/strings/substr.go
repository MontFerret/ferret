package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SUBSTRING returns a substring of value.
// @param {String} str - The source string.
// @param {Int} offset - Start at offset, offsets start at position 0.
// @param {Int} [length] - At most length characters, omit to get the substring from offset to the end of the string.
// @return {String} - A substring of value.
func Substring(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.EmptyString, err
	}

	offsetArg, err := runtime.CastInt(args[1])

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	offset := int(offsetArg)
	length := size

	if len(args) > 2 {
		arg2, ok := args[2].(runtime.Int)

		if ok {
			length = int(arg2)
		}
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
func Left(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)

	var pos int
	arg1, ok := args[1].(runtime.Int)

	if ok {
		pos = int(arg1)
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
func Right(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	pos := size

	arg1, ok := args[1].(runtime.Int)

	if ok {
		pos = int(arg1)
	}

	if len(text) < pos {
		return runtime.NewString(text), nil
	}

	return runtime.NewStringFromRunes(runes[size-pos : size]), nil
}
