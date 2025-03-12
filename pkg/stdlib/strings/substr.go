package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SUBSTRING returns a substring of value.
// @param {String} str - The source string.
// @param {Int} offset - Start at offset, offsets start at position 0.
// @param {Int} [length] - At most length characters, omit to get the substring from offset to the end of the string.
// @return {String} - A substring of value.
func Substring(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.EmptyString, err
	}

	err = core.ValidateType(args[1], types.Int)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	offset := int(args[1].(core.Int))
	length := size

	if len(args) > 2 {
		arg2, ok := args[2].(core.Int)

		if ok {
			length = int(arg2)
		}
	}

	var substr []rune

	if length == size {
		substr = runes[offset:]
	} else {
		end := offset + length

		if size > end {
			substr = runes[offset:end]
		} else {
			substr = runes[offset:]
		}
	}

	return core.NewStringFromRunes(substr), nil
}

// LEFT returns the leftmost characters of the string value by index.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The leftmost characters of the string value by index.
func Left(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)

	var pos int
	arg1, ok := args[1].(core.Int)

	if ok {
		pos = int(arg1)
	}

	if len(text) < pos {
		return core.NewString(text), nil
	}

	return core.NewStringFromRunes(runes[0:pos]), nil
}

// RIGHT returns the rightmost characters of the string value.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The rightmost characters of the string value.
func Right(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	pos := size

	arg1, ok := args[1].(core.Int)

	if ok {
		pos = int(arg1)
	}

	if len(text) < pos {
		return core.NewString(text), nil
	}

	return core.NewStringFromRunes(runes[size-pos : size]), nil
}
